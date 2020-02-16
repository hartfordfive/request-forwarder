package proxy

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/hartfordfive/request-forwarder/lib"
	//"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Proxy struct is the object used for the proxy serer
type proxy struct {
	AllowdMethods []string // Can include GET, POST, PUT, PATCH, DELETE, HEAD
	RemoteHost    string
	RemotePort    int
	RemoteTimeout int
}

// NewProxy returns a new proxy instance to intercept the requests
func NewProxy(addr string, port int, methods string) *proxy {
	var allowedMethods []string

	if strings.Trim(methods, " ") == "" {
		allowedMethods = append(allowedMethods, http.MethodGet)
		allowedMethods = append(allowedMethods, http.MethodPost)
		allowedMethods = append(allowedMethods, http.MethodHead)
		allowedMethods = append(allowedMethods, http.MethodPut)
		allowedMethods = append(allowedMethods, http.MethodDelete)
		allowedMethods = append(allowedMethods, http.MethodTrace)
	} else {
		methods := strings.Split(methods, ",")
		for _, method := range methods {
			m := strings.Trim(method, " ")
			if strings.ToLower(m) == "get" {
				allowedMethods = append(allowedMethods, http.MethodGet)
			} else if strings.ToLower(m) == "post" {
				allowedMethods = append(allowedMethods, http.MethodPost)
			} else if strings.ToLower(m) == "head" {
				allowedMethods = append(allowedMethods, http.MethodHead)
			} else if strings.ToLower(m) == "put" {
				allowedMethods = append(allowedMethods, http.MethodPut)
			} else if strings.ToLower(m) == "delete" {
				allowedMethods = append(allowedMethods, http.MethodDelete)
			} else if strings.ToLower(m) == "trace" {
				allowedMethods = append(allowedMethods, http.MethodTrace)
			}
		}
	}

	return &proxy{
		AllowdMethods: allowedMethods,
		RemoteHost:    addr,
		RemotePort:    port,
		RemoteTimeout: 2, // seconds
	}
}

func (p *proxy) ServeHTTP(wr http.ResponseWriter, inReq *http.Request) {
	log.Println(inReq.RemoteAddr, " ", inReq.Method, " ", inReq.URL)

	// if req.URL.Scheme != "http" && req.URL.Scheme != "https" {
	// 	msg := "Only HTTP and HTTPS are supported!"
	// 	http.Error(wr, msg, http.StatusBadRequest)
	// 	log.Print(msg)
	// 	return
	// }

	req, err := http.NewRequest(inReq.Method, fmt.Sprintf("http://%s:%d%s", p.RemoteHost, p.RemotePort, inReq.URL.RequestURI()), nil)
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}
	req.Header.Set("Cache-Control", "no-cache")
	client := &http.Client{Timeout: time.Second * time.Duration(p.RemoteTimeout)}

	//http: Request.RequestURI can't be set in client requests.
	//http://golang.org/src/pkg/net/http/client.go
	inReq.RequestURI = ""

	lib.DelHopHeaders(inReq.Header)

	if clientIP, _, err := net.SplitHostPort(inReq.RemoteAddr); err == nil {
		lib.AppendHostToXForwardHeader(req.Header, clientIP)
	}

	// Check if the method is in the list of AllowdMethods
	if _, exists := lib.ExistsInSlice(p.AllowdMethods, inReq.Method); !exists {
		msg := fmt.Sprintf("HTTP method %s not allowed", inReq.Method)
		http.Error(wr, msg, http.StatusBadRequest)
		log.Print(msg)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		http.Error(wr, "Server Error", http.StatusInternalServerError)
		log.Fatal("ServeHTTP:", err)
	}
	defer resp.Body.Close()

	log.Println(req.RemoteAddr, " ", resp.Status)

	lib.DelHopHeaders(resp.Header)

	lib.CopyHeader(wr.Header(), resp.Header)
	wr.WriteHeader(resp.StatusCode)
	io.Copy(wr, resp.Body)
}
