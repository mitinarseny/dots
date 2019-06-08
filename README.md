<p align="center">
    <a href="https://github.com/mitinarseny/dots">
        <img src="https://rzhao.io/img/dotfiles/dotfiles.png" alt="dots logo" />
    </a>
</p>

# dots

[![Build Status](https://travis-ci.org/mitinarseny/dots.svg?branch=master)](https://travis-ci.org/mitinarseny/dots)
[![Coverage Status](https://coveralls.io/repos/github/mitinarseny/dots/badge.svg?branch=master)](https://coveralls.io/github/mitinarseny/dots?branch=master)

Delivery tool for your `.dotfiles`.

## Install

### macOS

```bash
brew install mitinarseny/tap/dots
```

## Usage
```bash
dots up [hostName] [-c path/to/config.yaml]
```


## Config
Config file is a `.yaml` file with the following structure:

```yaml
<HOST>
hosts:
  hostName: <HOST>
```
### Hosts
Each `<HOST>` has the following structure:

```yaml
variables:
  - name1: value1
  - name2: $name1/value2 # use variables defined above
  # ...
  
links:
  /absolute/target/path: relative/path
  ~/relative/to/home/target/path: /absolute/path
  ~/another/target/path:
    path: some/path
    force: false
  ~/target/path/using/variable/$name1:
    path: source/path/using/variable/$name2:
    force: true
  # ...

commands:
  - echo "commands"
  - echo "are"
  - echo "executed"
  - echo "in consecutive order"
  # ...

defaults:
  apps:
    AppName:
      key: string_value
      anotherKey:
        type: bool
        value: true
      andAnotherKey:
        type: array
        value:
          - value1
          - value2
          # ...
      yetAnotherKey:
        type: dict
        value:
          inner_key1: value1
          innrer_key2: value2
          # ...
    AnotherAppName:
      key: value
      # ...
  domains:
    some.app.domain:
      key: value
      # ...
  globals:
    key: value
```
An example config is available [here](example.config.yaml).



