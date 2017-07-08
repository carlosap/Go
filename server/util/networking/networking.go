package networking

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/Novetta/common/color"
	"github.com/Novetta/common/util"

	"golang.org/x/net/ipv4"
)

const (
	// ChanSize is the size of the buffered channel for network data
	ChanSize = 100

	//This is the maximum number of bytes in a single udp packet
	maxPacketSize = 65535 * 2

	//MulticastInterface ENV Variable MULTICAST_INTERFACE controls the default interface
	MulticastInterface = "MULTICAST_INTERFACE"

	//The TCP version we are using (denied tcp6)
	tcpVer = "tcp4"

	//sleepTimeEnv allows setting the duration between tcp dial attempts
	sleepTimeEnv = "SLEEP_TIME"
)

var (
	ttl       = 30
	sleepTime time.Duration
)

func init() {
	if t := os.Getenv("TTL"); len(t) > 0 {
		tt, err := strconv.Atoi(t)
		if err != nil {
			log.Fatalf("ENV TTL is not in valid format: %s", t)
		}
		ttl = tt
	}

	var err error
	sleepTime, err = time.ParseDuration(util.GetSetEnv(sleepTimeEnv, "1s"))
	if err != nil {
		log.Fatalf("%q env is not in correct format %+v", sleepTimeEnv, err)
	}
}

// GetMulticastAddress gets the multicast address of the best multicast-enabled interface
func GetMulticastAddress() string {
	iface, _ := GetDefaultMulticastInterface()
	if iface != nil {
		addrs, err := iface.Addrs()
		if err != nil {
			log.Fatalf(color.Red.ColorString("Unable to get addresses for '%s': %+v"), iface.Name, err)
		} else if addrs == nil || len(addrs) == 0 {
			log.Fatalf(color.Red.ColorString("No addresses for given interface: %s"), iface.Name)
		} else {
			return addrs[0].String()
		}
	}

	return ""
}

//GetDefaultMulticastInterface Gets the default Multicast Interface used by the running machine
func GetDefaultMulticastInterface() (*net.Interface, error) {
	return getDefaultMulticastInterface()
}

func getDefaultMulticastInterface() (*net.Interface, error) {
	ifaceName := os.Getenv(MulticastInterface)
	if ifaceName != "" {
		ifi, err := net.InterfaceByName(ifaceName)
		log.Printf("Setting default multicast interface to: %s (%+v)", ifaceName, ifi)
		return ifi, err
	}

	// Check other interfaces
	ifaces, _ := net.Interfaces()
	for _, ni := range ifaces {
		if ni.Flags&net.FlagMulticast == net.FlagMulticast {
			log.Printf("Found Multicast interface: %s", ni.Name)
			return &ni, nil
		}
	}

	// That's weird
	return nil, fmt.Errorf("No multicast interface detected")
}

//ListenUDPPortChannel Listens on a udp4 port returning packet payloads on the result channel
func ListenUDPPortChannel(portgroup string, resultChan chan<- []byte) error {
	udpAddr, err := net.ResolveUDPAddr("udp4", portgroup)
	if err != nil {
		log.Printf("%v\n", err)
		return err
	}

	udpConn, err := net.ListenUDP("udp4", udpAddr)
	if err != nil {
		log.Printf("%v\n", err)
		return err
	}

	log.Printf("Listening on %v.", portgroup)

	stopChan := make(chan bool)
	go HandleReadConn(udpConn, resultChan, stopChan)
	return nil
}

//ListenUDPPort Listen on a udp4 port returning a channel that will have the packet payloads
func ListenUDPPort(portgroup string) <-chan []byte {
	resultChan := make(chan []byte)
	err := ListenUDPPortChannel(portgroup, resultChan)
	if err != nil {
		log.Printf("%+v", err)
		close(resultChan)
		return resultChan
	}

	return resultChan
}

