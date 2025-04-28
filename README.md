# ⚡ HyperScanner + CORS (hxscanner) v1.4+
  ![Go Version](https://img.shields.io/badge/Go-1.17+-00ADD8?style=flat&logo=go)
![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)
![Platform](https://img.shields.io/badge/Platform-linux%20%7C%20macOS%20%7C%20windows-blue)
<p align="center">
  <img src="https://img.shields.io/badge/HyperScanner-purple?style=for-the-badge&logoColor=white" alt="HyperScanner" />
</p>

<p align="center">
  
  <pre style="color: purple;">
    ██╗  ██╗██╗  ██╗      ███████╗ ██████╗ █████╗ ███╗   ██╗███╗   ██╗███████╗██████╗
    ██║  ██║╚██╗██╔╝      ██╔════╝██╔════╝██╔══██╗████╗  ██║████╗  ██║██╔════╝██╔══██╗
    ███████║ ╚███╔╝█████╗ ███████╗██║     ███████║██╔██╗ ██║██╔██╗ ██║█████╗  ██████╔╝
    ██╔══██║ ██╔██╗╚════╝ ╚════██║██║     ██╔══██║██║╚██╗██║██║╚██╗██║██╔══╝  ██╔══██╗
    ██║  ██║██╔╝ ██╗      ███████║╚██████╗██║  ██║██║ ╚████║██║ ╚████║███████╗██║  ██║
    ╚═╝  ╚═╝╚═╝  ╚═╝      ╚══════╝ ╚═════╝╚═╝  ╚═╝╚═╝  ╚═══╝╚═╝  ╚═══╝╚══════╝╚═╝  ╚═╝
  </pre>

  <h3 align="center">HyperScanner v1.4+CORS (IP/Domain/URL Scanner)</h3>

</p>

---

**HyperScanner** is a powerful and futuristic HTTP status code scanner that takes a list of IPs or URLs and classifies their responses into structured folders and files.  
Designed for **speed**, **clarity**, and **beauty**, HyperScanner simplifies HTTP response analysis with an organized and efficient workflow.

---

## ✨ Features

- 🔍 **Scan IP & URLs Lists:** Quickly scans any list of IPs or URLs and fetches their HTTP status codes.
- 🗂️ **Organized Output:** Automatically creates folder structures based on HTTP status codes (1xx, 2xx, 3xx, 4xx, 5xx).
- 📁 **File-Based Response Storage:** Saves each response into categorized files, such as `200.txt` for HTTP 200 OK responses.
- ✅ **Detailed Logs:** Cleanly tracks results:
  - `ip_exist.txt`: List of valid and reachable IPs/URLs.
  - `ip_invalid.txt`: List of invalid or unreachable IPs/URLs.
  - `log.txt`: Comprehensive full scanning log.
- 🎨 **Enhanced CLI (Terminal Output):** Color-coded status codes for better readability (upcoming: icons + detailed categories).
- 🌐 **CORS Integration (New!):** Detects and logs CORS headers like `Access-Control-Allow-Origin`.
- 💻 **Cross-Platform:** Works flawlessly on **Windows**, **Linux**, and **macOS**.

---

## 🚀 Installation

Make sure you have **Go 1.19+** installed.

Then install HyperScanner using:

```bash
go install github.com/hx-corp/hxscanner@latest
```

---

## 🛠️ Usage

Prepare a text file (`ips.txt`) containing a list of IPs or URLs, one per line.

Run HyperScanner with:

```bash
hxscanner -i ips.txt
```

You can also specify an output directory:

```bash
hxscanner -i ips.txt -o my_results
```

---

## 📋 CLI Options

| Option        | Description |
| ------------- | ------------ |
| `-i <file>`   | Input file with targets (IPs/Domains/URLs), one per line (required if `-f` not used) |
| `-f <file>`   | Alias for `-i` |
| `-w <number>` | Number of concurrent scanning workers (default: number of CPU cores) |
| `-t <duration>` | HTTP request timeout (default: 5s) |
| `-q`          | Quiet mode: suppress individual results (except errors/warnings) |
| `--cors`      | Perform basic CORS vulnerability check on successful targets |
| `-h`          | Show this help message |

---

## 📂 Output Structure

HyperScanner organizes the output beautifully:

```plaintext
output/
├── 1xx/
│   └── 100.txt
├── 2xx/
│   ├── 200.txt
│   ├── 204.txt
│   └── ...
├── 3xx/
│   └── 301.txt
├── 4xx/
│   └── 404.txt
├── 5xx/
│   └── 500.txt
├── ip_exist.txt
├── ip_invalid.txt
├── log.txt
└── cors_detected.txt   (new in v1.4+)
```

- `<status_code>.txt`: IPs/URLs returning that status code.
- `ip_exist.txt`: Valid, reachable IPs/URLs.
- `ip_invalid.txt`: Failed or unreachable IPs/URLs.
- `log.txt`: Full detailed log of scanning activities.
- `cors_detected.txt`: IPs/URLs where CORS headers were found (`Access-Control-Allow-Origin`).

---

## 🌐 Cross-Platform Compatibility

HyperScanner runs seamlessly on:

- 🪟 Windows
- 🐧 Linux
- 🍎 macOS

No additional setup needed — just **Go** installed.

---

## 🖥️ GUI Version (Coming Soon)

We are building a beautiful, lightweight **GUI version** using **Wails**.  
Stay tuned for an enhanced desktop experience with all HyperScanner features!

---

## 🏗️ Contributing

We welcome contributions! Here's how you can help:

1. **Fork** the repository.
2. **Create** a new branch (`feature/your-feature-name`).
3. **Commit** your changes with clear messages.
4. **Open** a Pull Request (PR) explaining your changes.

Let's make HyperScanner even better together! 🌟

---

## 📄 License

This project is licensed under the **MIT License** — see the [LICENSE](LICENSE) file for details.

---

## 📧 Contact

For support, feedback, or feature requests:

- GitHub: [nxneeraj/hxscanner](https://github.com/nxneeraj/hxscanner)
- Email: neerajsahnx@gmail.com

---

## 🔮 Future Plans

- 🎯 Add terminal icons and detailed status code categories.
- 🖥️ Complete and release the GUI version.
- ⚡ Optimize performance for extremely large input lists.
- 🌍 Add proxy and multi-threaded support.

------

> Build faster. Test smarter. Hack ethically.  
> With 💥 from Team HyperGod-X 👾
<p align="center"><strong> Keep Moving Forward </strong></p>
