package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/hartfordfive/request-forwarder/proxy"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var flagAddr = flag.String("a", "127.0.0.1", "The local address to listen on.")
	var flagPort = flag.Int("p", 8080, "The local port to listen on.")
	var flagRemoteAddr = flag.String("ra", "127.0.0.1", "The remoate address to bind to.")
	var flagRemotePort = flag.Int("rp", 0, "The remote port to bind to.")
	var flagMetricsPort = flag.Int("m", 9555, "The port on which to expose Prometheus metrics.")
	var flagAllowedMethods = flag.String("w", "", "Comma separated list of allowed methods. Empty means all.")
	flag.Parse()

	if *flagRemotePort < 1 {
		log.Fatal("The remote port flag (-rp) must be specified with a non-zero value!")
	}

	hostport := fmt.Sprintf("%s:%d", *flagAddr, *flagPort)
	prometheusHostport := fmt.Sprintf("%s:%d", *flagAddr, *flagMetricsPort)

	handler := proxy.NewProxy(*flagRemoteAddr, *flagRemotePort, *flagAllowedMethods)

	go func() {
		log.Println("Serving prometheus metrics on ", prometheusHostport)
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(prometheusHostport, nil))
	}()

	log.Println("Starting proxy server on", hostport)
	if err := http.ListenAndServe(hostport, handler); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