//ListenUDPPortGroup Listen on a udp4 port, that may be multicast, and returns a channel that will have the packet payloads
func ListenUDPPortGroup(portgroup string) (<-chan []byte, error) {
	resultChan := make(chan []byte, ChanSize)

	if runtime.GOOS == "linux" {
		return ListenMultiPortGroup(portgroup)
	}

	udpAddr, err := net.ResolveUDPAddr("udp4", portgroup)
	if err != nil {
		log.Printf("%v\n", err)
		close(resultChan)
		return resultChan, err
	}

	ifi, _ := GetDefaultMulticastInterface()

	udpConn, err := net.ListenMulticastUDP("udp4", ifi, udpAddr)
	if err != nil {
		log.Printf("%v\n", err)
		close(resultChan)
		return resultChan, err
	}

	log.Printf("Listening on %v.", portgroup)

	stopChan := make(chan bool)
	go HandleReadConn(udpConn, resultChan, stopChan)
	return resultChan, nil
}

//ListenMultiPortGroupChannel User provided result channel
func ListenMultiPortGroupChannel(portGroup string, resultChan chan<- []byte) error {
	stopChan := make(chan bool)
	return ListenMultiPortGroupStop(portGroup, resultChan, stopChan)
}

//ListenMultiPortGroupStop Includes a stop channel to stop listening on the port on request
func ListenMultiPortGroupStop(portGroup string, resultChan chan<- []byte, stopChan <-chan bool) error {
	udpAddr, err := net.ResolveUDPAddr("udp4", portGroup)
	if err != nil {
		log.Printf("%v\n", err)
		return err
	}

	log.Printf("Listening (readonly) on %v.", portGroup)

	wroteSize := false
	messageBytes := make([]byte, maxPacketSize*2)

	// -- File Handle
	f := getFileHandleFromPortGroup(portGroup, udpAddr)
	if f == nil {
		return fmt.Errorf("Error getting handle from getFileHandleFromPortGroup portGroup=%s", portGroup)
	}
	c, err := net.FilePacketConn(f)
	if err != nil {
		log.Printf("Error getting FilePacketConn for %v: %v", f, err)
		return err
	}
	f.Close()

	// Set the default multicast interface, if any
	ifi, _ := GetDefaultMulticastInterface()
	if err := ipv4.NewPacketConn(c).JoinGroup(ifi, udpAddr); err != nil {
		log.Printf(err.Error())
		return err
	}

	fillData := make([]byte, maxPacketSize)

	go func() {
		log.Printf("reading in ListenMultiPortGroupStop")
		skipped := 0
		sent := 0
		total := 0
		//Don't use time.Tick exit will mem-leak
		ticker := time.NewTicker(time.Second)
		defer c.Close()
		for {
			select {
			case <-stopChan:
				return
			default:
				c.SetReadDeadline(time.Now().Add(time.Second))
				numBytes, _, err := c.ReadFrom(messageBytes)
				if err != nil {
					if e, ok := err.(*net.OpError); ok && e.Timeout() {
						continue
					}
					log.Printf("Multicast read error %T: %v", err, err)
				}
				if numBytes > 0 {
					if !wroteSize {
						log.Printf("%s is %d bytes", portGroup, numBytes)
						wroteSize = true
					}
					select {
					case resultChan <- messageBytes[:numBytes:numBytes]:
						messageBytes = messageBytes[numBytes:]
						if len(messageBytes) < maxPacketSize {
							messageBytes = append(messageBytes, fillData...)
						}
						sent++
						total++
					case <-ticker.C:
						if skipped > 0 {
							log.Printf("%sSkipped %d of %d in last second%s", color.Red, skipped, sent+skipped, color.Normal)
							skipped = 0
							sent = 0
						}
						if total > 100000 {
							log.Printf("\n\n100000\n\n")
							total = 0
						}
					default:
						skipped++
					}
				}
			}
		}
	}()
	return nil
}

func getFileHandleFromPortGroup(portGroup string, udpAddr *net.UDPAddr) *os.File {
	s, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		log.Printf("Error getting socket: %v", err)
		return nil
	}

	syscall.SetsockoptInt(s, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	ipa := udpAddr.IP.To4()
	if ipa == nil || len(ipa) != 4 {
		log.Printf("Issue parsing address: %#v, %#v", *udpAddr, udpAddr.IP)
		return nil
	}
	ipa4 := [4]byte{ipa[0], ipa[1], ipa[2], ipa[3]}
	lsa := &syscall.SockaddrInet4{Port: udpAddr.Port, Addr: ipa4}
	if err := syscall.Bind(s, lsa); err != nil {
		log.Printf("Error binding socket: %v", err)
		return nil
	}

	return os.NewFile(uintptr(s), portGroup)
}

