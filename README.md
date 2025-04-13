# ⚡ HyperScanner (hxscanner) v1.0

**HyperScanner** is a powerful and futuristic HTTP status code scanner that takes a list of IPs and classifies their responses into structured folders and files. Designed with speed, clarity, and beauty in mind — it’s perfect for hackers, sysadmins, and bug bounty hunters ⚔️

---

## ✨ Features

- 🔍 Scans any IP list and fetches HTTP status codes
- 🗂️ Auto-creates folder structure by status code (1xx → 5xx)
- 📁 Saves every response IP into its status code file (e.g., `200.txt`)
- ✅ Clean CLI with logs: `ip_exist.txt`, `ip_invalid.txt`, `log.txt`
- 🎨 Terminal color-coded output (coming live with icons & categories)
- 🖥️ GUI version in progress (Wails frontend)
- 💻 Cross-platform (Windows/Linux/Mac)

---

## 🚀 Installation

Make sure you have **Go 1.19+** installed. Then run:

```bash
go install github.com/nxneeraj/hxscanner@latest
```
