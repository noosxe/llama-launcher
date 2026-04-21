package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/spf13/viper"
)

type Model struct {
	Name            string   `mapstructure:"name"`
	ModelPath       string   `mapstructure:"model_path"`
	ModelFile       string   `mapstructure:"model_file"`
	ModelDir        string   `mapstructure:"model_dir"`
	ContainerName   string   `mapstructure:"container_name"`
	ContainerImage  string   `mapstructure:"container_image"`
	HostPort        int      `mapstructure:"host_port"`
	ContainerPort   int      `mapstructure:"container_port"`
	GPULayers       int      `mapstructure:"gpu_layers"`
	ContextSize     int      `mapstructure:"context_size"`
	Threads         int      `mapstructure:"threads"`
	BatchSize       int      `mapstructure:"batch_size"`
	NPredict        int      `mapstructure:"n_predict"`
	ChatTemplate    string   `mapstructure:"chat_template"`
	KVCacheQuantKey string   `mapstructure:"ctk"`
	KVCacheQuantVal string   `mapstructure:"ctv"`
	Annotations     []string `mapstructure:"annotations"`
}

type Config struct {
	ContainerImage  string  `mapstructure:"container_image"`
	Port            int     `mapstructure:"port"`
	ModelDir        string  `mapstructure:"model_dir"`
	NPredict        int     `mapstructure:"n_predict"`
	ChatTemplate    string  `mapstructure:"chat_template"`
	KVCacheQuantKey string  `mapstructure:"ctk"`
	KVCacheQuantVal string  `mapstructure:"ctv"`
	Models          []Model `mapstructure:"models"`
}

func Load(path string) (*Config, error) {
	v := viper.New()

	if path != "" {
		v.SetConfigFile(path)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("toml")

		v.AddConfigPath(filepath.Join(xdg.ConfigHome, "llama-launcher"))
		v.AddConfigPath(".")
	}

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func ConfigPath() string {
	configDir := xdg.ConfigHome
	return fmt.Sprintf("%s/llama-launcher/config.toml", configDir)
}

func EnsureConfig() error {
	configPath := ConfigPath()
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("config not found at %s", configPath)
	}
	return nil
}
