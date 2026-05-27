# sgo

**sgo** is a tiny static file server written in Go — one binary, no config file, no dependencies at runtime. Point it at a folder and open the URL in a browser.

### Features

- **Single executable** — pure Go, statically linked (`CGO_ENABLED=0`)
- **Serves a directory over HTTP** — built on `net/http` file server
- **Default port 5678** — change with a positional port (`./sgo 8080`)
- **Smart default root** — serves the folder that contains the `sgo` binary (not the shell’s current working directory)
- **Listens on all interfaces** — use `127.0.0.1` locally or `<host-ip>` from phones and other machines on the same network
- **Startup summary** — prints resolved `DIR:`, `URL:`, and `LAN:` (when a LAN IPv4 is found) so you can confirm the right folder and copy a network URL
- **Cross-platform builds** — Linux (amd64, arm64), macOS (Intel & Apple Silicon), Windows (amd64)

---

## Usage

| Input | Effect |
|-------|--------|
| `./sgo` | Port **5678**, serve directory containing the binary |
| `./sgo 8080` | Port **8080** |

Place `sgo` in the folder you want to serve (or copy the binary there). Open `http://127.0.0.1:<port>/` in a browser. The server listens on all interfaces; at startup, **`LAN:`** prints a ready-to-use URL for other devices on the same network when a suitable IPv4 is found (otherwise use `http://<host-ip>:<port>/` manually).

Check the **`DIR:`** line at startup to confirm which folder is being served.

Example startup output:

```
sgo: static file server written in Go
  DIR: /path/to/site
  URL: http://127.0.0.1:5678/
  LAN: http://192.168.1.42:5678/
  Press Ctrl+C to stop
```

`LAN:` is omitted when no suitable LAN IPv4 is found.

---

## Install

Install scripts download the binary for your OS and CPU, then place it in your Downloads folder.

### Linux & macOS

Detects `x86_64` → amd64 and `aarch64` → arm64 (Linux or macOS). Saves `~/Downloads/sgo`.

```bash
curl -fsSL https://github.com/44devcom/sgo/raw/refs/heads/master/bin/install.sh | bash
```

Then:

```bash
chmod +x ~/Downloads/sgo
~/Downloads/sgo
```

### Windows

Saves `%USERPROFILE%\Downloads\sgo.exe`.

```powershell
irm https://github.com/44devcom/sgo/raw/refs/heads/master/bin/install.ps1 | iex
```

If execution policy blocks `iex`:

```powershell
irm https://github.com/44devcom/sgo/raw/refs/heads/master/bin/install.ps1 -OutFile install.ps1
powershell -ExecutionPolicy Bypass -File .\install.ps1
```

Then:

```powershell
cd $env:USERPROFILE\Downloads
.\sgo.exe
.\sgo.exe 8080
```

---

## Download

Pre-built binaries in `dist/` (also linked per platform below):

| Platform | Artifact |
|----------|----------|
| Linux amd64 | `dist/linux-amd64/sgo` |
| Linux arm64 (aarch64) | `dist/linux-arm64/sgo` |
| macOS Intel | `dist/darwin-amd64/sgo` |
| macOS Apple Silicon | `dist/darwin-arm64/sgo` |
| Windows amd64 | `dist/windows-amd64/sgo.exe` |

---

### Linux

**Download:** [linux-amd64](https://github.com/44devcom/sgo/raw/refs/heads/master/dist/linux-amd64/sgo) · [linux-arm64](https://github.com/44devcom/sgo/raw/refs/heads/master/dist/linux-arm64/sgo) (native Linux ARM64, e.g. Raspberry Pi)

```bash
chmod +x sgo
./sgo
./sgo 8080
./sgo 8080
```

Place `sgo` next to your site files. Open `http://127.0.0.1:5678/` (or your chosen port). From another device on the LAN, use the **`LAN:`** line from startup, or `http://<machine-ip>:5678/`.

---

### macOS

**Download:** [darwin-amd64](https://github.com/44devcom/sgo/raw/refs/heads/master/dist/darwin-amd64/sgo) · [darwin-arm64](https://github.com/44devcom/sgo/raw/refs/heads/master/dist/darwin-arm64/sgo)

```bash
cd ~/Downloads   # or the folder containing sgo
chmod +x ./sgo
xattr -d com.apple.quarantine ./sgo 2>/dev/null || true
./sgo
```

The installer runs the quarantine step for you. When launching from Finder, the working directory may be your home folder — sgo still serves the **executable’s directory** (the folder that contains `sgo`), not Finder’s current folder. Keep `sgo` in the site folder you want to serve.

---

### Windows

**Download:** [windows-amd64/sgo.exe](https://github.com/44devcom/sgo/raw/refs/heads/master/dist/windows-amd64/sgo.exe)

```powershell
cd $env:USERPROFILE\Downloads
.\sgo.exe
.\sgo.exe 8080
```

Place `sgo.exe` in the folder you want to serve.

Open `http://127.0.0.1:5678/` in a browser. Other devices on the LAN: use the **`LAN:`** line from startup, or `http://<pc-ip>:5678/` if it is not shown.

---

## Development

### Prerequisites

- [Go 1.24+](https://go.dev/dl/) (see `go.mod` for the exact toolchain)
- Git clone of this repo

```bash
git clone https://github.com/44devcom/sgo.git
cd sgo
```

### Run from source

The entry point is `src/start.go` (single-file `main` package).

`sgo` always serves the **directory that contains the executable**, not the shell’s current working directory. With `go run`, Go builds a temporary binary (often under the module cache), so the served folder is **not** your project tree. For local development, build a binary in the folder you want to serve:

```bash
go build -o sgo src/start.go
./sgo
./sgo 8080
```

Open the **`URL:`** line printed at startup (default port **5678**).

### Build a local binary

```bash
go build -o sgo src/start.go
./sgo
```

Copy or move `sgo` next to your site files if the binary is not already there.

### Tests

```bash
go test ./...
```

Run a single package or verbose output:

```bash
go test -v ./src/...
```

### Release builds (`bin/build.sh`)

Cross-compile all `dist/` targets (tests run first; build is skipped if tests fail):

```bash
./bin/build.sh
```

`build.sh` sets `CGO_ENABLED=0` and writes:

| Output | `GOOS` / `GOARCH` | Notes |
|--------|-------------------|--------|
| `dist/linux-amd64/sgo` | linux / amd64 | |
| `dist/linux-arm64/sgo` | linux / arm64 | native Linux ARM64 |
| `dist/darwin-amd64/sgo` | darwin / amd64 | |
| `dist/darwin-arm64/sgo` | darwin / arm64 | |
| `dist/windows-amd64/sgo.exe` | windows / amd64 | |

To build one platform manually:

```bash
CGO_ENABLED=0 go build -o sgo src/start.go
```
