package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	port = flag.Int("port", 8080, "port for server to listen on")
	dir  = flag.String("dir", "./web", "directory to serve")
)

func main() {
	flag.Parse()
	log.Printf("HTTP server started: http://localhost:%d...", *port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), http.FileServer(http.Dir(*dir)))
	log.Fatalln(err)
}
