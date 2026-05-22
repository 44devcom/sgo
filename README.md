# sgo

A minimal static file server written in Go. It serves the current working directory over HTTP.

## Usage

Run from the directory you want to serve:

```bash
./sgo
```

### Port

The default port is **5678**. Pass a port as the first argument:

```bash
./sgo 8080
```

Then open `http://localhost:<port>/` in your browser.

## Download

| Platform | Path | Download |
|----------|------|----------|
| Linux (amd64) | `dist/linux-amd64/sgo` | [download](https://github.com/44devcom/sgo/raw/refs/heads/master/dist/linux-amd64/sgo) |
| macOS (Intel) | `dist/darwin-amd64/sgo` | [download](https://github.com/44devcom/sgo/raw/refs/heads/master/dist/darwin-amd64/sgo) |
| macOS (Apple Silicon) | `dist/darwin-arm64/sgo` | [download](https://github.com/44devcom/sgo/raw/refs/heads/master/dist/darwin-arm64/sgo) |
| Windows (amd64) | `dist/windows-amd64/sgo.exe` | [download](https://github.com/44devcom/sgo/raw/refs/heads/master/dist/windows-amd64/sgo.exe) |