//ConnectTCP Connects to a server over raw TCP
func ConnectTCP(portGroup string, stopChan <-chan bool) <-chan []byte {
	resultChan := make(chan []byte, ChanSize)

	conn, err := net.Dial("tcp", portGroup)
	if err != nil {
		log.Printf("Error dialing %s: %v", portGroup, err)
		close(resultChan)
		return resultChan
	}

	go HandleReadConn(conn, resultChan, stopChan)

	return resultChan
}

//RedialSendTCP continuously dials the network connection and sends data when it is established
func RedialSendTCP(portGroup string, data <-chan []byte) {
	isOpen := true
	for isOpen {
		conn, err := net.Dial("tcp", portGroup)
		if err != nil {
			log.Printf("Error dialing %s: %v", portGroup, err)
			time.Sleep(sleepTime)
			continue
		}
		isOpen = HandleSendConn(conn, data)
		time.Sleep(sleepTime)
	}
}

//ListenSendTCP listens for incoming connections and sends the data to them
func ListenSendTCP(portGroup string, data <-chan []byte) error {
	lAddr, err := net.ResolveTCPAddr(tcpVer, portGroup)
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP(tcpVer, lAddr)
	if err != nil {
		return err
	}

	acceptSendConns(listener, data)

	return err
}

func acceptSendConns(l *net.TCPListener, data <-chan []byte) {
	defer l.Close()
	isOpen := true
	for isOpen {
		clientConn, err := l.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %+v", err)
		} else {
			isOpen = HandleSendConn(clientConn, data)
		}
	}
}

//HandleSendConn sends data on the given conn and returns if the conn is broken or the channel closes
func HandleSendConn(conn net.Conn, data <-chan []byte) (open bool) {
	var err error
	for d := range data {
		_, err = conn.Write(d)
		if err != nil {
			log.Printf("Error sending to %q", conn.RemoteAddr().String())
			return true
		}
	}

	return false
}

//ListenTCP listens for incoming packets on a tcp channel
func ListenTCP(portGroup string, stopChan <-chan bool) <-chan []byte {
	outChan := make(chan []byte)
	err := startServer(portGroup, outChan, stopChan)
	if err != nil {
		log.Printf("Error starting TCP server: %+v", err)
	}
	return outChan
}

func startServer(portGroup string, outChan chan<- []byte, stopChan <-chan bool) error {
	lAddr, err := net.ResolveTCPAddr(tcpVer, portGroup)
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP(tcpVer, lAddr)
	if err != nil {
		return err
	}

	go acceptConns(listener, outChan, stopChan)

	return nil
}

func acceptConns(l *net.TCPListener, outChan chan<- []byte, stopChan <-chan bool) {
	defer l.Close()
	defer close(outChan)
	for {
		clientConn, err := l.Accept()
		select {
		case <-stopChan:
			log.Printf("Stopping")
			return
		default:
			if err != nil {
				log.Printf("Error accepting connection: %+v", err)
			} else {
				go handleConn(clientConn, outChan, stopChan)
			}
		}
	}
}

func handleConn(conn net.Conn, outChan chan<- []byte, stopChan <-chan bool) {
	log.Printf("Opening %s", conn.RemoteAddr())
	//You need to make an extra channel here because the underlying HandleReadConn
	//closes its output channel on return, and we have multiple connectios on
	//different threads running.  We don't want to break all threads when 1 exits
	dataChan := make(chan []byte)
	go HandleReadConn(conn, dataChan, stopChan)
	forwardData(dataChan, outChan)
	conn.Close()
	log.Printf("Closing %s", conn.RemoteAddr())
}

func forwardData(data <-chan []byte, outChan chan<- []byte) {
	var d []byte
	for d = range data {
		outChan <- d
	}
}

