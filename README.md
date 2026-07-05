<h1 align="center">🏮 Lantern</h1>

<p align="center">
  <b>A lightweight LAN music streaming server written in Go.</b>
</p>

<p align="center">
Stream your personal music library from one machine to any device on your network — phone, laptop, TV, or a friend’s browser — instantly, with zero setup headaches.
</p>

## Problem

Lantern comes from personal ick. Modern music streaming services are often inconvenient - subscriptions, licensing, region locks, internet connection dependency, etc.

So for some people it is more appealing to collect and store their favorite albums locally as files.

**Here comes the problem - your library lives only on one device.**

If you want to browse your lib on any other device you simply have to sync them - that means double or triple the storage.

What Lantern does is it turns your collection into a personal streaming server accessible anywhere on LAN.

## Quick Start

### 1. Clone the repository

```bash
git clone https://github.com/<your-username>/lantern.git
cd lantern
````

### 2. Build

```bash
go build -o lantern
```

*You can specify GOOS and GOARCH to build it for different devices.*

### 3. Run

```bash
./lantern
```

## Configuration

Lantern uses a simple config file. You will be prompted to specify some values on first launch. After that it will be accessible in:

```
~/.config/lantern/config.json
```

Example:

```json
{
  "music_dir": "/home/user/Music",
  "port": 8080
}
```

## Usage

Once running, Lantern will generate a connection QR-code as well as URL.  
You can access it on host device via:

```
http://localhost:8080
```

Or from another device on the same network:

```
http://<your-local-ip>:8080
```

Use the built-in QR code to open it instantly on mobile.

## Roadmap

### Storage & Performance

- [x] SQLite-based persistent library index
- [ ] File system watcher (fsnotify) for live updates

### Features

- [ ] Playlists
- [ ] Playback caching
- [ ] Better search & filtering

### Misc

- [ ] Vanilla CSS only

### Platform

- [ ] Responsive UI for comfortable mobile experience

## Features

  - Local music streaming over HTTP (LAN)
  - Minimal server-rendered web UI
  - Recursive folder scanning
  - QR code sharing for instant access on mobile devices
  - LAN-open access (no authentication by design)
  - Supports common audio formats:
      - MP3
      - FLAC
      - WAV
      - OGG
  - Simple JSON configuration
  
## Tech stack
- Go
- go-taglib (metadata extraction)
- modernc/sqlite (SQLite driver for local persistence)
- Tailwind CSS *(will be replaced with vanilla css in future)*
- Vanilla server-rendered HTML

## License

Project is MIT licensed.
