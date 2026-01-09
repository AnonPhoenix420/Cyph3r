```

  _____    __     __   ______    _    _    _____   ______
 / ____|  \ \   / /  | ___ \  | |  | |  |___  |  | ___ \
| |       \ \_/ /   | |_/ /  | |__| |    / /   | |_/ /
| |        \   /    |  __/   |  __  |  |_ \   |  _  \
| |____     | |     | |      | |  | |  ___) |  | | \ \
 \_____|    |_|     \_|      |_|  |_|  |____/   \_|  \_|

```

## ⚠️ Disclaimer

**For Educational and Professional Services Use Only**
**The creator is NOT responsible for misuse of this tool**

Always ensure you have **explicit permission** before testing any network, host, or service you do not own.

---

## 🧠 Overview

`CYPH3R` is a modular, Go-based network diagnostics and monitoring utility designed for learning, troubleshooting, and professional network validation. It supports **continuous monitoring**, **downtime tracking**, and multiple protocols while remaining lightweight and Termux-friendly.

Runs cleanly on:

* Parrot OS
* Linux (x86_64 & ARM64)
* Termux (Android / aarch64)

---

## ✨ Features

* TCP connectivity testing
* UDP probing (local echo friendly)
* HTTP / HTTPS status & latency checks
* Continuous monitor mode (tracks downtime)
* GeoIP / DNS resolution (IP or localhost)
* Optional phone number metadata lookup
* Colored terminal output (UP = blue, DOWN = red)
* JSON output mode for scripting & automation
* Version & author metadata flags
* Single static binary build support
* ARM64 / aarch64 compatible

---

## 📦 Requirements

* Go **1.23.0**
* git

### Supported Platforms

* Parrot OS (Linux)
* Linux PCs (Debian-based recommended)
* Termux on Android

---

## 🛠️ Installation

### 🔹 Parrot OS / Linux

```bash
curl -LO https://go.dev/dl/go1.23.0.linux-arm64.tar.gz
sudo tar -C /usr/local -xzf go1.23.0.linux-arm64.tar.gz
```
```
tar -xzf go1.23.0.linux-arm64.tar.gz
```
Verify Go:

```bash
go version
```

---

### 🔹 Termux (Android)

```bash
apt update && apt upgrade
```
```
apt install git
```

---

## 📥 Clone Repository

```bash
git clone https://github.com/AnonPhoenix420/cyph3r.git
```
```
cd cyph3r
```
```
git clone https://github.com/nyaruka/phonenumbers.git
```
---

## 🔧 Build & Install

### Local build (recommended for Termux)

```bash
go mod tidy
go build -o cyph3r ./cmd/cyph3r
```

### System-wide install (Linux only)

```bash
sudo install -m 755 cyph3r /usr/local/bin/cyph3r
```

---

## 🚀 Usage

### Basic TCP test

```bash
./cyph3r --target localhost --port 80 --proto tcp
```

### GeoIp Lookup Localhost/IP
```
./cyph3r --target example.com --geoip
```

### Continuous monitoring with downtime tracking

```bash
./cyph3r --target example.com --port 443 --proto https --monitor --interval 5
```

### UDP test

```bash
./cyph3r --target 127.0.0.1 --port 53 --proto udp
```

### HTTP / HTTPS test

```bash
./cyph3r --target example.com --port 443 --proto https
```

### JSON output mode

```bash
./cyph3r --target example.com --proto https --json
```

### Phone metadata lookup (optional)

```bash
./cyph3r --phone +14155552671
```

### How To Use Port Scanner 
```

./cyph3r --target 127.0.0.1 --portscan

```
### Change To Custom Port Range
```
./cyph3r --target 192.168.1.10 --portscan --scanstart 1 --scanend 10000

```

### Help shown on 
```
-h

```
```
--help

```
or no arg

### Version info

```bash

./cyph3r --version

```

---

## 🧪 Monitor Mode Explained

When `--monitor` is enabled, CYPH3R:

* Continues running even if the target goes DOWN
* Tracks how long the target remains unavailable
* Reports downtime when the target comes back UP
* Runs indefinitely until interrupted (`Ctrl+C`)

---

## 📁 Project Structure

```
cyph3r/
├── cmd/cyph3r/main.go      # Entry point
├── internal/
│   ├── netcheck/           # TCP / UDP / HTTP logic
│   ├── geo/                # GeoIP & DNS lookup
│   ├── phone/              # Phone metadata lookup
│   ├── output/             # Color & JSON output
│   └── version/            # Version metadata
├── install.sh
├── go.mod
├── LICENSE
└── README.md
```

---

## 📸 Example Output

```
[UP] Target is UP
[DOWN] Target went DOWN
[UP] Target is UP again (downtime: 1m42s)
```

---

## ⚖️ License

```
Educational & Professional Use Only

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND.
THE CREATOR IS NOT RESPONSIBLE FOR MISUSE OR DAMAGE CAUSED BY THIS TOOL.
```

---

## 🧠 Notes

* No `--i-own-this` flag required
* No capped duration limits
* Designed for diagnostics, learning, and professional use
* **Do not** use this tool for unauthorized testing

---

Happy hacking — responsibly 🧠🚀
