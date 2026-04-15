package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"llama-launcher/config"
)

var runCmd = &cobra.Command{
	Use:   "run [model-name]",
	Short: "Run a selected model",
	Args:  cobra.ExactArgs(1),
	Run:   runModel,
}

func runModel(cmd *cobra.Command, args []string) {
	modelName := args[0]

	cfg, err := config.Load(cfgFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	var modelConfig *config.Model
	for i := range cfg.Models {
		if cfg.Models[i].Name == modelName {
			modelConfig = &cfg.Models[i]
			break
		}
	}

	if modelConfig == nil {
		fmt.Fprintf(os.Stderr, "Model '%s' not found\n", modelName)
		os.Exit(1)
	}

	containerImage := modelConfig.ContainerImage
	if containerImage == "" {
		containerImage = cfg.ContainerImage
	}
	if containerImage == "" {
		fmt.Fprintf(os.Stderr, "container_image not configured\n")
		os.Exit(1)
	}

	if err := startContainer(cfg, containerImage, modelConfig); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start container: %v\n", err)
		os.Exit(1)
	}
}

func startContainer(cfg *config.Config, image string, m *config.Model) error {
	modelDir := m.ModelDir
	if modelDir == "" {
		modelDir = cfg.ModelDir
	}
	if modelDir == "" && m.ModelPath != "" {
		modelDir = filepath.Dir(m.ModelPath)
	}
	if modelDir == "" {
		return fmt.Errorf("model directory not configured")
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
	containerPort := m.ContainerPort
	if containerPort == 0 {
		containerPort = 8080
	}

	dockerArgs := []string{
		"run",
		"--name", m.ContainerName,
		"--gpus", "all",
		"-p", fmt.Sprintf("%d:%d", hostPort, containerPort),
		"-v", fmt.Sprintf("%s:/models:ro", modelDir),
		"--rm",
		image,
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

	fmt.Printf("Starting container '%s' on port %d...\n", m.ContainerName, hostPort)
	fmt.Printf("Docker args: docker %v\n", dockerArgs)

	execCmd := exec.Command("docker", dockerArgs...)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Stdin = os.Stdin

	return execCmd.Run()
}

func init() {
	rootCmd.AddCommand(runCmd)
}
