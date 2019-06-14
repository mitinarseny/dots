<p align="center">
    <a href="https://github.com/mitinarseny/dots">
        <img src="assets/logo.png" alt="dots logo" height="140" />
    </a>
    <h1 align="center">dots</h1>
    <p align="center">Delivery tool for your <code>.dotfiles</code></p>
    <p align="center">
      <a href="https://github.com/mitinarseny/dots/releases/latest"><img alt="Release" src="https://img.shields.io/github/release/mitinarseny/dots.svg?style=flat-square"></a>
      <a href="/LICENSE.md"><img alt="Software License" src="https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square"></a>
      <a href="https://travis-ci.org/mitinarseny/dots"><img alt="TravisCI" src="https://img.shields.io/travis/mitinarseny/dots/master.svg?style=flat-square"></a>
      <a href="https://codecov.io/gh/mitinarseny/dots"><img alt="Codecov branch" src="https://img.shields.io/codecov/c/github/mitinarseny/dots/master.svg?style=flat-square"></a> 
    </p>
</p>

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
An example config is available [here](https://github.com/mitinarseny/dotfiles/blob/master/.dots.yaml).



