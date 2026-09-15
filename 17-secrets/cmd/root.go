package cmd

import (
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "secret",
	Short: "Secret is an API key and secrets manager",
}

var encodingKey string

func init() {
	RootCmd.PersistentFlags().StringVarP(&encodingKey, "key", "k", "", "the key for encoding and decoding strings")
}

func secretsPath() string {
	home, _ := homedir.Dir()
	return filepath.Join(home, ".secrets")
}
