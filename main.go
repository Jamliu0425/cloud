package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"strings"
)

func index(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("<h1>Welcome to CloudNative.</h1>"))

	os.Setenv("VERSION", "v111111")
	version := os.Getenv("VERSION")
	w.Header().Set("VERSION", version)
	fmt.Printf("OS Version: %s\n", version)

	for k, v := range r.Header {
		fmt.Println(k, v)
		for _, vv := range v {
			fmt.Printf("header key: %s, header value: %s\n", k, vv)
			w.Header().Set(k, vv)
		}
	}

	clientIp := ClientIP(r)

	log.Printf("Client Ip: %s\n", clientIp)
	// HTTP return code
	log.Printf("Code: %d\n", http.StatusAccepted)

}

func ClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}
	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}

func healthz(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "UP")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/healthz", healthz)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Starting HttpServer Failed, error: %s", err.Error())
	}
}
