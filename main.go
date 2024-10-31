package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func main() {
	if !isMultipassInstalled() {
		log.Fatal("Multipass is not installed. Please install Multipass and try again.")
	}

	config := map[string]string{
		"nodeName":       "ebpf",
		"sourcePath":     "/Users/naveen/go/src/github.com/naveensrinivasan",
		"sshKeyPath":     "/Users/naveen/.ssh/id_ed25519.pub",
		"remoteName":     "ebpf",
		"privateKeyPath": "/Users/naveen/.ssh/id_ed25519",
		"processor":      "arm64",
		"goVersion":      "1.23.2",
		"diskSize":       "50G",
		"cpuCount":       "2",
		"memory":         "4G",
		"ubuntuVersion":  "22.04",
		"gitEmail":       "172697+naveensrinivasan@users.noreply.github.com",
		"gitName":        "naveensrinivasan",
	}

	config["goDownload"] = fmt.Sprintf("https://go.dev/dl/go%s.linux-%s.tar.gz", config["goVersion"], config["processor"])

	packages := []string{
		"build-essential",
		"clang",
		"llvm",
		"libbpf-dev",
		"linux-headers-$(uname -r)",
		"linux-tools-common",
		"linux-tools-generic",
		"linux-tools-$(uname -r)",
		"zsh",
		"ca-certificates",
		"curl",
		"vim",
	}

	execCommand("multipass", "launch",
		config["ubuntuVersion"],
		"-n", config["nodeName"],
		"-m", config["memory"],
		"-c", config["cpuCount"],
		"-d", config["diskSize"],
		"--mount", fmt.Sprintf("%s:%s", config["sourcePath"], config["sourcePath"]))

	ipAddress := getIPAddress(config["nodeName"])

	execCommand("bash", "-c", fmt.Sprintf("cat %s | multipass exec %s -- bash -c 'cat >> ~/.ssh/authorized_keys'", config["sshKeyPath"], config["nodeName"]))
	sshConfig := fmt.Sprintf(`Host %s
    HostName %s
    User ubuntu
    IdentityFile %s
    ForwardAgent yes`, config["remoteName"], ipAddress, config["privateKeyPath"])

	execCommand("bash", "-c", fmt.Sprintf(`
		sed -i '/^Host %s$/,/^$/{/./d}' ~/.ssh/config 2>/dev/null || true
		echo '%s' >> ~/.ssh/config`,
		config["remoteName"], sshConfig))

	execCommand("multipass", "exec", config["nodeName"], "--", "bash", "-c", fmt.Sprintf("sudo apt-get update && sudo apt-get install -y %s", strings.Join(packages, " ")))

	// Download and install Go (with cleanup)
	execCommand("multipass", "exec", config["nodeName"], "--", "bash", "-c", fmt.Sprintf(`
		wget %s && \
		sudo tar -C /usr/local -xzf go%s.linux-%s.tar.gz && \
		rm go%s.linux-%s.tar.gz`,
		config["goDownload"],
		config["goVersion"],
		config["processor"],
		config["goVersion"],
		config["processor"]))

	// Update PATH for Go (remove separate source command)
	execCommand("multipass", "exec", config["nodeName"], "--", "bash", "-c",
		"echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc && . ~/.bashrc")

	// Install Oh My Zsh
	execCommand("multipass", "exec", config["nodeName"], "--", "bash", "-c", "sh -c \"$(wget https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh -O -)\" --unattended")

	// Add Go to PATH in .zshrc
	execCommand("multipass", "exec", config["nodeName"], "--", "bash", "-c",
		"echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.zshrc")

	execCommand("multipass", "exec", config["nodeName"], "--", "sudo", "chsh", "-s", "$(which zsh)", "ubuntu")

	execCommand("multipass", "exec", config["nodeName"], "--", "bash", "-c", `
		sudo install -m 0755 -d /etc/apt/keyrings && \
		sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc && \
		sudo chmod a+r /etc/apt/keyrings/docker.asc`)

	execCommand("multipass", "exec", config["nodeName"], "--", "bash", "-c", `
		echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
		$(. /etc/os-release && echo "$VERSION_CODENAME") stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null`)

	dockerPackages := []string{
		"docker-ce",
		"docker-ce-cli",
		"containerd.io",
		"docker-buildx-plugin",
		"docker-compose-plugin",
	}
	execCommand("multipass", "exec", config["nodeName"], "--", "bash", "-c",
		fmt.Sprintf("sudo apt-get update && sudo apt-get install -y %s", strings.Join(dockerPackages, " ")))

	execCommand("multipass", "exec", config["nodeName"], "--", "sudo", "docker", "run", "hello-world")

	// Add Git configuration after Docker setup
	execCommand("multipass", "exec", config["nodeName"], "--", "git", "config", "--global", "user.email", config["gitEmail"])
	execCommand("multipass", "exec", config["nodeName"], "--", "git", "config", "--global", "user.name", config["gitName"])

	// Set default editor to vim in shell
	execCommand("multipass", "exec", config["nodeName"], "--", "bash", "-c",
		"echo 'export EDITOR=vim' >> ~/.bashrc && . ~/.bashrc")
	execCommand("multipass", "exec", config["nodeName"], "--", "bash", "-c",
		"echo 'export EDITOR=vim' >> ~/.zshrc && . ~/.zshrc")

	execCommand("multipass", "exec", config["nodeName"], "--", "git", "config", "--global", "core.editor", "vim")
}

func execCommand(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()
	if err := cmd.Run(); err != nil {
		log.Fatalf("Command failed: %s\n", err)
	}
}

func getIPAddress(nodeName string) string {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("multipass info %s | grep IPv4 | awk '{print $2}'", nodeName))
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Failed to get IP address: %s\n", err)
	}
	return strings.TrimSpace(string(output))
}

func isMultipassInstalled() bool {
	cmd := exec.Command("multipass", "--version")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
