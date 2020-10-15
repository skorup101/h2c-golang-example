package main

import (
	"fmt"
	"golang.org/x/net/http2"
	"net"
	"net/http"
	"os"
)

func checkErr(err error, msg string) {
	if err == nil {
		return
	}
	fmt.Printf("ERROR: %s: %s\n", msg, err)
	os.Exit(1)
}

const serverAddr = "localhost:7080"
const XFF = "1.1.1.1"

func main() {
	server := http2.Server{}

	l, err := net.Listen("tcp", serverAddr)
	checkErr(err, "while listening")

	fmt.Printf("Listening [%s]...\n", serverAddr)
	for {
		conn, err := l.Accept()
		checkErr(err, "during accept")

		server.ServeConn(conn, &http2.ServeConnOpts{
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Printf("\nNew Request from client: %+v\n\n", r)
				w.Header().Add("X-Forwarded-For", XFF)
				fmt.Fprintf(w, "Hello, %v, http: %v", r.URL.Path, r.TLS == nil)
			}),
		})
	}
}
