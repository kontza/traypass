/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"embed"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var assets embed.FS

type AppConfig struct {
	ScanDirectory string `mapstructure:"scan-directory"`
	cfgFile       string
	verbose       bool
}

var appConfig AppConfig

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "waitwhat",
	Short: "Trials on Wails + HTMX",
	Run:   rootRunner,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(appAssets embed.FS) {
	assets = appAssets
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&appConfig.cfgFile, "config", "", fmt.Sprintf("config file (default is $HOME/.%s.yaml)", rootCmd.Use))
	rootCmd.Flags().StringVarP(&appConfig.ScanDirectory, "scan-directory", "d", "$HOME/.local/share/gopass/stores", "The directory to scan for GPG files")
	rootCmd.Flags().BoolVarP(&appConfig.verbose, "verbose", "v", false, "Show verbose logging information.")

	viper.BindPFlags(rootCmd.Flags())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if appConfig.cfgFile != "" {
		if configPath, err := filepath.Abs(appConfig.cfgFile); err != nil {
			slog.Error("Failed to get absolute path for '%s': %v", appConfig.cfgFile, err)
			os.Exit(1)
		} else {
			viper.SetConfigFile(configPath)
		}
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".changeme" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(fmt.Sprintf(".%s", rootCmd.Use))
	}

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
		if err := viper.Unmarshal(&appConfig); err != nil {
			slog.Error("Failed to unmarshal config", "error", err)
			os.Exit(1)
		}
	}
}
