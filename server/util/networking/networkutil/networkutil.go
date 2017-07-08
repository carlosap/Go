package networkutil

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Novetta/common/networking"
	"github.com/Novetta/common/services/wsclient"
	"github.com/Novetta/common/services/wsserver"
	"github.com/Novetta/common/util/logging"
	"github.com/Novetta/kerbproxy/kerbtypes"
)

func init() {
	log.SetOutput(os.Stdout)
}

//GzipResponseWriter returns gzipped data
type GzipResponseWriter struct {
	httpWriter http.ResponseWriter
	gzipWriter io.WriteCloser
}

//Init starts a writer
func (g *GzipResponseWriter) Init(w http.ResponseWriter, wr io.WriteCloser) {
	g.httpWriter = w
	g.gzipWriter = wr
}

func (g GzipResponseWriter) Write(b []byte) (int, error) {
	return g.gzipWriter.Write(b)
}

//Close a gzip writer
func (g *GzipResponseWriter) Close() error {
	return g.gzipWriter.Close()
}

//Header returns the http headers
func (g GzipResponseWriter) Header() http.Header {
	return g.httpWriter.Header()
}

//WriteHeader writes the http headers
func (g GzipResponseWriter) WriteHeader(p int) {
	g.httpWriter.WriteHeader(p)
}

//LogRequests logs http requests
func LogRequests(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Header.Get(kerbtypes.RemoteUser)
		var remote string
		if c := r.Header.Get(kerbtypes.XForwardedFor); len(c) > 0 {
			remote = c
		} else {
			remote = r.RemoteAddr
		}
		log.Printf(": %s : %s : %s", u, remote, r.RequestURI)
		f(w, r)
	}
}

//ProtectXSS stops xss vulnerabilites
func ProtectXSS(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")
		f(w, r)
	}
}

//GzipResponseHandler gzips http responses
func GzipResponseHandler(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		encoding := r.Header.Get("Accept-Encoding")
		writer := w
		var gzipResponse *GzipResponseWriter
		if encoding != "" && strings.Contains(strings.ToLower(encoding), "gzip") {
			gzipResponse = new(GzipResponseWriter)
			g, _ := gzip.NewWriterLevel(w, gzip.BestCompression)
			gzipResponse.Init(w, g)
			//defer gzipResponse.Close()
			writer = gzipResponse
			writer.Header().Set("Content-Encoding", "gzip")
		} /*else if encoding != "" && strings.Contains(strings.ToLower(encoding), "deflate") {
			gzipResponse = new(GzipResponseWriter)
			flater, _ := flate.NewWriter(w, flate.BestCompression)
			gzipResponse.Init(w, flater)
			writer = gzipResponse
			writer.Header().Set("Content-Encoding", "deflate")
		}*/
		f(writer, r)
		if gzipResponse != nil {
			gzipResponse.Close()
		}
	}
}

// GzipBytes Gzip a byte array
func GzipBytes(b []byte) []byte {
	buf := bytes.NewBuffer(b)
	return gzipByteBuffer(buf)
}

func gzipByteBuffer(buf *bytes.Buffer) []byte {
	outBuffer := new(bytes.Buffer)
	gzipWriter := gzip.NewWriter(outBuffer)
	if _, err := gzipWriter.Write(buf.Bytes()); err != nil {
		log.Printf("Unable to write to gzip: %v", err)
	}
	buf.Reset()
	gzipWriter.Close()

	return outBuffer.Bytes()
}

//NewPacketCompressor compresses streaming data
func NewPacketCompressor(inChan <-chan []byte, timeout time.Duration, maxPackets int, doLog bool) <-chan []byte {
	buf := new(bytes.Buffer)
	outChan := make(chan []byte, networking.ChanSize)
	timeoutTicker := time.Tick(timeout)
	logTicker := time.Tick(time.Minute)
	written := 0
	inSize := 0
	outSize := 0

	go func() {
		for {
			select {
			case b := <-inChan:
				if _, err := buf.Write(b); err != nil {
					log.Fatalf("Unable to write to buffer: %v", err)
				}
				inSize += len(b)
				written++
				if written >= maxPackets {
					outPacket := gzipByteBuffer(buf)
					outSize += len(outPacket)
					outChan <- outPacket
					buf.Reset()
					written = 0
				}
			case <-timeoutTicker:
				if written > 0 {
					outPacket := gzipByteBuffer(buf)
					outSize += len(outPacket)
					outChan <- outPacket
					buf.Reset()
					written = 0
				}
			case <-logTicker:
				if doLog {
					log.Printf("Total in: %d bytes. Total Out: %d bytes. Ratio: %.5f times compression ratio", inSize, outSize, float64(inSize)/float64(outSize))
				}
			}
		}
	}()

	return outChan
}

