package main
import (
	"net/http"

)

func main() {
	http.Handle("/", http.FileServer(http.Dir("public")))
	http.ListenAndServeTLS(":3000", "cert.pem", "key.pem", nil)
	
}