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
	. "github.com/mitinarseny/dots/config"
	"github.com/spf13/cobra"
	"gopkg.in/mattes/go-expand-tilde.v1"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Install dotfiles",
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

func up() {
	data, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		log.Fatalf("An error occured while opening file '%v': %v", cfgFile, err.Error())
	}
	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		log.Fatalf("An error occured while parsing '%v': %v", cfgFile, err.Error())
	}

	fmt.Println("Setting variables...")
	for _, v := range c.Variables {
		if err := os.Setenv(v.Name, v.Value); err != nil {
			log.Fatalf("An error occured while setting '$%v=%v': %v", v.Name, v.Value, err.Error())
		}
		fmt.Printf("Variable was successfully set: %v=%v", v.Name, v.Value)
	}
	fmt.Println("Creating symlinks...")
	for _, l := range c.Links {
		s, err := tilde.Expand(l.Source)
		if err != nil {
			log.Fatalf("An error occured while resolving '%v': %v", l.Source, err.Error())
		}
		s = os.ExpandEnv(s)
		absSource, err := filepath.Abs(s)
		if err != nil {
			log.Fatalf("An error occured while trying to get absolete path of '%v': %v", s, err.Error())
		}

		t, err := tilde.Expand(l.Target)
		if err != nil {
			log.Fatalf("An error occured while resolving '%v': %v", l.Target, err.Error())
		}
		t = os.ExpandEnv(t)
		absTarget, err := filepath.Abs(t)
		if err != nil {
			log.Fatalf("An error occured while trying to get absolete path of '%v': %v", t, err.Error())
		}

		if _, err := os.Stat(absSource); os.IsNotExist(err) {
			log.Fatalf("Symlink source '%v' does not exist", absSource)
		}
		if _, err := os.Stat(absTarget); err == nil {
			fmt.Printf("Symlink source already exists: '%v' -> ", absTarget)
			if !l.Force {
				fmt.Println("omitted")
			} else {
				if err := os.Remove(absTarget); err != nil {
					log.Fatalf("An error occurred while removing existing target '%v': %v", absSource, err.Error())
				} else {
					fmt.Println("removed")
				}

				if err := os.Symlink(absSource, absTarget); err != nil {
					log.Fatalf("An error occured while creating symlink '%v' -> '%v': %v", absSource, absTarget, err.Error())
				}
				fmt.Printf("Symlink was successfully created: '%v' -> '%v'\n", absSource, absTarget)
			}
		} else {
			if err := os.Symlink(absSource, absTarget); err != nil {
				log.Fatalf("An error occured while creating symlink '%v' -> '%v': %v", absSource, absTarget, err.Error())
			}
			fmt.Printf("Symlink was successfully created: '%v' -> '%v'\n", absSource, absTarget)
		}


	}
}
