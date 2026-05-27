package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type config struct {
	port int
}

func parseConfig(args []string) (config, error) {
	const defaultPort = 5678

	switch len(args) {
	case 0:
		return config{port: defaultPort}, nil
	case 1:
		p, err := strconv.Atoi(args[0])
		if err != nil {
			return config{}, fmt.Errorf("unexpected arguments: %s", args[0])
		}
		if p < 1 || p > 65535 {
			return config{}, fmt.Errorf("invalid port: %d", p)
		}
		return config{port: p}, nil
	default:
		return config{}, fmt.Errorf("unexpected arguments: %s", strings.Join(args, " "))
	}
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

	servePath, err := executableDir()
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