//HandleReadConn handles a network connection, reading off of the connection 1 packet at a time
//and sending it up the dataChan.  This will exit when stopChan is closed
func HandleReadConn(conn net.Conn, dataChan chan<- []byte, stopChan <-chan bool) {
	defer close(dataChan)
	if conn == nil {
		return
	}
	defer conn.Close()
	messageBytes := make([]byte, maxPacketSize*2)
	for {
		select {
		case <-stopChan:
			return
		default:
			//Set a ReadDeadline to handle an idle timeout.  This will ensure it doesn't hang forever
			//if the connection should be closed
			conn.SetReadDeadline(time.Now().Add(time.Second * 30))
			n, err := conn.Read(messageBytes)
			if err != nil {
				if err == io.EOF {
					//Don't log, but return because conn is closed
					return
				} else if e, ok := err.(*net.OpError); ok && (e.Timeout() || e.Temporary()) {
					//Do nothing, we just don't want to wait forever to close connection
					//if stopChan had been closed
					//or
					//Do nothing because the error should resolve itself (the error is temporary)
				} else {
					log.Printf("Error reading from connection local: %s remote: %s: %+v", conn.LocalAddr().String(), conn.RemoteAddr().String(), err)
					return
				}
			} else if n > 0 {
				dataChan <- messageBytes[:n:n]
				messageBytes = messageBytes[n:]
				if len(messageBytes) <= maxPacketSize {
					messageBytes = append(messageBytes, make([]byte, maxPacketSize)...)
				}
			}
		}
	}
}

//ListenMultiPortGroup Listens to UDP multicast ipv4 for the specified ip group and port
func ListenMultiPortGroup(portGroup string) (<-chan []byte, error) {
	stopChan := make(chan bool)
	resultChan := make(chan []byte, ChanSize)
	err := ListenMultiPortGroupStop(portGroup, resultChan, stopChan)
	return resultChan, err
}

//NetworkChanReader Wrap our network channel into actual bytes
type NetworkChanReader struct {
	byteChan <-chan []byte
	buf      []byte
}

// Comply with the io.Reader interface
func (nc *NetworkChanReader) Read(p []byte) (n int, err error) {
	for len(nc.buf) < len(p) {
		// Grab the next packet and append the bytes to our buffer
		nc.buf = append(nc.buf, <-nc.byteChan...)
	}

	n = copy(p, nc.buf)
	nc.buf = nc.buf[n:]
	return
}

//NewNetworkChanReader Constructor
func NewNetworkChanReader(inputChan <-chan []byte) io.Reader {
	return &NetworkChanReader{
		byteChan: inputChan,
		buf:      make([]byte, 0),
	}
}

//SendUDPPortGroup sends data over a udp4 socket.
//This socket may be a multicast socket
func SendUDPPortGroup(portGroup string, outputChan <-chan []byte) {
	log.Printf("Output Feed: %s", portGroup)
	udpAddr, err := net.ResolveUDPAddr("udp4", portGroup)
	if err != nil {
		log.Panicf("UDP Error: %v\n", err)
	}

	rawUDPConn, err := net.ListenMulticastUDP("udp4", nil, udpAddr)
	if err != nil {
		log.Panicf("Raw UDP Error: %v\n", err)
	}

	udpConn := ipv4.NewPacketConn(rawUDPConn)
	udpConn.SetMulticastLoopback(true)
	udpConn.SetMulticastTTL(ttl)
	udpConn.SetTOS(0x28) // AF11

	ifi, _ := GetDefaultMulticastInterface()
	if ifi != nil {
		udpConn.SetMulticastInterface(ifi)
	}

	for chanResult := range outputChan {
		if chanResult != nil && len(chanResult) > 0 {
			_, err := udpConn.WriteTo(chanResult, nil, udpAddr)
			if err != nil {
				log.Printf("Error when sending %v", err)
			}
		}
	}
}

//SendUDP sends data over a udp4 socket.
func SendUDP(portGroup string, outputChan <-chan []byte) {
	log.Printf("Output Feed: %s", portGroup)

	udpConn, err := net.Dial("udp4", portGroup)
	if err != nil {
		log.Panicf("Raw UDP Error: %v\n", err)
	}

	for chanResult := range outputChan {
		if chanResult != nil && len(chanResult) > 0 {
			_, err := udpConn.Write(chanResult)
			if err != nil {
				log.Printf("Error when sending %v", err)
			}
		}
	}
}
