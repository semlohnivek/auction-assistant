// config/config.go
package config

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml"
)

type Config struct {
	Resize struct {
		MaxDimension int `toml:"max_dimension"`
	} `toml:"resize"`

	SFTP struct {
		Host       string `toml:"host"`
		Port       int    `toml:"port"`
		User       string `toml:"user"`
		Pass       string `toml:"pass"`
		RemotePath string `toml:"remote_path"`
	} `toml:"sftp"`

	OpenAI struct {
		Model        string  `toml:"model"`
		Key          string  `toml:"key"`
		MaxTokens    int     `toml:"max_tokens"`
		Temperature  float64 `toml:"temperature"`
		SystemPrompt string  `toml:"system_prompt"`
		ImageBaseUrl string  `toml:"image_base_url"`
	} `toml:"openai"`

	Flatten struct {
		IncludeBarcodeInUpload   bool   `toml:"include_barcode_in_upload"`
		IncludeBarcodeInAnalysis bool   `toml:"include_barcode_in_analysis"`
		BarcodeFirst             bool   `toml:"barcode_first"`
		LotPrefix                string `toml:"lot_prefix"`
		BarcodeSuffix            string `toml:"barcode_suffix"`
		ImagePrefix              string `toml:"image_prefix"`
		ImageExt                 string `toml:"image_ext"`
	} `toml:"flatten"`

	Defaults struct {
		AuctionSize int `toml:"auction_size"`
	} `toml:"defaults"`
}

var Current Config

func Load(path string) error {
	configFile, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("unable to read config file: %w", err)
	}
	if err := toml.Unmarshal(configFile, &Current); err != nil {
		return fmt.Errorf("unable to parse config file: %w", err)
	}
	return nil
}