//IsCompressed determines if a slice of bytes contains compressed data.
//The return value is true if the data is compressed
//
func IsCompressed(data []byte) bool {
	if len(data) < 2 {
		return false
	}
	//Conforming to RFC1952 GZIP file format header
	if data[0] == 0x1f && data[1] == 0x8b {
		return true
	}
	return false
}

//Decompressor Defines a Decompressor type
//
type Decompressor struct {
	outChan       chan<- []byte
	netData       <-chan []byte
	stopChan      <-chan bool
	isCompressed  bool
	isInitialized bool
	gReader       *gzip.Reader
}

//NewNetworkDecompressor intitalizes a decompressor
func NewNetworkDecompressor(netChan string, outChan chan<- []byte, stopChan <-chan bool) {
	netData := make(chan []byte, cap(outChan))

	dec := &Decompressor{
		outChan:  outChan,
		netData:  netData,
		stopChan: stopChan,
	}

	err := networking.ListenMultiPortGroupStop(netChan, netData, stopChan)
	if err != nil {
		close(dec.outChan)
	}

	go dec.sendData()
}

func (dec *Decompressor) sendData() {
	for {
		select {
		case d := <-dec.netData:
			if !dec.isInitialized {
				if IsCompressed(d) {
					err := dec.Decompress(d)
					if err == nil {
						dec.isCompressed = true
					}
				}
			} else if dec.isCompressed {
				dec.Decompress(d)
			} else {
				dec.outChan <- d
			}
		case <-dec.stopChan:
			close(dec.outChan)
			return
		}
	}
}

//Decompress slice of bytes and pass it to a channel.
//If error returned is not nil, no bytes will be passed to the channel.
//
func (dec *Decompressor) Decompress(data []byte) error {
	var err error
	if dec.gReader == nil {
		dec.gReader, err = gzip.NewReader(bytes.NewReader(data))
	} else {
		err = dec.gReader.Reset(bytes.NewReader(data))
	}
	if err != nil {
		return err
	}
	output, err := ioutil.ReadAll(dec.gReader)
	if err != nil {
		return err
	}
	//No errors, send data to channel.
	dec.outChan <- output

	return nil
}

//ListenToAddress Recieve data on a given address and sends on the returned channel.
//Closes network connection on stop closed
func ListenToAddress(portgroup string, stop <-chan bool) (<-chan []byte, error) {
	u, err := ParseToURI(portgroup)
	if err != nil {
		return nil, err
	}

	var ch <-chan []byte
	switch u.Scheme {
	case "ws", "wss":
		//do websocket stuff
		tempCh := make(chan []byte, 1)
		go persistWebSocket(u, tempCh, stop)
		ch = tempCh
	case "tcp":
		//do tcp sutff
		if strings.Contains(u.Host, "127.0.0.1") || strings.Contains(u.Host, "0.0.0.0") || isLocal(u.Host) {
			ch = networking.ListenTCP(u.Host, stop)
		} else {
			ch = networking.ConnectTCP(u.Host, stop)
		}
	default:
		//do udp stuff
		tempCh := make(chan []byte, 10)
		i, err := IsMulticast(u.Host)
		if err == nil {
			if i {
				err = networking.ListenMultiPortGroupStop(u.Host, tempCh, stop)
				ch = tempCh
			} else {
				ch = networking.ListenUDPPort(u.Host)
			}
		}
	}
	return ch, err
}

func isLocal(uri string) bool {
	host, _, _ := net.SplitHostPort(uri)
	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {
		logging.Debug("Checking local address %q for %q", addr.String(), host)
		if strings.Contains(addr.String(), host) {
			return true
		}
	}
	return false
}

func persistWebSocket(u *url.URL, dataChan chan<- []byte, stop <-chan bool) {
	for {
		select {
		case <-stop:
			close(dataChan)
			return
		default:
			//create a temporary channel
			//it will be closed if something happens to the websocket in the underlying library
			s := make(chan []byte, 1)
			go fanWsChan(dataChan, s)

			err := wsclient.StartWsClient(u.String(), s, stop)
			if err != nil {
				log.Printf("Error reading from websocket, disconnected: '%s': %+v", u.String(), err)
				time.Sleep(time.Second * 5)
			} else {
				log.Printf("Disconnected from websocket: '%s': no error", u.String())
				time.Sleep(time.Second * 5)
			}
		}
	}
}

