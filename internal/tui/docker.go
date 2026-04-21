package tui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/noosxe/llama-launcher/internal/config"
)

func expandPath(pathStr string) string {
	if pathStr == "~" {
		if home, err := os.UserHomeDir(); err == nil {
			return home
		}
	} else if strings.HasPrefix(pathStr, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, pathStr[2:])
		}
	}
	return pathStr
}

func BuildDockerCmd(cfg *config.Config, m *config.Model) (*exec.Cmd, error) {
	containerImage := m.ContainerImage
	if containerImage == "" {
		containerImage = cfg.ContainerImage
	}
	if containerImage == "" {
		return nil, fmt.Errorf("container_image not configured")
	}

	modelDir := m.ModelDir
	if modelDir == "" {
		modelDir = cfg.ModelDir
	}
	if modelDir == "" && m.ModelPath != "" {
		modelDir = filepath.Dir(m.ModelPath)
	}
	if modelDir == "" {
		return nil, fmt.Errorf("model directory not configured")
	}
	
	modelDir = expandPath(modelDir)
	if absDir, err := filepath.Abs(modelDir); err == nil {
		modelDir = absDir
	}

	modelFile := m.ModelFile
	if modelFile == "" {
		modelFile = m.ModelPath
		if modelFile != "" {
			modelFile = filepath.Base(m.ModelPath)
		}
	}

	hostPort := m.HostPort
	if hostPort == 0 {
		hostPort = cfg.Port
	}
	if hostPort == 0 {
		hostPort = 8080 // default
	}

	containerPort := m.ContainerPort
	if containerPort == 0 {
		containerPort = 8080
	}

	dockerArgs := []string{
		"run", "-d",
		"--name", m.ContainerName,
		"--gpus", "all",
		"-p", fmt.Sprintf("%d:%d", hostPort, containerPort),
		"-v", fmt.Sprintf("%s:/models:ro", modelDir),
		"--rm",
		containerImage,
		"-m", fmt.Sprintf("/models/%s", modelFile),
		"--port", fmt.Sprintf("%d", containerPort),
		"--host", "0.0.0.0",
	}

	if m.GPULayers > 0 {
		dockerArgs = append(dockerArgs, "--n-gpu-layers", fmt.Sprintf("%d", m.GPULayers))
	}
	if m.ContextSize > 0 {
		dockerArgs = append(dockerArgs, "-c", fmt.Sprintf("%d", m.ContextSize))
	}
	if m.Threads > 0 {
		dockerArgs = append(dockerArgs, "-t", fmt.Sprintf("%d", m.Threads))
	}
	if m.BatchSize > 0 {
		dockerArgs = append(dockerArgs, "-b", fmt.Sprintf("%d", m.BatchSize))
	}

	nPredict := m.NPredict
	if nPredict == 0 {
		nPredict = cfg.NPredict
	}
	if nPredict != 0 {
		dockerArgs = append(dockerArgs, "-n", fmt.Sprintf("%d", nPredict))
	}

	chatTemplate := m.ChatTemplate
	if chatTemplate == "" {
		chatTemplate = cfg.ChatTemplate
	}
	if chatTemplate != "" {
		dockerArgs = append(dockerArgs, "--chat-template", chatTemplate)
	}

	kvCacheQuantKey := m.KVCacheQuantKey
	if kvCacheQuantKey == "" {
		kvCacheQuantKey = cfg.KVCacheQuantKey
	}
	if kvCacheQuantKey != "" {
		dockerArgs = append(dockerArgs, "-ctk", kvCacheQuantKey)
	}

	kvCacheQuantVal := m.KVCacheQuantVal
	if kvCacheQuantVal == "" {
		kvCacheQuantVal = cfg.KVCacheQuantVal
	}
	if kvCacheQuantVal != "" {
		dockerArgs = append(dockerArgs, "-ctv", kvCacheQuantVal)
	}

	return exec.Command("docker", dockerArgs...), nil
}

func BuildDockerLogsCmd(containerName string) *exec.Cmd {
	return exec.Command("docker", "logs", "-f", containerName)
}
