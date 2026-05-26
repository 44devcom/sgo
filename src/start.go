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
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "get working directory: %v\n", err)
		os.Exit(1)
	}
	addr := fmt.Sprintf(":%d", port)
	url := fmt.Sprintf("http://localhost:%d/", port)

	fmt.Printf("sgo: static file server\n")
	fmt.Printf("  directory: %s\n", cwd)
	fmt.Printf("  URL:       %s\n", url)
	fmt.Printf("  stop:      press Ctrl+C\n")

	fs := http.FileServer(http.Dir("."))
	if err := http.ListenAndServe(addr, fs); err != nil {
		fmt.Fprintf(os.Stderr, "listen on %s: %v\n", addr, err)
		os.Exit(1)
	}
}
