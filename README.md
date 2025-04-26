# ⚡ HyperScanner + CORS (hxscanner) v1.4+

**HyperScanner** is a powerful and futuristic HTTP status code scanner that takes a list of IPs and classifies their responses into structured folders and files. Designed with speed, clarity, and beauty, HyperScanner simplifies the process of analyzing HTTP responses.

---

## ✨ Features

- 🔍 **Scan IP & URLs Lists**: Quickly scans any list of IPs and fetches their HTTP status codes.
- 🗂️ **Organized Output**: Automatically creates folder structures based on HTTP status codes (e.g., `1xx`, `2xx`, `3xx`, `4xx`, `5xx`).
- 📁 **File-Based Response Storage**: Saves each IP's response into corresponding status code files, such as `200.txt` for HTTP 200 responses.
- ✅ **Detailed Logs**: Tracks results in clean and structured log files:
  - `ip_exist.txt`: List of valid IPs.
  - `ip_invalid.txt`: List of invalid IPs.
  - `log.txt`: Comprehensive log of the scan.
- 🎨 **Enhanced CLI**: Provides color-coded terminal outputs for improved readability (coming soon: icons and categories).
- 💻 **Cross-Platform**: Compatible with Windows, Linux, and Mac.

---

## 🚀 Installation

Ensure you have **Go 1.19+** installed on your system. Then, install HyperScanner using:

```bash
go install github.com/nxneeraj/hxscanner@latest
```

---

## 🛠️ Usage

1. Prepare a text file containing a list of IPs, one IP per line (e.g., `ips.txt`).
2. Run the scanner with the following command:

```bash
hxscanner -i ips.txt
```

3. Results will be saved into structured folders and files based on HTTP status codes.

### CLI Options

- `-i <file>`: Specify the input file containing the list of IPs.
- `-o <directory>`: (Optional) Specify the output directory for results. Default is the current directory.
- `--help`: Displays help information.

---

## 📂 Output Structure

After running HyperScanner, the output will be structured as follows:

```
output/
├── 1xx/
├── 2xx/
│   ├── 200.txt
│   ├── ...
├── 3xx/
├── 4xx/
├── 5xx/
├── ip_exist.txt
├── ip_invalid.txt
└── log.txt
```

- **`<status_code>.txt`**: Contains the list of IPs that returned the corresponding HTTP status code.
- **`ip_exist.txt`**: List of successfully scanned IPs.
- **`ip_invalid.txt`**: List of invalid or unreachable IPs.
- **`log.txt`**: Detailed log of the scanning process.

---

## 🌐 Cross-Platform Compatibility

HyperScanner is built in **Go**, making it compatible with the following platforms:

- **Windows**
- **Linux**
- **MacOS**

---

## 🖥️ GUI Version (In Progress)

The GUI version of HyperScanner is underway, built using **Wails** for a modern desktop experience. Stay tuned for updates!

---

## 🏗️ Contributing

Contributions are welcome! If you'd like to improve HyperScanner or fix any issues, please follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Submit a pull request with a detailed description of your changes.

---

## 📄 License

HyperScanner is open-source software licensed under the [MIT License](LICENSE).
---

## 📧 Contact

For support, feature requests, or any queries, feel free to reach out:

- GitHub: [nxneeraj/hxscanner](https://github.com/nxneeraj/hxscanner)
- Email: neerajsahnx@gmail.com 

---

## 🔮 Future Plans

- Add terminal icons and detailed categories for better CLI output.
- Complete and release the GUI version.
- Enhance performance and scalability for large IP lists.

> Build faster. Test smarter. Hack ethically.  
> With 💥 from Team HyperGod-X 👾
<p align="center"><strong> Keep Moving Forward </strong></p>
