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
	"errors"
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

	log.Println("Setting variables...")
	for _, v := range c.Variables {
		if err := os.Setenv(v.Name, v.Value); err != nil {
			log.Fatalf("An error occured while setting '$%v=%v': %v", v.Name, v.Value, err.Error())
		}
	}
	log.Println("All variables are set...")

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

		if err := symlink(absSource, absTarget, true); err != nil {
			log.Fatalf("An error occured while making symlink '%v' -> '%v': %v", absSource, absTarget, err.Error())
		}

		log.Printf("Symlinked successfully: '%v' -> '%v'", absSource, absTarget)
	}
}

func symlink(source, target string, force bool) error {
	if _, err := os.Stat(source); os.IsNotExist(err) {
		return errors.New(fmt.Sprintf("Symlink source '%v' does not exist", source))
	}
	if _, err := os.Stat(target); err == nil {
		if !force {
			return errors.New(fmt.Sprintf("Symlink target '%v' already exists", target))
		}
		if err := os.Remove(target); err != nil {
			return err
		}
	}
	if err := os.Symlink(source, target); err != nil {
		return err
	}
	return nil
}
