# sgo

A minimal static file server written in Go. It serves files from a directory over HTTP.

## Usage

Serve the current working directory (default port **5678**):

```bash
./sgo
```

Pass a port as the first argument (backward compatible):

```bash
./sgo 8080
```

Or use flags:

```bash
./sgo -port 8080
```

Then open `http://localhost:<port>/` in your browser.

### Serve directory

By default, sgo serves the directory you run it from. To serve another folder (recommended when the path contains **spaces** or **non-ASCII** characters), use `-dir`:

```bash
./sgo -dir=/Users/me/My Project/site
./sgo -dir "/Users/me/My Project/site"
./sgo -port 8080 -dir="/Users/me/Café/demo"
```

On macOS and in most shells, prefer **`-dir=/full/path`** (equals form) so the path is not split at spaces before sgo runs.

Without `-dir`, sgo always serves the **current working directory** (where you run the command). Use `cd` with a quoted path, or pass `-dir`.

For both a custom directory and port, use flags:

```bash
./sgo -port 8080 -dir="/path/with spaces"
```

Startup prints the **resolved absolute path** as `directory:` — check that line if the wrong folder is served.

### Troubleshooting paths with spaces

If `directory:` shows the parent of a folder whose name contains a space (e.g. you expected `.../My Project/site` but see `.../My`), the shell split the path. Fix it with quotes or `-dir=`:

```bash
cd "/Users/me/My Project/site"   # quoted cd
./sgo -dir="/Users/me/My Project/site"
```

## Download

| Platform | Path | Download |
|----------|------|----------|
| Linux (amd64) | `dist/linux-amd64/sgo` | [download](https://github.com/44devcom/sgo/raw/refs/heads/master/dist/linux-amd64/sgo) |
| macOS (Intel) | `dist/darwin-amd64/sgo` | [download](https://github.com/44devcom/sgo/raw/refs/heads/master/dist/darwin-amd64/sgo) |
| macOS (Apple Silicon) | `dist/darwin-arm64/sgo` | [download](https://github.com/44devcom/sgo/raw/refs/heads/master/dist/darwin-arm64/sgo) |
| Windows (amd64) | `dist/windows-amd64/sgo.exe` | [download](https://github.com/44devcom/sgo/raw/refs/heads/master/dist/windows-amd64/sgo.exe) |

## macOS

Open Terminal (Cmd + Space, type **Terminal**, Enter). From your Downloads folder (or wherever you placed `sgo`):

```bash
chmod +x ./sgo
xattr -d com.apple.quarantine ./sgo 2>/dev/null || true
./sgo -dir="$HOME/My Project/site"
```

Replace `My Project/site` with your folder. Use `-dir=` or quotes so spaces in the path are preserved.
