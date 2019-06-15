<p align="center">
    <a href="https://github.com/mitinarseny/dots">
        <img src="assets/logo.png" alt="dots logo" width="30%" />
    </a>
    <h1 align="center">dots</h1>
    <p align="center">Delivery tool for your <code>.dotfiles</code></p>
    <p align="center">
      <a href="https://github.com/mitinarseny/dots/releases/latest"><img alt="Release" src="https://img.shields.io/github/release/mitinarseny/dots.svg?style=flat-square"></a>
      <a href="/LICENSE.md"><img alt="Software License" src="https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square"></a>
      <a href="https://travis-ci.org/mitinarseny/dots"><img alt="TravisCI" src="https://img.shields.io/travis/mitinarseny/dots/master.svg?style=flat-square"></a>
      <a href="https://codecov.io/gh/mitinarseny/dots"><img alt="Codecov branch" src="https://img.shields.io/codecov/c/github/mitinarseny/dots/master.svg?style=flat-square"></a>
      <a href="http://godoc.org/github.com/mitinarseny/dots"><img alt="Go Doc" src="https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square"></a>
      <a href="https://saythanks.io/to/mitinarseny"><img alt="SayThanks.io" src="https://img.shields.io/badge/SayThanks.io-♥-1EAEDB.svg?style=flat-square"></a>
      <a href="https://houndci.com"><img alt="Reviewed by Hound" src="https://img.shields.io/badge/Reviewed_by-Hound-8E64B0.svg?style=flat-square&longCache=true"></a> 
    </p>
</p>

---

## Install

### macOS
```bash
brew install mitinarseny/tap/dots
```

## Usage
To set variables, create symlinks and execute commands listed in [config](#config) file run this: 
```bash
dots up
```

### Custom config file
By default [config](#config) file is `.dots.yaml` in current working directory.
If you want to specify other location, pass it to `-c` flag:
```bash
dots up -c custom/config.yaml
```
 
### Specify hostname
You are free to choose what host to up:
```bash
dots up hostname
```
The default hostname is `default`.

## Config
Config is a YAML file (defaults to `.dots.yaml` in current working directory).
An example config can be found in my [dotfiles](https://github.com/mitinarseny/dotfiles/blob/master/.dots.yaml).
  
On top-level config consists of mapping of hostnames to `host`s:
```yaml
hostname1:
  # ...
hostname2:
  # ...
```
If not specifying the hostname when running `dots up`, default will be `default`.

### `host` 
Each `host` has the following structure:
```yaml
hostname:
  extends: hostname
  variables:
    # ...
  links:
    # ...
  commands:
    # ...
```
Each of fields can be omitted. 

### `variables`
This section consists of list of stages, while each stage consists of environment variables.
```yaml
variables:
  - name1: value1
    name2: value2
  - name3: $name1/$name2
  - name4:
      command: echo value4
      description: this is equal to exporting name4=$(echo value4)
  - name5:
      command: echo "description can be omitted"
  # ...
```

### `links`
This section is a mapping of target paths to destination paths. Targets must be absolute. 
```yaml
links:
  ~/abs/target: rel/dest
  
  ~/another/abs/tagret:
    path: another/rel/dest
    force: true # ensures to recreate link if already exists
    
  # link directory
  ~/dir1: dir1
  
  # create in ~/dir2/ the same structure as in .dir/ and link all .sh files
  ~/dir2/: .dir/*.sh 
  # ...
``` 

### `commands`
This section is a list of commands to execute. Each command may have a `description`.
```yaml
commands:
  - echo command1
  - command: echo command2
    description: this
  # ... 
```
