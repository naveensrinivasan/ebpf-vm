<div align="center">

# 🚀 eBPF Dev Environment Launcher for macOS

*Your instant eBPF development environment, ready in minutes! 🎯*

</div>

---

🌟 Features

- 🔥 **Zero to eBPF** in one command
- 🐋 **Docker** pre-configured and ready to roll
- 🛠️ **Complete toolchain** with LLVM, Clang, and libbpf
- 🚄 **Lightning-fast** setup with Multipass
- 🔒 **Secure** SSH configuration out of the box
- 🎨 **Oh My Zsh** for that extra dev productivity

🏃‍♂️ Quick Start
```bash
git clone https://github.com/naveensrinivasan/ebpf-vm.git
cd ebpf-vm
go run main.go
```


🎮 What You Get

This tool automatically sets up a complete Ubuntu VM with:

- 🔧 Complete eBPF development toolchain
- 🐋 Docker CE with Compose and Buildx
- 🔨 Go development environment
- 🎨 Developer-friendly shell with Oh My Zsh
- 🔒 Pre-configured SSH access

🛠️ Tech Stack

| Category           | Components                            |
|--------------------|---------------------------------------|
| 🔧 **Build Tools** | `clang`, `llvm`, `build-essential`    |
| 🐧 **Linux Tools** | `libbpf-dev`, `linux-headers`, `linux-tools` |
| 🐋 **Containers**  | `docker-ce`, `buildx`, `compose`      |
| 🔨 **Dev Tools**   | `vim`, `curl`, `zsh`                   |

## ⚡ Connect to Your VM

Once the VM is up and running, you can connect to it using SSH:

```bash
ssh ebpf
```

## 🎛️ Configuration

Customize your environment in seconds:
```go
config := map[string]string{
"memory": "4G", // Boost it!
"cpuCount": "2", // More power!
"diskSize": "50G", // Store more!
}
```


## 🚀 System Requirements

- 🖥️ **Multipass** installed
- 🔑 SSH key pair (ed25519)

## 🤝 Contributing

Got ideas? We love them! Here's how you can contribute:

1. 🍴 Fork it
2. 🌿 Create your feature branch
3. 💾 Commit your changes
4. 🚀 Push to the branch
5. 🎉 Create a PR

## 📝 License

Apache 2.0

---

<div align="center">

Made with ❤️ for eBPF developers

*Happy Coding! 🎉*

</div>
