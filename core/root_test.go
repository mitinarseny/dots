package core

import (
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestUnmarshalRoot(t *testing.T) {
	r := require.New(t)
	data := `
default:
  links:
    ~/.zshrc: .zshrc
    ~/.gitexcludes:
      path: .gitexcludes
      force: true
    ~/.zsh/:
      path: .zsh/*
      force: true

  commands:
    - command: antibody bundle < zsh_plugins.txt > ~/.zsh_plugins.sh
      description: load zsh plugins
    - command: bat cache --source .config/bat/ --build
      description: init bat cache

macos:
  extends: default
  variables:
    - HOME:
        command: echo ~
  links:
    ~/Library/Application Support/Sublime Text 3/Packages/User/Package Control.sublime-settings:
      path: subl/Package Control.sublime-settings
      force: true
    ~/Library/Application Support/Sublime Text 3/Packages/User/Preferences.sublime-settings:
      path: subl/Preferences.sublime-settings
      force: true
  commands:
    - defaults write -app iTerm PrefsCustomFolder $HOME/.iterm2_profile
    - defaults write -app iTerm LoadPrefsFromCustomFolder -bool true
    - command: defaults write -app Safari AllowJavaScriptFromAppleEvents 1
      description: fix BeardedSpice
`
	dflt := Host{
		Name: "default",
		Links: &Links{
			"~/.zshrc": {
				Target: "~/.zshrc",
				Path:   ".zshrc",
				Force:  false,
			},
			"~/.gitexcludes": {
				Target: "~/.gitexcludes",
				Path:   ".gitexcludes",
				Force:  true,
			},
			"~/.zsh/": {
				Target: "~/.zsh/",
				Path:   ".zsh/*",
				Force:  true,
			},
		},
		Commands: &Commands{
			{
				String:      "antibody bundle < zsh_plugins.txt > ~/.zsh_plugins.sh",
				Description: sp("load zsh plugins"),
			},
			{
				String:      "bat cache --source .config/bat/ --build",
				Description: sp("init bat cache"),
			},
		},
	}
	expected := Config{
		"default": &dflt,
		"macos": {
			Name:    "macos",
			Extends: &dflt,
			Variables: &Variables{
				{
					"HOME": {
						Name: "HOME",
						Command: &Command{
							String: "echo ~",
						},
					},
				},
			},
			Links: &Links{
				"~/Library/Application Support/Sublime Text 3/Packages/User/Package Control.sublime-settings": {
					Target: "~/Library/Application Support/Sublime Text 3/Packages/User/Package Control.sublime-settings",
					Path:   "subl/Package Control.sublime-settings",
					Force:  true,
				},
				"~/Library/Application Support/Sublime Text 3/Packages/User/Preferences.sublime-settings": {
					Target: "~/Library/Application Support/Sublime Text 3/Packages/User/Preferences.sublime-settings",
					Path:   "subl/Preferences.sublime-settings",
					Force:  true,
				},
			},
			Commands: &Commands{
				{
					String: "defaults write -app iTerm PrefsCustomFolder $HOME/.iterm2_profile",
				},
				{
					String: "defaults write -app iTerm LoadPrefsFromCustomFolder -bool true",
				},
				{
					String:      "defaults write -app Safari AllowJavaScriptFromAppleEvents 1",
					Description: sp("fix BeardedSpice"),
				},
			},
		},
	}
	var root Config
	err := yaml.Unmarshal([]byte(data), &root)

	r.NoError(err)
	r.Equal(expected, root)
}

func sp(s string) *string {
	return &s
}
