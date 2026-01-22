package cmd

import (
	"fmt"
	"gen-cert/internal"
	"os"

	"github.com/spf13/cobra"
)

var (
	// some global variables
	certName string
	keyName  string

	InitError error

	rootCmd = &cobra.Command{
		SilenceUsage: true,
		Use:          "gen-cert",
		Short:        "generate certificate authority for https",
		Long:         fmt.Sprintf(``),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if InitError != nil {
				os.Exit(1)
			}
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// configuration variables
	rootCmd.PersistentFlags().StringVarP(&certName, "cert", "c", "dns-api.crt",
		"generate certificate")
	rootCmd.PersistentFlags().StringVarP(&keyName, "key", "k", "dns-api.key",
		"generate certificate key")

	// Generating certificate and a key, we use self-signed
	// certificates and custom header http authorization with
	// shared secret
	InitError = internal.GenCertKeyPair(certName, keyName)
	if InitError != nil {
		fmt.Println(InitError)
	}

	return
}
