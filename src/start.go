package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type config struct {
	port int
	dir  string
}

func parseConfig(args []string) (config, error) {
	fs := flag.NewFlagSet("sgo", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	port := 5678
	dir := "."
	fs.IntVar(&port, "port", 5678, "HTTP listen port")
	fs.StringVar(&dir, "dir", ".", "directory to serve")

	if err := fs.Parse(args); err != nil {
		return config{}, err
	}

	portExplicit := false
	dirExplicit := false
	fs.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "port":
			portExplicit = true
		case "dir":
			dirExplicit = true
		}
	})

	pos := fs.Args()
	if len(pos) == 1 && !portExplicit {
		if p, err := strconv.Atoi(pos[0]); err == nil {
			if p < 1 || p > 65535 {
				return config{}, fmt.Errorf("invalid port: %d", p)
			}
			port = p
			pos = nil
		}
	}

	if len(pos) > 0 {
		return config{}, fmt.Errorf("unexpected arguments: %s (use -dir for a serve directory)", strings.Join(pos, " "))
	}

	if !dirExplicit {
		dir = "."
	}

	if port < 1 || port > 65535 {
		return config{}, fmt.Errorf("invalid port: %d", port)
	}

	return config{port: port, dir: dir}, nil
}

func resolveRootDir(dir string) (string, error) {
	abs, err := filepath.Abs(filepath.Clean(dir))
	if err != nil {
		return "", fmt.Errorf("resolve directory: %w", err)
	}
	info, err := os.Stat(abs)
	if err != nil {
		return "", fmt.Errorf("directory %q: %w", dir, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("%q is not a directory", abs)
	}
	return abs, nil
}

func main() {
	cfg, err := parseConfig(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	root, err := resolveRootDir(cfg.dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	addr := fmt.Sprintf(":%d", cfg.port)
	url := fmt.Sprintf("http://localhost:%d/", cfg.port)

	fmt.Printf("sgo: static file server written in Go\n")
	fmt.Printf("  DIR: %s\n", root)
	fmt.Printf("  URL:       %s\n", url)
	fmt.Printf("  Press Ctrl+C to stop\n")

	fs := http.FileServer(http.Dir(root))
	if err := http.ListenAndServe(addr, fs); err != nil {
		fmt.Fprintf(os.Stderr, "listen on %s: %v\n", addr, err)
		os.Exit(1)
	}
}
