package main

import (
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func mustIPNet(t *testing.T, cidr string) net.Addr {
	t.Helper()
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		t.Fatal(err)
	}
	return ipnet
}

func TestFindLANIPv4FromAddrs(t *testing.T) {
	tests := []struct {
		name  string
		addrs []net.Addr
		want  string
	}{
		{
			name: "prefers private over loopback and public",
			addrs: []net.Addr{
				mustIPNet(t, "127.0.0.1/32"),
				mustIPNet(t, "203.0.113.1/32"),
				mustIPNet(t, "192.168.1.10/32"),
			},
			want: "192.168.1.10",
		},
		{
			name: "loopback only",
			addrs: []net.Addr{
				mustIPNet(t, "127.0.0.1/8"),
			},
			want: "",
		},
		{
			name: "ipv6 only",
			addrs: []net.Addr{
				mustIPNet(t, "fe80::1/64"),
			},
			want: "",
		},
		{
			name: "private over link-local",
			addrs: []net.Addr{
				mustIPNet(t, "169.254.12.34/32"),
				mustIPNet(t, "10.0.0.5/32"),
			},
			want: "10.0.0.5",
		},
		{
			name: "public fallback",
			addrs: []net.Addr{
				mustIPNet(t, "203.0.113.9/32"),
			},
			want: "203.0.113.9",
		},
		{
			name:  "empty",
			addrs: nil,
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findLANIPv4FromAddrs(tt.addrs)
			if tt.want == "" {
				if got != nil {
					t.Fatalf("got %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("got nil, want IP")
			}
			if got.String() != tt.want {
				t.Fatalf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestFormatLANURL(t *testing.T) {
	got := formatLANURL(net.ParseIP("192.168.1.42"), 5678)
	want := "http://192.168.1.42:5678/"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestResolveRootDir(t *testing.T) {
	base := t.TempDir()
	withSpace := filepath.Join(base, "My Project")
	if err := os.Mkdir(withSpace, 0o755); err != nil {
		t.Fatal(err)
	}
	cafe := filepath.Join(base, "Café")
	if err := os.Mkdir(cafe, 0o755); err != nil {
		t.Fatal(err)
	}
	missing := filepath.Join(base, "nope")
	filePath := filepath.Join(base, "file.txt")
	if err := os.WriteFile(filePath, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		dir     string
		want    string
		wantErr bool
	}{
		{"dot under cwd", withSpace, withSpace, false},
		{"absolute with space", withSpace, withSpace, false},
		{"utf8", cafe, cafe, false},
		{"missing", missing, "", true},
		{"file not dir", filePath, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.dir
			if tt.name == "dot under cwd" {
				if err := os.Chdir(withSpace); err != nil {
					t.Fatal(err)
				}
				dir = "."
				tt.want = withSpace
			}
			got, err := resolveRootDir(dir)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			wantAbs, err := filepath.Abs(tt.want)
			if err != nil {
				t.Fatal(err)
			}
			if got != wantAbs {
				t.Fatalf("got %q, want %q", got, wantAbs)
			}
		})
	}
}

func TestExecutableDir(t *testing.T) {
	got, err := executableDir()
	if err != nil {
		t.Fatal(err)
	}
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Dir(exe)
	wantAbs, err := filepath.Abs(want)
	if err != nil {
		t.Fatal(err)
	}
	gotAbs, err := filepath.Abs(got)
	if err != nil {
		t.Fatal(err)
	}
	if gotAbs != wantAbs {
		t.Fatalf("got %q, want %q", gotAbs, wantAbs)
	}
}

func TestParseConfig(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    config
		wantErr bool
	}{
		{"defaults", nil, config{port: 5678}, false},
		{"positional port", []string{"8080"}, config{port: 8080}, false},
		{"rejected -port flag", []string{"-port", "3000"}, config{}, true},
		{"unknown -dir flag", []string{"-dir", "/tmp"}, config{}, true},
		{"positional path rejected", []string{"/tmp/foo"}, config{}, true},
		{"invalid port positional", []string{"99999"}, config{}, true},
		{"unquoted path splits into multiple args", []string{"/tmp/My", "Project", "site"}, config{}, true},
		{"port with extra args", []string{"8080", "extra"}, config{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseConfig(tt.args)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Fatalf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func startTestServer(t *testing.T, root string) *httptest.Server {
	t.Helper()
	abs, err := resolveRootDir(root)
	if err != nil {
		t.Fatal(err)
	}
	return httptest.NewServer(newFileServer(abs))
}

func getStatusBody(t *testing.T, rawURL string) (int, string) {
	t.Helper()
	resp, err := http.Get(rawURL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	return resp.StatusCode, string(body)
}

func urlPath(segments ...string) string {
	if len(segments) == 0 {
		return "/"
	}
	parts := make([]string, len(segments))
	for i, seg := range segments {
		parts[i] = url.PathEscape(seg)
	}
	return "/" + path.Join(parts...)
}

func TestRuntimeFilesystemChanges(t *testing.T) {
	t.Run("runtime file change", func(t *testing.T) {
		root := t.TempDir()
		hello := filepath.Join(root, "hello.txt")
		if err := os.WriteFile(hello, []byte("v1"), 0o644); err != nil {
			t.Fatal(err)
		}

		srv := startTestServer(t, root)
		defer srv.Close()

		u := srv.URL + "/hello.txt"
		code, body := getStatusBody(t, u)
		if code != http.StatusOK || body != "v1" {
			t.Fatalf("initial: status=%d body=%q", code, body)
		}

		if err := os.WriteFile(hello, []byte("v2"), 0o644); err != nil {
			t.Fatal(err)
		}
		code, body = getStatusBody(t, u)
		if code != http.StatusOK || body != "v2" {
			t.Fatalf("after overwrite: status=%d body=%q", code, body)
		}

		if err := os.Remove(hello); err != nil {
			t.Fatal(err)
		}
		code, _ = getStatusBody(t, u)
		if code != http.StatusNotFound {
			t.Fatalf("after delete: status=%d, want 404", code)
		}
	})

	t.Run("runtime content change", func(t *testing.T) {
		root := t.TempDir()
		page := filepath.Join(root, "index.html")
		if err := os.WriteFile(page, []byte("<p>first</p>"), 0o644); err != nil {
			t.Fatal(err)
		}

		srv := startTestServer(t, root)
		defer srv.Close()
		u := srv.URL + "/index.html"

		changes := []struct {
			name string
			body []byte
		}{
			{"initial", []byte("<p>first</p>")},
			{"grow", []byte("<p>first</p><p>added section</p>")},
			{"shrink", []byte("ok")},
			{"utf8", []byte("<title>Café — demo</title>")},
		}

		for i, ch := range changes {
			if i > 0 {
				if err := os.WriteFile(page, ch.body, 0o644); err != nil {
					t.Fatal(err)
				}
			}
			code, got := getStatusBody(t, u)
			if code != http.StatusOK {
				t.Fatalf("%s: status=%d, want 200", ch.name, code)
			}
			if got != string(ch.body) {
				t.Fatalf("%s: body=%q, want %q", ch.name, got, string(ch.body))
			}
		}
	})

	t.Run("runtime folder change", func(t *testing.T) {
		root := t.TempDir()
		if err := os.WriteFile(filepath.Join(root, "index.txt"), []byte("root"), 0o644); err != nil {
			t.Fatal(err)
		}

		srv := startTestServer(t, root)
		defer srv.Close()

		tests := []struct {
			name string
			run  func(t *testing.T)
		}{
			{
				name: "add folder with space",
				run: func(t *testing.T) {
					dir := filepath.Join(root, "New Folder")
					if err := os.Mkdir(dir, 0o755); err != nil {
						t.Fatal(err)
					}
					if err := os.WriteFile(filepath.Join(dir, "page.html"), []byte("new"), 0o644); err != nil {
						t.Fatal(err)
					}
					u := srv.URL + urlPath("New Folder", "page.html")
					code, body := getStatusBody(t, u)
					if code != http.StatusOK || body != "new" {
						t.Fatalf("status=%d body=%q", code, body)
					}
				},
			},
			{
				name: "rename folder",
				run: func(t *testing.T) {
					oldDir := filepath.Join(root, "old")
					if err := os.Mkdir(oldDir, 0o755); err != nil {
						t.Fatal(err)
					}
					if err := os.WriteFile(filepath.Join(oldDir, "x.txt"), []byte("x"), 0o644); err != nil {
						t.Fatal(err)
					}
					oldURL := srv.URL + urlPath("old", "x.txt")
					code, body := getStatusBody(t, oldURL)
					if code != http.StatusOK || body != "x" {
						t.Fatalf("before rename: status=%d body=%q", code, body)
					}

					renamed := filepath.Join(root, "renamed")
					if err := os.Rename(oldDir, renamed); err != nil {
						t.Fatal(err)
					}
					code, _ = getStatusBody(t, oldURL)
					if code != http.StatusNotFound {
						t.Fatalf("old URL after rename: status=%d, want 404", code)
					}
					newURL := srv.URL + urlPath("renamed", "x.txt")
					code, body = getStatusBody(t, newURL)
					if code != http.StatusOK || body != "x" {
						t.Fatalf("new URL after rename: status=%d body=%q", code, body)
					}
				},
			},
			{
				name: "delete folder",
				run: func(t *testing.T) {
					gone := filepath.Join(root, "gone")
					if err := os.Mkdir(gone, 0o755); err != nil {
						t.Fatal(err)
					}
					if err := os.WriteFile(filepath.Join(gone, "y.txt"), []byte("y"), 0o644); err != nil {
						t.Fatal(err)
					}
					u := srv.URL + urlPath("gone", "y.txt")
					code, body := getStatusBody(t, u)
					if code != http.StatusOK || body != "y" {
						t.Fatalf("before delete: status=%d body=%q", code, body)
					}
					if err := os.RemoveAll(gone); err != nil {
						t.Fatal(err)
					}
					code, _ = getStatusBody(t, u)
					if code != http.StatusNotFound {
						t.Fatalf("after delete: status=%d, want 404", code)
					}
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tt.run(t)
			})
		}
	})

	t.Run("runtime subfolder change", func(t *testing.T) {
		root := t.TempDir()
		if err := os.WriteFile(filepath.Join(root, "marker.txt"), []byte("up"), 0o644); err != nil {
			t.Fatal(err)
		}

		srv := startTestServer(t, root)
		defer srv.Close()

		nested := filepath.Join(root, "My Project", "sub", "deep")
		dataFile := filepath.Join(nested, "data.txt")
		u := srv.URL + urlPath("My Project", "sub", "deep", "data.txt")

		if err := os.MkdirAll(nested, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(dataFile, []byte("v1"), 0o644); err != nil {
			t.Fatal(err)
		}
		code, body := getStatusBody(t, u)
		if code != http.StatusOK || body != "v1" {
			t.Fatalf("create nested: status=%d body=%q", code, body)
		}

		if err := os.WriteFile(dataFile, []byte("v2"), 0o644); err != nil {
			t.Fatal(err)
		}
		code, body = getStatusBody(t, u)
		if code != http.StatusOK || body != "v2" {
			t.Fatalf("modify nested: status=%d body=%q", code, body)
		}

		if err := os.RemoveAll(filepath.Join(root, "My Project")); err != nil {
			t.Fatal(err)
		}
		code, _ = getStatusBody(t, u)
		if code != http.StatusNotFound {
			t.Fatalf("remove nested: status=%d, want 404", code)
		}
	})
}
