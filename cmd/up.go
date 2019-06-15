package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path"
)

const (
	DefaultHostName = "default"
)

// upCmd represents the upHost command
var upCmd = &cobra.Command{
	Use:   "up [host]",
	Short: "Install dotfiles",
	Args:  cobra.RangeArgs(0, 1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		data, err := ioutil.ReadFile(cfgFile)
		if err != nil {
			fmt.Printf("An error occurred while openning file '%v': %v\n", cfgFile, err)
			os.Exit(1)
		}
		if err := yaml.Unmarshal(data, &dc); err != nil {
			fmt.Printf("An error occurred while parsing '%v': %v\n", cfgFile, err)
			os.Exit(1)
		}
		if err := dc.Inspect(); err != nil {
			fmt.Println("An error occurred while inspecting config ", err)
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := os.Chdir(path.Dir(cfgFile)); err != nil {
			fmt.Printf("An error occured while changing work directory: %s\n", err)
			os.Exit(1)
		}
		switch {
		case len(args) == 0:
			up(DefaultHostName)
		case len(args) == 1:
			hostName := args[0]
			up(hostName)
		}
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}

func up(hostName string) {
	if err := dc.Up(hostName); err != nil {
		fmt.Println(err)
	}
}
