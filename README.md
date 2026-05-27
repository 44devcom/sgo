# sgo

A minimal static file server written in Go. It serves files from a directory over HTTP.

## Usage

Serve the directory containing the `sgo` binary (default port **5678**):

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

By default, sgo serves the **directory that contains the `sgo` executable** (not the shell’s current working directory). Place `sgo` next to your site files, or use `-dir` to serve another folder:

```bash
./sgo -dir=/Users/me/My Project/site
./sgo -dir "/Users/me/My Project/site"
./sgo -port 8080 -dir="/Users/me/Café/demo"
```

On macOS and in most shells, prefer **`-dir=/full/path`** (equals form) so the path is not split at spaces before sgo runs.

Use `-dir` when `sgo` is on your `PATH` (for example `/usr/local/bin/sgo`) but the site lives elsewhere, or when you want a folder other than the executable’s directory.

For both a custom directory and port, use flags:

```bash
./sgo -port 8080 -dir="/path/with spaces"
```

Startup prints the **resolved absolute path** as `DIR:` — check that line if the wrong folder is served.

### Troubleshooting paths with spaces

If `DIR:` shows the parent of a folder whose name contains a space (e.g. you expected `.../My Project/site` but see `.../My`), the shell split the path. Fix it with quotes or `-dir=`:

```bash
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

Put `sgo` in the folder you want to serve (for example your project or Downloads). Open Terminal (Cmd + Space, type **Terminal**, Enter):

```bash
cd ~/Downloads
chmod +x ./sgo
xattr -d com.apple.quarantine ./sgo 2>/dev/null || true
```

When you launch `sgo` from Finder, macOS sets the working directory to your home folder; sgo still serves the folder where the binary lives. Use `-dir=` when the site is not next to the binary:

```bash
./sgo -dir="/Users/me/My Project/site"
```

Use `-dir=` or quotes so spaces in the path are preserved.

### Installer

```bash
curl -fsSL https://https://github.com/44devcom/sgo/raw/refs/heads/master/bin/install.sh | bash
```