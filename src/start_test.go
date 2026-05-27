package main

import (
	"os"
	"path/filepath"
	"testing"
)

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

func TestParseConfig(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    config
		wantErr bool
	}{
		{"defaults", nil, config{port: 5678, dir: "."}, false},
		{"positional port", []string{"8080"}, config{port: 8080, dir: "."}, false},
		{"flag port", []string{"-port", "3000"}, config{port: 3000, dir: "."}, false},
		{"flag dir", []string{"-dir", "/tmp/foo"}, config{port: 5678, dir: "/tmp/foo"}, false},
		{"dir equals form", []string{"-dir=/tmp/My Project"}, config{port: 5678, dir: "/tmp/My Project"}, false},
		{"port and dir flags", []string{"-port", "9000", "-dir", "/srv/www"}, config{port: 9000, dir: "/srv/www"}, false},
		{"invalid port positional", []string{"99999"}, config{}, true},
		{"positional path without dir flag", []string{"/tmp/My", "Project", "site"}, config{}, true},
		{"unexpected with explicit dir", []string{"-dir", "/tmp", "extra"}, config{}, true},
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
