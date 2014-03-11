package main

import (
	"fmt"
	"net/http"

	"flag"
)

var port = flag.Int("port", 5555, "port of the webserver")
var host = flag.String("host", "localhost", "host of the webserver")

func main() {
	flag.Parse()
	setupRouter()

	fmt.Printf("listening on %s:%d\n", *host, *port)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port), nil)
	if err != nil {
		panic(err.Error())
	}
}
