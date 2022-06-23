/*
Copyright © 2022 Michael Reichenbach

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/spf13/pflag"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const defaultTimeout = 5 * time.Second
const envPrefix = "IMAP_CLI"

var (
	cfgFile       string
	server        string
	username      string
	password      string
	mailbox       string
	useTls        bool
	skipVerifyTls bool
	timeout       time.Duration
)

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Version: "1.0.0",
		Use:     "go-imap-client",
		Short:   "A CLI to connect to an IMAP mailbox and search it.",
		Long: `The go-imap-client is a CLI that enables quick searching of an IMAP mailbox.
This can be useful in automated environments, like CI/CD pipelines, to check if a mail arrived in the given inbox.

Usage Example:
imap --server "my-server:993" --username "username" --password "password" --tls search my mail`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			initConfig(cmd)
		},
	}

	configureFlags(rootCmd)
	bindFlagsToConfig(rootCmd)

	addChildCommands(rootCmd)

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := newRootCmd().Execute()
	if err != nil {
		os.Exit(1)
	}
}

func configureFlags(rootCmd *cobra.Command) {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.imap-cli.yaml)")

	rootCmd.PersistentFlags().StringVarP(&server, "server", "s", "", "mail server including port, e.g. --server=imap.my-server.com:143")
	rootCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "username to use for the connection, e.g. --username=admin")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "password of the username, e.g. --password=my-password")

	rootCmd.MarkFlagRequired("server")
	rootCmd.MarkFlagRequired("username")
	rootCmd.MarkFlagRequired("password")

	rootCmd.PersistentFlags().StringVarP(&mailbox, "mailbox", "m", imap.InboxName, "name of the mailbox")
	rootCmd.PersistentFlags().DurationVar(&timeout, "timeout", defaultTimeout, "timeout for the connection to the mail server")

	rootCmd.PersistentFlags().BoolVar(&useTls, "tls", false, "set to connect using tls (default is false)")
	rootCmd.PersistentFlags().BoolVar(&skipVerifyTls, "skip-verify", false, "set to skip the verification of the server certificate (default is false)")
}

func bindFlagsToConfig(cmd *cobra.Command) {
	viper.BindPFlags(cmd.PersistentFlags())
	viper.BindPFlags(cmd.Flags())
}

func addChildCommands(rootCmd *cobra.Command) {
	rootCmd.AddCommand(newSearchCommand())
}

// initConfig reads in config file and ENV variables if set.
func initConfig(cmd *cobra.Command) {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".go-imap-client" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".imap-cli")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		cmd.Println(fmt.Sprintf("Using config file: %q", viper.ConfigFileUsed()))
	}

	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv() // read in environment variables that match

	bindFlags(cmd)
}

// Bind each cobra flag to its associated viper configuration (config file and environment variable)
func bindFlags(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Environment variables can't have dashes in them, so bind them to their equivalent
		// keys with underscores, e.g. --favorite-color to STING_FAVORITE_COLOR
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			viper.BindEnv(f.Name, fmt.Sprintf("%s_%s", envPrefix, envVarSuffix))
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && viper.IsSet(f.Name) {
			val := viper.Get(f.Name)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}
