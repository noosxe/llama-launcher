# 🦙 Llama Launcher

**Llama Launcher** is a highly polished, interactive Terminal User Interface (TUI) built in Go to seamlessly orchestrate and manage `llama.cpp` Docker containers. 

Instead of juggling massive shell commands to spin up different local LLMs, Llama Launcher serves as a centralized dashboard. It allows you to rapidly deploy, supervise, monitor, and selectively tear down highly-configured local inference nodes!

## ✨ Features

- **Detached Orchestration**: Your inference instances are fully decoupled! You can launch a container inside the dashboard, press `Q` to completely exit the launcher, and your underlying `docker` workload will stay actively serving! Opening the dashboard again instantly rediscovers and re-attaches to the background logs.
- **Hardware Monitoring**: Check real-time Global Host Stats via the interactive footer, which utilizes native Linux tools (`top`, `free`, and `nvidia-smi`) to poll accurate CPU, RAM, GPU, and VRAM utilization continuously!
- **Catppuccin Themes**: Ships natively with full support for the gorgeous [Catppuccin](https://github.com/catppuccin/catppuccin) color standard. Customize your entire interface on the fly with a dedicated settings menu!
- **Dynamic Configuration**: Hot-swap configurations utilizing standard TOML tables. Expand standard pathing variables seamlessly into your volumes!

## ⚙️ Prerequisites

- **Go 1.22+**
- **Docker**
- **NVIDIA Container Toolkit** (for GPU acceleration/rendering)

## 🚀 Installation & Usage

1. **Clone the repo**
   ```bash
   git clone https://github.com/your-username/llama-launcher.git
   cd llama-launcher
   ```
2. **Build the dashboard**
   ```bash
   go build -o llama-launcher .
   ```
3. **Set up the config**
   Copy the example config and adjust the paths to point toward your native `.gguf` directories.
   ```bash
   cp config.example.toml config.toml
   ```
4. **Launch**
   ```bash
   ./llama-launcher tui
   ```

## 🎮 Interface Controls

- `Up` / `Down` or `J` / `K`: Navigate models in the sidebar.
- `Enter`: Boot or gracefully kill a model endpoint.
- `Escape`: Toggle the central Overlay Menu to switch internal Themes or confidently exit the application.
- `Q` or `Ctrl+C`: Gracefully sever the log listeners and exit the application instantly *(Note: Background processes stay completely unbothered!)*

## 📝 Configuration (`config.toml` structure)

Configurations are natively processed in TOML. The global variables allow for generic defaults, whilst individual model configurations explicitly override base defaults!

```toml
container_image = "local/llama.cpp:server-cuda"
host_port = 8080
container_port = 8080
model_dir = "~/models" # Universal Expansion fully supported!

[[models]]
name = "tiny-llama"
model_file = "tiny-llama-1b.q4_k_m.gguf"
container_name = "llama-tiny"
gpu_layers = 32
context_size = 2048

[[models]]
name = "mistral-7b"
model_file = "mistral-7b-instruct-v0.1.q4_k_m.gguf"
container_name = "llama-mistral"
host_port = 8081 # Override the default host-port!
gpu_layers = 32
context_size = 4096
```