func fanWsChan(send chan<- []byte, recieve <-chan []byte) {
	for d := range recieve {
		send <- d
	}
}

//SendToAddress Send data from the data chan to the given address.
//This can handle websockets and udp
//channel expected to be a websocket SecureMessage or []byte
func SendToAddress(portgroup string, data <-chan interface{}) error {
	u, err := ParseToURI(portgroup)
	if err != nil {
		return err
	}

	if u.Scheme == "ws" || u.Scheme == "wss" {
		outChan := make(chan wsserver.SecureMessage)
		go func() {
			for d := range data {
				if secMsg, ok := d.(wsserver.SecureMessage); ok {
					outChan <- secMsg
				} else {
					logging.Errorf("SendToAddress data is not wsserver.SecureMessage for scheme %s, got %T", u.Scheme, d)
				}
			}
		}()
		wsserver.SendOnWebsocket(u.String(), outChan)
	} else {
		outChan := make(chan []byte)
		go func() {
			for d := range data {
				if b, ok := d.([]byte); ok {
					outChan <- b
				} else {
					logging.Errorf("SendToAddress data is not []byte for scheme %s, got %T", u.Scheme, b)
				}
			}
		}()
		switch u.Scheme {
		case "tcp":
			splitGroup := strings.Split(portgroup, "://")
			if !strings.Contains(portgroup, "127.0.0.1") && isLocal(splitGroup[len(splitGroup)-1]) {
				logging.Info("TCP listening on %q", portgroup)
				networking.ListenSendTCP(u.Host, outChan)
			} else {
				logging.Info("TCP sending on %q", portgroup)
				networking.RedialSendTCP(u.Host, outChan)
			}
		default:
			m, err := IsMulticast(u.Host)
			if err == nil {
				if m {
					networking.SendUDPPortGroup(u.Host, outChan)
				} else {
					networking.SendUDP(u.Host, outChan)
				}
			}
		}
	}
	return err
}

//IsMulticast takes a string reprensting a ip:port and returns if it represents a multicast address, and an error
func IsMulticast(ipString string) (bool, error) {
	//make sure there is a port in the address
	ip := ""
	if !strings.Contains(ipString, ":") {
		ip = ipString
	} else {
		var err error
		ip, _, err = net.SplitHostPort(ipString)
		if err != nil {
			return false, err
		}
	}
	i := net.ParseIP(ip)
	if i == nil {
		return false, fmt.Errorf("IP '%s' is not a valid IP address", ip)
	}
	return i.IsMulticast(), nil
}

//ParseToURI the given string to a URI
func ParseToURI(portgroup string) (*url.URL, error) {
	portgroup = strings.TrimSpace(portgroup)
	if !strings.Contains(portgroup, "://") {
		portgroup = "udp://" + portgroup
	}
	return url.Parse(portgroup)
}

//WriteUserHeader a custom user header for local auth purposes
func WriteUserHeader(r *http.Request, user *kerbtypes.User) {
	uData, _ := json.MarshalIndent(user, "", " ")
	uEncoded := base64.StdEncoding.EncodeToString(uData)
	r.Header.Set(kerbtypes.RemoteUser, uEncoded)
}

//DisableTLSCheck the default TLS check for all client connections
func DisableTLSCheck() {
	if t, ok := http.DefaultTransport.(*http.Transport); ok {
		t.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
		http.DefaultTransport = t
	}
}

//GetBasePath returns the host from an http request
func GetBasePath(r *http.Request) *url.URL {
	u := &url.URL{}
	fowardHost := r.Header.Get(kerbtypes.XForwardedHost)

	if fowardHost != "" {
		u.Host = fowardHost
	} else {
		u.Host = r.Host
	}

	scheme := r.Header.Get(kerbtypes.XForwardedProto)
	if scheme != "" {
		u.Scheme = r.Header.Get(kerbtypes.XForwardedProto)
	} else if kerbtypes.EnableSSL {
		u.Scheme = "https"
	} else {
		u.Scheme = "http"
	}

	path := strings.Split(strings.TrimPrefix(r.Header.Get(kerbtypes.XForwardedURI), "/"), "/")
	if path[0] != "" {
		u.Path = "/" + path[0]
	} else {
		u.Path = "/" + strings.Split(r.URL.Path, "/")[0]
	}

	return u
}
