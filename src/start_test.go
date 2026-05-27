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

func TestResolveServePath(t *testing.T) {
	t.Run("default uses executable dir", func(t *testing.T) {
		got, err := resolveServePath(config{dir: ".", dirExplicit: false})
		if err != nil {
			t.Fatal(err)
		}
		want, err := executableDir()
		if err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})

	t.Run("explicit dir", func(t *testing.T) {
		got, err := resolveServePath(config{dir: "/tmp/foo", dirExplicit: true})
		if err != nil {
			t.Fatal(err)
		}
		if got != "/tmp/foo" {
			t.Fatalf("got %q, want /tmp/foo", got)
		}
	})
}

func TestParseConfig(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    config
		wantErr bool
	}{
		{"defaults", nil, config{port: 5678, dir: ".", dirExplicit: false}, false},
		{"positional port", []string{"8080"}, config{port: 8080, dir: ".", dirExplicit: false}, false},
		{"flag port", []string{"-port", "3000"}, config{port: 3000, dir: ".", dirExplicit: false}, false},
		{"flag dir", []string{"-dir", "/tmp/foo"}, config{port: 5678, dir: "/tmp/foo", dirExplicit: true}, false},
		{"dir equals form", []string{"-dir=/tmp/My Project"}, config{port: 5678, dir: "/tmp/My Project", dirExplicit: true}, false},
		{"port and dir flags", []string{"-port", "9000", "-dir", "/srv/www"}, config{port: 9000, dir: "/srv/www", dirExplicit: true}, false},
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
