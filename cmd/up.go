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
	"github.com/mitinarseny/dots/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
)

var (
	hostName string
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Install dotfiles",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		data, err := ioutil.ReadFile(cfgFile)
		if err != nil {
			errLogger.Fatalf("An error occurred while openning file '%v': %v", cfgFile, err)
		}
		if err := yaml.Unmarshal(data, &dc); err != nil {
			errLogger.Fatalf("An error occurred while parsing '%v': %v", cfgFile, err)
		}
		if err := dc.Revise(filepath.Dir(cfgFile)); err != nil {
			errLogger.Fatalln("An error occurred while revising ", err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		up()
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
	upCmd.Flags().StringVarP(&hostName, "host_name", "H", "", "Host to use")

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
	link()
	// commands

	//fmt.Printf("Executing commands (%d):\n", len(dc.Commands))
	//
	//if err := os.Chdir(path.Dir(dc.Source)); err != nil {
	//	fmt.Printf("An error occured while changing work directory: %s\n", err)
	//	os.Exit(1)
	//}
	//nw := int(math.Log10(float64(len(dc.Commands))))
	//for i, c := range dc.Commands {
	//	fmt.Printf("%[1]*[2]d/%[1]*[3]d: %[4]s\n", nw, i+1, len(dc.Commands), *c)
	//
	//	cmd := exec.Command("sh", "-c", string(*c))
	//	cmdReader, err := cmd.StdoutPipe()
	//	if err != nil {
	//		_, _ = fmt.Fprintf(os.Stderr, "An error occured while acquiring pipe: %s\n", err)
	//		continue
	//	}
	//
	//	scanner := bufio.NewScanner(cmdReader)
	//	go func() {
	//		for scanner.Scan() {
	//			fmt.Println(scanner.Text())
	//		}
	//	}()
	//
	//	err = cmd.Start()
	//	if err != nil {
	//		_, _ = fmt.Fprintf(os.Stderr, "An error occurred while starting command execution: %s\n", err)
	//		continue
	//	}
	//
	//	err = cmd.Wait()
	//	if err != nil {
	//		_, _ = fmt.Fprintf(os.Stderr, "An error occurred while waitng: %s\n", err)
	//		continue
	//	}
	//}
}

func link() {
	ll := dc.Links
	if hostName != "" {
		h, exists := dc.Hosts[hostName]
		if !exists {
			fmt.Printf("There is no host '%s' in %s\n", hostName, cfgFile)
			return
		}
		ll = append(ll, h.Links...)
	}
	fmt.Printf("Creating symlinks (%d):\n", len(ll))
	for _, l := range ll {
		fmt.Printf("%s <- %s: ", l.Target.Original, l.Source.Original)
		st, err := l.Link()
		if err != nil {
			fmt.Println(err)
			continue
		}
		switch st {
		case config.Omitted:
			fmt.Println("omitted (already exists)")
		case config.Created:
			fmt.Println("created")
		case config.Replaced:
			fmt.Println("replaced")
		}
	}
}
