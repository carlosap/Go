package httputil

import (
	"crypto/tls"
	"fmt"
	"github.com/Novetta/common/networking/networkutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

//HTTPProxy handles proxying requests incoming on port to the remote machine
func HTTPProxy(remoteURI string, port int) {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("Unable to get hostname: %v", err)
	}

	remoteParsedURI, err := url.Parse(remoteURI)
	if err != nil {
		log.Fatalf("Unable to parse remote URI %s: %v", remoteURI, err)
	}

	hostURI := fmt.Sprintf("%s:%d", hostname, port)

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Proxy:           http.ProxyFromEnvironment,
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		r.Host = hostURI
		proxy := httputil.NewSingleHostReverseProxy(remoteParsedURI)
		if remoteParsedURI.Scheme == "https" {
			proxy.Transport = transport
		}

		proxy.ServeHTTP(w, r)
	}

	http.HandleFunc("/", networkutil.LogRequests(handler))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

//BasicHTMLWriter writes the minimum necessary headers with the given body
func BasicHTMLWriter(body string, w http.ResponseWriter) {
	messageHeader := `<!DOCTYPE html>
<html>
    <head>
        <meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
        <meta http-equiv="X-UA-Compatible" content="IE-Edge"/>
    </head>
    <body>`
	messageTail := `    </body>
</html>`
	w.Header().Set("X-Frame-Options", "SAMEORIGIN")
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "%s\n%s\n%s", messageHeader, body, messageTail)
}
