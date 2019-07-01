package core

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
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
      path: .zsh/**
      force: true
      dirs: true

  commands:
    - command: antibody bundle < zsh_plugins.txt > ~/.zsh_plugins.sh
      description: load zsh plugins
    - command: bat cache --source .config/bat/ --build
      description: init bat cache

macos:
  extends: 
    - default
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
    - description: Setting defaults
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
				Path:   ".zsh/**",
				Force:  true,
				Dirs:   true,
			},
		},
		Commands: &Commands{
			{
				String:      sp("antibody bundle < zsh_plugins.txt > ~/.zsh_plugins.sh"),
				Description: sp("load zsh plugins"),
			},
			{
				String:      sp("bat cache --source .config/bat/ --build"),
				Description: sp("init bat cache"),
			},
		},
	}
	expected := Config{
		"default": &dflt,
		"macos": {
			Name: "macos",
			Extends: []*Host{
				&dflt,
			},
			Variables: &Variables{
				{
					"HOME": {
						Name: "HOME",
						Command: &Command{
							String: sp("echo ~"),
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
					Description: sp("Setting defaults"),
					Commands: []*Command{
						{
							String: sp("defaults write -app iTerm PrefsCustomFolder $HOME/.iterm2_profile"),
						},
						{
							String: sp("defaults write -app iTerm LoadPrefsFromCustomFolder -bool true"),
						},
						{
							String:      sp("defaults write -app Safari AllowJavaScriptFromAppleEvents 1"),
							Description: sp("fix BeardedSpice"),
						},
					},
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
