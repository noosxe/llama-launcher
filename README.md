# 🦙 Llama Launcher

[![Release](https://github.com/noosxe/llama-launcher/actions/workflows/release.yml/badge.svg)](https://github.com/noosxe/llama-launcher/actions/workflows/release.yml) [![Go Reference](https://pkg.go.dev/badge/github.com/noosxe/llama-launcher/cmd/llama-launcher.svg)](https://pkg.go.dev/github.com/noosxe/llama-launcher/cmd/llama-launcher) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> [!CAUTION]
> **Warning**: This project is fully **vibe-coded**. Proceed with appropriate levels of curiosity and caution. The auther (me) will be adding various comments throughout this readme and maybe other documents or code manually. You will probably notice them immediately because they will not align with general "vibe" of llms. 

**Llama Launcher** is a highly polished (not really!), interactive Terminal User Interface (TUI) built in Go to seamlessly orchestrate and manage `llama.cpp` Docker containers. 

Instead of juggling massive shell commands to spin up different local LLMs, Llama Launcher serves as a centralized dashboard. It allows you to rapidly deploy, supervise, monitor, and selectively tear down highly-configured local inference nodes!

## ✨ Features

- **Detached Orchestration**: Your inference instances are fully decoupled! You can launch a container inside the dashboard, press `Q` to completely exit the launcher, and your underlying `docker` workload will stay actively serving! Opening the dashboard again instantly rediscovers and re-attaches to the background logs.
- **Hardware Monitoring**: Check real-time Global Host Stats via the interactive footer, which utilizes native Linux tools (`top`, `free`, and `nvidia-smi`) to poll accurate CPU, RAM, GPU, and VRAM utilization continuously!
- **Catppuccin Themes**: Ships natively with full support for the gorgeous [Catppuccin](https://github.com/catppuccin/catppuccin) color standard. Customize your entire interface on the fly with a dedicated settings menu! (Disclaimer: I've tested the "mocha" variant only, as it is the color scheme I use in my terminal. Other variants of catppuccin may not work as intended. I'm an idiot, not a masochist. (ok, I've left the idiot not a masochist part here because gemini suggested that when I was typing. :shrug:))
- **Dynamic Configuration**: Hot-swap configurations utilizing standard TOML tables. Expand standard pathing variables seamlessly into your volumes! (no idea what it meant by hot-swap here, not tested at all!)

## ⚙️ Prerequisites

- **Go 1.22+**
- **Docker**
- **NVIDIA Container Toolkit** (for GPU acceleration/rendering)

## 🚀 Installation & Usage

1. **Install directly**
   ```bash
   go install github.com/noosxe/llama-launcher/cmd/llama-launcher@latest
   ```
   *Alternatively, clone and build:*
   ```bash
   git clone https://github.com/noosxe/llama-launcher.git
   cd llama-launcher
   go build -o llama-launcher ./cmd/llama-launcher
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

## 📝 Configuration (`config.toml`)

Llama Launcher utilizes a global/local cascading configuration. Values defined at the top level act as defaults and can be overridden by individual model definitions.

### Global Options

| Option | Description |
|--------|-------------|
| `container_image` | Default Docker image to use (e.g., `ghcr.io/ggerganov/llama.cpp:server`) |
| `model_dir` | The base directory on your host where models are stored. Supports `~` expansion. |
| `port` | Default host port to map if not specified in the model. |
| `n_predict` | Default number of tokens to predict (`-n`). |
| `chat_template` | Default chat template name or content (`--chat-template`). |
| `ctk` | Default KV cache key quantization (`-ctk`). |
| `ctv` | Default KV cache value quantization (`-ctv`). |

### Model Options (`[[models]]`)

Each entry in the `[[models]]` list defines a specific model configuration.

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `name` | string | **Yes** | Display name in the TUI sidebar. |
| `container_name`| string | **Yes** | Unique name for the Docker container (`--name`). |
| `model_file` | string | No* | The GGUF filename inside `model_dir`. |
| `model_path` | string | No* | Full path to the model file. (Alternative to `model_file`). |
| `container_image`| string | No | Override for the global `container_image`. |
| `model_dir` | string | No | Override for the global `model_dir`. |
| `host_port` | int | No | The port on the host machine to bind to. |
| `container_port`| int | No | The port inside the container (default: 8080). |
| `gpu_layers` | int | No | Number of layers to offload to GPU (`--n-gpu-layers`). |
| `context_size`| int | No | Context window size (`-c`). |
| `threads` | int | No | Number of threads to use (`-t`). |
| `batch_size` | int | No | Physical batch size (`-b`). |
| `n_predict` | int | No | Number of tokens to predict (`-n`). |
| `chat_template`| string | No | Chat template override (`--chat-template`). |
| `ctk` | string | No | KV cache key quantization override (`-ctk`). |
| `ctv` | string | No | KV cache value quantization override (`-ctv`). |

*\* Note: You must provide either `model_file` (if `model_dir` is set) or `model_path` so the launcher can find the model.*

### Full Example

```toml
container_image = "local/llama.cpp:server-cuda"
model_dir = "~/models"
port = 8080

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
host_port = 8081
gpu_layers = 33
ctk = "q8_0"
```
