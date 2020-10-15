package main

import (
	"crypto/tls"
	"fmt"
	"golang.org/x/net/http2"
	"net"
	"net/http"
	"os"
	"strings"
	"strconv"
)

const URL = "http://localhost:7080"
const XFF = "1.1.1.1"

func main() {
	client := http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
		},
	}

	numOfGet := 0
	httpGet := func(){
		numOfGet++
		req, errReq := http.NewRequest("GET", URL, &strings.Reader{})
		checkErr(errReq, "during newRequest" + strconv.Itoa(numOfGet))
		req.Header.Add("X-Forwarded-For", XFF)
		resp, errResp := client.Do(req)
		
		// resp, err := client.Get(URL)
		checkErr(errResp, "during get" + strconv.Itoa(numOfGet))
		fmt.Println("\nResponse from server nr:", strconv.Itoa(numOfGet), "\n", resp, "\n")
	};

	// multiple requests on single connection
	httpGet()
	httpGet()
	httpGet()
}

func checkErr(err error, msg string) {
	if err == nil {
		return
	}
	fmt.Printf("ERROR: %s: %s\n", msg, err)
	os.Exit(1)
}
