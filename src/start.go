package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func main() {
	port := 5678
	if len(os.Args) > 1 {
		p, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid port: %q\n", os.Args[1])
			os.Exit(1)
		}
		port = p
	}
	addr := fmt.Sprintf(":%d", port)
	fs := http.FileServer(http.Dir("."))
	http.ListenAndServe(addr, fs)
}
