package cmd

import (
	"fmt"

	"github.com/AbePlays/gophercises/17-secrets/secret"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Sets a secret in your secrets file",
	Run: func(cmd *cobra.Command, args []string) {
		v := secret.File(encodingKey, secretsPath())
		key, value := args[0], args[1]
		err := v.Set(key, value)
		if err != nil {
			panic(err)
		}
		fmt.Println("Secret set successfully")
	},
}

func init() {
	RootCmd.AddCommand(setCmd)
}
