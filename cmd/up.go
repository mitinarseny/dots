// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Install dotfiles",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		data, err := ioutil.ReadFile(cfgFile)
		if err != nil {
			errLogger.Fatalf("An error occurred while openning file '%v': %v", cfgFile, err.Error())
		}
		dc.Source = cfgFile
		if err := yaml.Unmarshal(data, &dc); err != nil {
			errLogger.Fatalf("An error occurred while parsing '%v': %v", cfgFile, err.Error())
		}
		if err := dc.Revise(); err != nil {
			errLogger.Fatalln("An error occurred while revising ", err.Error())
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		up()
	},
}

func init() {
	rootCmd.AddCommand(upCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// upCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// upCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func center(s string, w int) string {
	return fmt.Sprintf("%[1]*s", -w, fmt.Sprintf("%[1]*s", (w+len(s))/2, s))
}

func left(s string, w int) string {
	return fmt.Sprintf("%[1]*s", -w, s)
}

func up() {
	fmt.Println("Creating symlinks...")

	//w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, ' ', 0)
	// find tab sizes
	var (
		maxTargetWidth int
		maxSourceWidth int
		maxStageWidth  = 11
	)
	for _, l := range dc.Links {
		if len(l.Target.Original) > maxTargetWidth {
			maxTargetWidth = len(l.Target.Original)
		}
		if len(l.Source.Original) > maxSourceWidth {
			maxSourceWidth = len(l.Source.Original)
		}
	}
	maxTargetWidth += 2
	maxSourceWidth += 2

	for _, l := range dc.Links {
		fmt.Print(fmt.Sprintf("%s\t<-\t %s\t|", left(l.Target.Original, maxTargetWidth), left(l.Source.Original, maxSourceWidth)))
		var targetBackup []byte
		if _, err := os.Lstat(l.Target.Absolute); err == nil {
			if !l.Force {
				fmt.Println(left("\t->\tomitted", maxStageWidth))
				continue
			} else {
				// backup
				//targetBackup, err = ioutil.ReadFile(l.Target.Absolute)
				//if err != nil {
				//	fmt.Printf(" -> failed to backup: %v\n", err.Error())
				//	continue
				//}
				if err := os.Remove(l.Target.Absolute); err != nil {
					fmt.Printf("\t->\tfailed to remove: %v\n", err.Error())
					continue
				}
				fmt.Print(left("\t->\tremoved", maxStageWidth))
			}
		}
		if err := os.Symlink(l.Source.Relative, l.Target.Absolute); err != nil {
			if targetBackup != nil {
				//os.NewFile() TODO: restore
			}
			fmt.Printf("\t->\terror: %v\n", err.Error())
			continue
		}
		fmt.Println("\t->\tcreated")
	}
	fmt.Println("Symlinks created!")
}
