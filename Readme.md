# 🚀 Skipper

> A powerful tunnel and reverse proxy that allows you to expose your localhost projects to the internet through a subdomain. Simple, fast, and secure - your local development environment, accessible from anywhere.

<div align="center">

<p align="center">
  <img src="https://pub-f0e906fcf8d54ea98e4cbbd15a55e147.r2.dev/skipperlogo.png" alt="Skipper Logo" width="250"/>
</p>


**Made with ❤️ by Juan Diego Diaz**

[![GitHub stars](https://img.shields.io/github/stars/juandi247/skipper?style=social)](https://github.com/juandi247/skipper)
[![GitHub forks](https://img.shields.io/github/forks/juandi247/skipper?style=social)](https://github.com/juandi247/skipper)
[![GitHub issues](https://img.shields.io/github/issues/juandi247/skipper)](https://github.com/juandi247/skipper/issues)
[![GitHub license](https://img.shields.io/github/license/juandi247/skipper)](https://github.com/juandi247/skipper/blob/main/LICENSE)

</div>

---

## ✨ Features

- 🚀 **Simple Setup** - One command to expose your localhost
- ⚡ **Fast** - Built with Go for high performance
- 🌐 **Global Access** - Access your projects from anywhere
- 🎯 **Custom Subdomains** - Choose your own subdomain

## 🛠️ Installation

### Package Managers

#### Homebrew (macOS)
```bash
brew tap juandi247/skipper
brew install skipper
```

#### Chocolatey (Windows)
```bash
choco install skipper
```
> ⚠️ **Pending for approval**

### Direct Downloads

Download the latest version (v0.0.7) directly for your platform:

| Platform | Download |
|----------|----------|
| Linux AMD64 | [skipper_0.0.7_linux_amd64.tar.gz](https://github.com/juandi247/Skipper/releases/download/v0.0.7/skipper_0.0.7_linux_amd64.tar.gz) |
| Linux ARM64 | [skipper_0.0.7_linux_arm64.tar.gz](https://github.com/juandi247/Skipper/releases/download/v0.0.7/skipper_0.0.7_linux_arm64.tar.gz) |
| Windows AMD64 | [skipper_0.0.7_windows_amd64.tar.gz](https://github.com/juandi247/Skipper/releases/download/v0.0.7/skipper_0.0.7_windows_amd64.tar.gz) |
| Windows ARM64 | [skipper_0.0.7_windows_arm64.tar.gz](https://github.com/juandi247/Skipper/releases/download/v0.0.7/skipper_0.0.7_windows_arm64.tar.gz) |
| macOS AMD64 | [skipper_0.0.7_darwin_amd64.tar.gz](https://github.com/juandi247/Skipper/releases/download/v0.0.7/skipper_0.0.7_darwin_amd64.tar.gz) |
| macOS ARM64 (M1) | [skipper_0.0.7_darwin_arm64.tar.gz](https://github.com/juandi247/Skipper/releases/download/v0.0.7/skipper_0.0.7_darwin_arm64.tar.gz) |

📦 **[View all releases on GitHub](https://github.com/juandi247/Skipper/releases/latest)**

## 🚀 Quick Start

Getting started with Skipper is as simple as running a single command:

```bash
skipper -port 3000 -subdomain myapp
```

This will expose your local application running on port 3000 at:
```
https://myapp.skipper.lat
```

### Parameters
- `-port`: The port where your local application is running
- `-subdomain`: The subdomain you want to use (will be available at subdomain.skipper.lat)

## 🏗️ Architecture

Skipper operates through a sophisticated architecture with advanced Go patterns:

### 🚀 Core Architecture
- **TCP Persistent Connection**: Bidirectional communication between proxy and tunnel
- **Custom Framing Protocol**: 20-byte binary header for efficient message routing
- **Stream Multiplexing**: Multiple flows over single TCP connection

### 🧠 Concurrency Patterns
- **Worker Pools**: Efficient concurrent task management
- **Fan-In/Fan-Out**: Dynamic message distribution patterns
- **Reactor Pattern**: Message dispatching by type and StreamID

### 🎛️ State Machine Design
- **Functional FSM**: Tunnel as concurrent state machine
- **Graceful Shutdown**: Controlled resource cleanup with context
- **Error Handling**: Robust state transitions and recovery

### 📦 Protocol & Serialization
- **Protocol Buffers**: Compact, fast serialization
- **Go Idioms**: Context, select, interfaces
- **Extensible Design**: Decoupled handlers for testing

## 🐕 Meet Skipper

<div align="center">

<p align="center">
  <img src="https://pub-f0e906fcf8d54ea98e4cbbd15a55e147.r2.dev/skipper2.jpeg" alt="Skipper Photos" width="600"/>
</p>

**Skipper**: looks cute, contributes nothing to the repo
</div>

## 🤝 Contributing

We welcome contributions from the community! Here's how you can help:

- 🐛 **Report bugs** on our [GitHub Issues](https://github.com/juandi247/skipper/issues) page
- 💡 **Suggest new features** or improvements
- 📝 **Help improve our documentation**
- 🔧 **Submit a Pull Request** with your improvements

## 📞 Contact

- **GitHub**: [@juandi247](https://github.com/juandi247)
- **Email**: juand.diaza@gmail.com
- **LinkedIn**: [Juan Diego Diaz](https://www.linkedin.com/in/juan-diego-diaz-92a857238/)

---

<div align="center">

**⭐ Star this repository if you found it helpful!**

[![GitHub stars](https://img.shields.io/github/stars/juandi247/skipper?style=for-the-badge&logo=github)](https://github.com/juandi247/skipper)

</div>
