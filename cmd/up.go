package cmd

import (
	"fmt"
	"github.com/mitinarseny/dots/config"
	"github.com/mitinarseny/dots/config/defaults"
	"github.com/spf13/cobra"
	"go.uber.org/zap/buffer"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"math"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// upCmd represents the upHost command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Install dotfiles",
	Args:  cobra.RangeArgs(0, 1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		data, err := ioutil.ReadFile(cfgFile)
		if err != nil {
			errLogger.Fatalf("An error occurred while openning file '%v': %v", cfgFile, err)
		}
		if err := yaml.Unmarshal(data, &dc); err != nil {
			errLogger.Fatalf("An error occurred while parsing '%v': %v", cfgFile, err)
		}

		if err := dc.Host.Revise(filepath.Dir(cfgFile)); err != nil {
			errLogger.Fatalln("An error occurred while revising ", err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := os.Chdir(path.Dir(cfgFile)); err != nil {
			fmt.Printf("An error occured while changing work directory: %s\n", err)
			os.Exit(1)
		}

		upHost(&dc.Host)

		if len(args) > 0 {
			hostName := args[0]
			h, exists := dc.Hosts[hostName]
			if !exists {
				fmt.Printf("there is no host '%s'", hostName)
				os.Exit(1)
			}
			if err := h.Revise(filepath.Dir(cfgFile)); err != nil {
				errLogger.Fatalln("An error occurred while revising ", err)
			}
			upHost(dc.Hosts[hostName])
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

func center(s string, w int) string {
	return fmt.Sprintf("%[1]*s", -w, fmt.Sprintf("%[1]*s", (w+len(s))/2, s))
}

func left(s string, w int) string {
	return fmt.Sprintf("%[1]*s", -w, s)
}

func upHost(host *config.Host) {
	setVariables(host.Variables)
	createLinks(host.Links)
	execCommands(host.Commands)
	setDefaults(host.Defaults)
}

func setVariables(vars config.Variables) {
	if len(vars) > 0 {
		fmt.Printf("Setting variables (%d stages):\n", len(vars))
	}
	sw := int(math.Log10(float64(len(vars))))
	for i, stage := range vars {
		if len(stage) > 0 {
			fmt.Printf("[stage %[1]*[2]d/%[1]*[3]d]: %d variables\n", sw, i+1, len(vars), len(stage))
		}
		vw := int(math.Log10(float64(len(stage))))
		for varName, variable := range stage {
			fmt.Printf("[%[1]*[2]d/%[1]*[3]d]: %s=", vw, i+1, len(stage), varName)
			if variable.Command != nil {
				fmt.Printf("$(%s) -> ", variable.Command.String)
				var out buffer.Buffer
				variable.Command.Stdout = &out
				if err := variable.Command.Run(); err != nil {
					fmt.Println(err)
					continue
				}
				varVal := strings.TrimSpace(out.String())
				variable.Value = &varVal
			}
			if variable.Value != nil {
				fmt.Println(*variable.Value)
				if err := os.Setenv(varName, *variable.Value); err != nil {
					fmt.Println(err)
					continue
				}
				continue
			}
		}
	}
}

func createLinks(links config.Links) {
	fmt.Printf("Creating symlinks (%d):\n", len(links))
	lw := int(math.Log10(float64(len(links))))
	for i, l := range links {
		fmt.Printf("[%[1]*[2]d/%[1]*[3]d]: %s <- %s: ", lw, i+1, len(links), l.Target.Original, l.Source.Original)
		st, err := l.Link()
		if err != nil {
			fmt.Println(err)
			continue
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
}

func execCommands(cmds config.Commands) {
	nw := int(math.Log10(float64(len(cmds))))
	if len(cmds) > 0 {
		fmt.Printf("Executing commands (%d):\n", len(cmds))
	}
	for i, c := range cmds {
		fmt.Printf("[%[1]*[2]d/%[1]*[3]d]: %[4]s\n", nw, i+1, len(cmds), c.String)
		c.Stdin = nil
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr

		if err := c.Run(); err != nil {
			fmt.Printf("An error occured while running: %s\n", err)
		}
	}

}

func setDefaults(d defaults.Defaults) {
	fmt.Println("Setting defaults:")

	if len(d.Globals) > 0 {
		fmt.Printf("GLOBAL (%d):\n", len(d.Globals))
	}
	for keyName, key := range d.Globals {
		cmdStr := fmt.Sprintf(
			"defaults write -globalDomain %s %s",
			keyName,
			key.Value.String())
		fmt.Printf("[[%s]]: %s\n", keyName, cmdStr)
		cmd := new(config.Command)
		cmd.String = cmdStr

		cmd.Stdin = nil
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Println(err)
		}
	}

	for typ, domains := range map[string]defaults.Domains{
		" -app": d.Apps,
		"":      d.Domains,
	} {
		for domainName, domain := range domains {
			fmt.Printf("%s (%d):\n", domainName, len(domain))
			for keyName, key := range domain {
				cmdStr := fmt.Sprintf(
					"defaults write%s %s %s %s",
					typ,
					domainName,
					keyName,
					key.Value.String())
				fmt.Printf("[[%s]]: %s\n", keyName, cmdStr)
				cmd := new(config.Command).WithString(cmdStr)

				cmd.Stdin = nil
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				if err := cmd.Run(); err != nil {
					fmt.Println(err)
				}
			}
		}
	}

}
