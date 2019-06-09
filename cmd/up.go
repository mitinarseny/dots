package cmd

import (
	"fmt"
	"github.com/mitinarseny/dots/config"
	"github.com/mitinarseny/dots/config/defaults"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"math"
	"os"
	"path"
	"path/filepath"
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

		if err := dc.Host.Revise(filepath.Dir(cfgFile)); err != nil {
			fmt.Println("An error occurred while revising ", err)
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := os.Chdir(path.Dir(cfgFile)); err != nil {
			fmt.Printf("An error occured while changing work directory: %s\n", err)
			os.Exit(1)
		}
		fmt.Println("Delivering common:")
		if err := upHost(&dc.Host); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if len(args) > 0 {
			hostName := args[0]
			h, exists := dc.Hosts[hostName]
			if !exists {
				fmt.Printf("there is no host '%s'", hostName)
				os.Exit(1)
			}
			if err := h.Revise(filepath.Dir(cfgFile)); err != nil {
				fmt.Println("An error occurred while revising ", err)
				os.Exit(1)
			}
			fmt.Printf("Delivering host '%s':\n", hostName)
			if err := upHost(dc.Hosts[hostName]); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
	//upCmd.Flags().StringVar(hostName, "host_name", "", "Host to use")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// upCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// upCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func upHost(host *config.Host) error {
	if err := setVariables(host.Variables); err != nil {
		return err
	}
	if err := createLinks(host.Links); err != nil {
		return err
	}
	if err := execCommands(host.Commands); err != nil {
		return err
	}
	if err := setDefaults(host.Defaults); err != nil {
		return err
	}
	return nil
}

func setVariables(vars config.Variables) error {
	if len(vars) == 0 {
		return nil
	}
	fmt.Printf("Variables (%d stages):\n", len(vars))
	for i, stage := range vars {
		vw := int(math.Log10(float64(len(stage))))
		for varName, variable := range stage {
			fmt.Printf("[%[1]*[2]d/%[1]*[3]d] %s=", vw, i+1, len(stage), varName)
			if variable.Command != nil {
				fmt.Printf("$(%s) -> ", variable.Command.String)
				if err := variable.FromCommand(); err != nil {
					return err
				}
			}
			if variable.Value != nil {
				fmt.Println(*variable.Value)
				if err := os.Setenv(varName, *variable.Value); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func createLinks(links config.Links) error {
	if len(links) == 0 {
		return nil
	}
	fmt.Printf("Links (%d):\n", len(links))
	lw := int(math.Log10(float64(len(links))))
	for i, l := range links {
		fmt.Printf("[%[1]*[2]d/%[1]*[3]d] %s <- %s: ", lw, i+1, len(links), l.Target.Original, l.Source.Original)
		st, err := l.Link()
		if err != nil {
			return err
		}
		switch st {
		case config.Created:
			fmt.Println("created")
		case config.Omitted:
			fmt.Printf("omitted (already exists, force: %t)\n", l.Force)
		case config.Replaced:
			fmt.Printf("replaced (already exists, force: %t)\n", l.Force)
		}
	}
	return nil
}

func execCommands(cmds config.Commands) error {
	if len(cmds) == 0 {
		return nil
	}
	nw := int(math.Log10(float64(len(cmds))))
	fmt.Printf("Commands (%d):\n", len(cmds))
	for i, c := range cmds {
		fmt.Printf("[%[1]*[2]d/%[1]*[3]d]", nw, i+1, len(cmds))
		if c.Description != nil {
			fmt.Print(" ", *c.Description)
		}
		fmt.Printf(":\n%s\n", c.String)
		c.Stdin = nil
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr

		if err := c.Run(); err != nil {
			return err
		}
	}
	return nil
}

func setDefaults(d defaults.Defaults) error {
	if len(d.Apps) == 0 && len(d.Domains) == 0 && len(d.Globals) == 0 {
		return nil
	}
	fmt.Println("Defaults:")
	for keyName, key := range d.Globals {
		if key.Description != nil {
			fmt.Printf("[%s] ", *key.Description)
		}
		cmdStr := fmt.Sprintf(
			"defaults write -globalDomain %s %s",
			keyName,
			key.Value.String())
		fmt.Println(cmdStr)
		cmd := new(config.Command).WithString(cmdStr)

		cmd.Stdin = nil
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return err
		}
	}

	for typ, domains := range map[string]defaults.Domains{
		" -app": d.Apps,
		"":      d.Domains,
	} {
		if len(domains) == 0 {
			continue
		}
		for domainName, domain := range domains {
			for keyName, key := range domain {
				if key.Description != nil {
					fmt.Printf("[%s] ", *key.Description)
				}
				cmdStr := fmt.Sprintf(
					"defaults write%s %s %s %s",
					typ,
					domainName,
					keyName,
					key.Value.String())
				fmt.Println(cmdStr)
				cmd := new(config.Command).WithString(cmdStr)

				cmd.Stdin = nil
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				if err := cmd.Run(); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
