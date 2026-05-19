<h1 align="center">🏮 Lantern</h1>

<p align="center">
  <b>A lightweight LAN music streaming server written in Go.</b>
</p>

<p align="center">
Stream your personal music library from one machine to any device on your network — phone, laptop, TV, or a friend’s browser — instantly, with zero setup headaches.
</p>

---

## ✨ Why Lantern?

Modern music streaming services are convenient until they aren’t.

Your music disappears behind subscriptions, licensing, region locks, or simply lives on one device that never quite follows you around.

Lantern was built out of a simple idea:

> **Your music library should live with you — not in someone else’s cloud.**

If you already have a local collection (like a 9.5k+ track library), Lantern turns it into a personal streaming server accessible anywhere on your LAN.

No accounts. No cloud. No sync. Just your music.

---

## ⚡ Features

- 🎵 Stream local music over HTTP (LAN)
- 🌐 Minimal server-rendered web UI
- 📁 Recursive folder scanning (configurable music directory)
- 📷 QR code sharing for instant access on mobile devices
- 🔓 LAN-open access (no authentication by design)
- 🎧 Supports common audio formats:
  - MP3
  - FLAC
  - WAV
  - OGG
- ⚙️ Simple JSON configuration
- 🚀 Written in Go using `net/http`

---

## 🛠 Tech Stack

- Go (`net/http`)
- go-taglib (metadata extraction)
- Tailwind CSS (UI styling)
- Vanilla server-rendered HTML

---

## 🧠 Architecture

Lantern is intentionally simple:

```

+----------------------+
| Music Directory      |
| (filesystem)         |
+----------+-----------+
|
v
+----------------------+
| Scanner (runtime)    |
| recursive scan       |
+----------+-----------+
|
v
+----------------------+
| In-memory Library    |
| (current version)    |
+----------+-----------+
|
v
+----------------------+
| HTTP Server          |
| net/http (Go)        |
+----------+-----------+
|
v
+----------------------+
| Web UI (SSR HTML)    |
| + QR sharing         |
+----------------------+

````

---

## 🚀 Quick Start

### 1. Clone the repository

```bash
git clone https://github.com/<your-username>/lantern.git
cd lantern
````

### 2. Build

You can name the binary whatever you want:

```bash
go build -o lantern
```

### 3. Run

```bash
./lantern
```

---

## ⚙️ Configuration

Lantern uses a simple config file:

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

You can also use relative paths:

```json
{
  "music_dir": "./music",
  "port": 8080
}
```

---

## 🌍 Accessing Lantern

Once running, Lantern will generate a connection QR-code as well as URL:

```
http://localhost:8080
```

Or from another device on the same network:

```
http://<your-local-ip>:8080
```

Use the built-in QR code to open it instantly on mobile.

---

## 🧭 Roadmap

Lantern is still early-stage and evolving.

### Storage & Performance

- [ ] SQLite-based persistent library index
- [ ] File system watcher (fsnotify) for live updates
- [ ] Faster incremental rescans

### Features

- [ ] Playlists
- [ ] Playback caching
- [ ] Better search & filtering
- [ ] Album/artist metadata improvements (lazy loads mostly)

### Platform

- [ ] PWA mobile experience
- [ ] Optional internet access via reverse proxy

---

## 📦 Supported Audio Formats

| Format | Status |
| ------ | ------ |
| MP3    | ✅      |
| FLAC   | ✅      |
| WAV    | ✅      |
| OGG    | ✅      |

---

## 🧪 Future Ideas

- Lightweight authentication (optional toggle)
- Multi-user support
- Music rooms (stream music for everyone connected)

---

## ❓ FAQ

### Can I expose Lantern to the internet?

Yes, but it's designed for LAN use. Internet exposure is possible via port forwarding or reverse proxy, but authentication is not built-in yet.

### Does Lantern modify my music files?

No. Everything is read-only at the moment.

### Does it support streaming large libraries?

Yes — current version uses runtime scanning and in-memory indexing.

---

## 💭 Final Note

>Lantern is a personal project — built for the joy of turning a private music collection into something that feels alive across devices.
>
>It’s minimal by design, but intentionally leaves room to grow.

---

## 📜 License

Project is MIT licensed.
