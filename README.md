# sgo

**sgo** is a tiny static file server written in Go — one binary, no config file, no dependencies at runtime. Point it at a folder and open the URL in a browser.

### Features

- **Single executable** — pure Go, statically linked (`CGO_ENABLED=0`)
- **Serves a directory over HTTP** — built on `net/http` file server
- **Default port 5678** — change with `-port` or a positional port (`./sgo 8080`)
- **Smart default root** — serves the folder that contains the `sgo` binary (not the shell’s current working directory)
- **Change root directory** — serve any path (required on Termux, when `sgo` is on `PATH`, or when the site lives elsewhere)
- **Listens on all interfaces** — use `127.0.0.1` locally or `<host-ip>` from phones and other machines on the same network
- **Startup summary** — prints resolved `DIR:` and `URL:` so you can confirm the right folder
- **Cross-platform builds** — Linux (amd64, arm64), macOS (Intel & Apple Silicon), Windows (amd64)

---

## Usage

| Input | Effect |
|-------|--------|
| `./sgo` | Port **5678**, serve directory containing the binary |
| `./sgo 8080` | Port **8080** (positional, backward compatible) |
| `./sgo -port 8080` | Port **8080** |
| `./sgo -dir=/path/to/site` | Serve that directory |
| `./sgo -port 8080 -dir="/path/with spaces"` | Custom port and root |

Open `http://127.0.0.1:<port>/` in a browser. The server listens on all interfaces, so other devices on the same network can use `http://<host-ip>:<port>/`.

Check the **`DIR:`** line at startup. If a path with spaces was split (e.g. you see `.../My` instead of `.../My Project/site`), quote the path or use `-dir=`:

```bash
./sgo -dir="/path/with spaces"
```

On macOS and in most shells, prefer **`-dir=/full/path`** (equals form) so the path is not split before sgo runs.

---

## Install

Install scripts download the binary for your OS and CPU, then place it in your Downloads folder.

### Linux & macOS

Detects `x86_64` → amd64 and `aarch64` → arm64. Saves `~/Downloads/sgo`.

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
.\sgo.exe -port 8080
```

**Termux (Android)** — use the **Linux & macOS** installer inside Termux (same script; Termux is Linux + aarch64 → `linux-arm64` binary):

```bash
curl -fsSL https://github.com/44devcom/sgo/raw/refs/heads/master/bin/install.sh | bash
termux-setup-storage   # once, for ~/storage/* access
```

Always pass **`-dir`** to a path under `~/storage/...` (see [Termux (Android)](#termux-android) below).

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

**Download:** [linux-amd64](https://github.com/44devcom/sgo/raw/refs/heads/master/dist/linux-amd64/sgo) · [linux-arm64](https://github.com/44devcom/sgo/raw/refs/heads/master/dist/linux-arm64/sgo)

```bash
chmod +x sgo
./sgo
./sgo 8080
./sgo -port 8080 -dir=/var/www/myproject
```

Place `sgo` next to your site files, or use `-dir`. Open `http://127.0.0.1:5678/` (or your chosen port). From another device: `http://<machine-ip>:5678/`.

---

### macOS

**Download:** [darwin-amd64](https://github.com/44devcom/sgo/raw/refs/heads/master/dist/darwin-amd64/sgo) · [darwin-arm64](https://github.com/44devcom/sgo/raw/refs/heads/master/dist/darwin-arm64/sgo)

```bash
cd ~/Downloads   # or the folder containing sgo
chmod +x ./sgo
xattr -d com.apple.quarantine ./sgo 2>/dev/null || true
./sgo
```

The installer runs the quarantine step for you. When launching from Finder, the working directory may be your home folder — sgo still serves the **executable’s directory** unless you set `-dir`:

```bash
./sgo -dir="/Users/me/My Project/site"
```

Prefer **`-dir=/full/path`** (equals form) so spaces are not split by the shell.

---

### Windows

**Download:** [windows-amd64/sgo.exe](https://github.com/44devcom/sgo/raw/refs/heads/master/dist/windows-amd64/sgo.exe)

```powershell
cd $env:USERPROFILE\Downloads
.\sgo.exe
.\sgo.exe -port 8080
.\sgo.exe -port 8080 -dir "C:\Users\me\My Project\site"
```

Open `http://127.0.0.1:5678/` in a browser. Other devices on the LAN: `http://<pc-ip>:5678/`.

---

### Termux (Android)

**Download:** [linux-arm64/sgo](https://github.com/44devcom/sgo/raw/refs/heads/master/dist/linux-arm64/sgo) (or use the [installer](#install) above).

1. `termux-setup-storage` (once)
2. Put your site under shared storage, e.g. `~/storage/downloads/my-site/`
3. Run with **`-dir`**:

```bash
chmod +x ~/Downloads/sgo
~/Downloads/sgo -dir="$HOME/storage/downloads/my-site" -port 5678
```

| Termux path | Typical location |
|-------------|------------------|
| `~/storage/downloads` | Download |
| `~/storage/shared` | Internal storage root |
| `~/storage/documents` | Documents |

- On the phone: `http://127.0.0.1:5678/`
- From another device on Wi‑Fi: `http://<phone-ip>:5678/`

---

## Development

Requires Go 1.24+ (see `go.mod`).

```bash
go test ./...
./bin/build.sh
```

`build.sh` runs tests first, then cross-compiles all `dist/` targets with `CGO_ENABLED=0`. Builds are skipped if tests fail.
