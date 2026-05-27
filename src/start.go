package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type config struct {
	port        int
	dir         string
	dirExplicit bool
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
		} else if !dirExplicit {
			dir = pos[0]
			dirExplicit = true
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

	return config{port: port, dir: dir, dirExplicit: dirExplicit}, nil
}

func executableDir() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("locate executable: %w", err)
	}
	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		return "", fmt.Errorf("resolve executable: %w", err)
	}
	return filepath.Dir(exe), nil
}

func resolveServePath(cfg config) (string, error) {
	if cfg.dirExplicit {
		return cfg.dir, nil
	}
	return executableDir()
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

func newFileServer(root string) http.Handler {
	return http.FileServer(http.Dir(root))
}

func isPrivateIPv4(ip net.IP) bool {
	ip4 := ip.To4()
	if ip4 == nil {
		return false
	}
	switch {
	case ip4[0] == 10:
		return true
	case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
		return true
	case ip4[0] == 192 && ip4[1] == 168:
		return true
	}
	return false
}

func isLinkLocalIPv4(ip net.IP) bool {
	ip4 := ip.To4()
	return ip4 != nil && ip4[0] == 169 && ip4[1] == 254
}

func lanIPv4Score(ip net.IP) int {
	ip4 := ip.To4()
	if ip4 == nil || ip.IsLoopback() {
		return -1
	}
	if isPrivateIPv4(ip4) {
		return 3
	}
	if isLinkLocalIPv4(ip4) {
		return 1
	}
	return 2
}

func findLANIPv4FromAddrs(addrs []net.Addr) net.IP {
	var best net.IP
	bestScore := -1
	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		ip := ipnet.IP
		score := lanIPv4Score(ip)
		if score > bestScore {
			bestScore = score
			best = ip.To4()
		}
	}
	if bestScore < 0 {
		return nil
	}
	return best
}

func findLANIPv4() net.IP {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}
	return findLANIPv4FromAddrs(addrs)
}

func formatLANURL(ip net.IP, port int) string {
	return fmt.Sprintf("http://%s:%d/", ip, port)
}

func main() {
	cfg, err := parseConfig(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	servePath, err := resolveServePath(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	root, err := resolveRootDir(servePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	addr := fmt.Sprintf(":%d", cfg.port)
	url := fmt.Sprintf("http://127.0.0.1:%d/", cfg.port)

	fmt.Printf("sgo: static file server written in Go\n")
	fmt.Printf("  DIR: %s\n", root)
	fmt.Printf("  URL: %s\n", url)
	if ip := findLANIPv4(); ip != nil {
		fmt.Printf("  LAN: %s\n", formatLANURL(ip, cfg.port))
	}
	fmt.Printf("  Press Ctrl+C to stop\n")

	if err := http.ListenAndServe(addr, newFileServer(root)); err != nil {
		fmt.Fprintf(os.Stderr, "listen on %s: %v\n", addr, err)
		os.Exit(1)
	}
}
