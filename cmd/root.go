package cmd

import (
	"fmt"
	"github.com/mitinarseny/dots/core"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	dc      core.Config
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dots",
	Short: "Delivery tool for dotfiles",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
		os.Exit(1)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "core-file", "c", ".dots.yaml", "core file")
}

// initConfig reads in core file and ENV variables if set.
func initConfig() {
	//if cfgFile != "" {
	//	// Use core file from the flag.
	//	viper.SetConfigFile(cfgFile)
	//} else {
	//	// Find home directory.
	//	home, err := homedir.Dir()
	//	if err != nil {
	//		fmt.Println(err)
	//		os.Exit(1)
	//	}
	//
	//	// Search core in home directory with name ".dots" (without extension).
	//	viper.AddConfigPath(home)
	//	viper.SetConfigName(".dots")
	//}

	viper.AutomaticEnv() // read in environment variables that match

	// If a core file is found, read it in.
	//if err := viper.ReadInConfig(); err == nil {
	//	fmt.Println("Using core file:", viper.ConfigFileUsed())
	//}
}
