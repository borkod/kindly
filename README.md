<div align="center">
	<br>
	<img src="www/img/kindly.png" alt="Logo" width="200">
	<br>
</div>

# Kindly

[![Release](https://img.shields.io/github/v/release/borkod/kindly?sort=semver)](https://github.com/borkod/kindly/releases/latest)
![Build Status](https://github.com/borkod/kindly/workflows/build/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/borkod/kindly?style=flat-square)](https://goreportcard.com/report/github.com/borkod/kindly)
![GitHub](https://img.shields.io/github/license/borkod/kindly)
[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg?style=flat-square)](https://www.paypal.me/borkodj)
[![Buy me a coffee](https://img.shields.io/badge/buy%20me-a%20coffee-orange.svg?style=flat-square)](https://www.buymeacoffee.com/borkod)

[Kindly](https://kindly.sh/) is a free and open-source software package management system that simplifies the installation of software on Linux..

## Documentation

All documentation at [kindly.sh](https://kindly.sh).

## Install

### Download

1. Download a prebuilt executable binary for your operating system from the [GitHub releases page](https://github.com/borko/kindly/releases).
2. Unzip the archive and place the executable binary wherever you would like to run it from. Additionally consider adding the location directory in the `PATH` variable if you would like the `kindly` command to be available everywhere.

### Homebrew

TODO 

```sh
brew install kindly
```

### Compile

**Clone**

```sh
git clone https://github.com/borkod/kindly
```

**Build using make**

TODO

```sh
make build
```

**Build using go**

```sh
cd kindly
go build .
```

## Usage

```
Usage:
  kindly [command]

Available Commands:
  check       Check if a package is available.
  help        Help about any command
  install     Installs one or many packages.
  list        Lists available packages.
  remove      Removes a previously installed package.
  template    Generate a Kindly YAML spec template for a GitHub repo.

Flags:
      --Arch string               Architecture (default is current architecture)
      --ManifestDir string        Default kindly manifests directory (default is $HOME/.kindly/manifests/)
      --OS string                 Operating System (default is current OS)
      --OutBinDir string          Default binary file output directory (default is $HOME/.kindly/bin/)
      --OutCompletionDir string   Default completions file output directory (default is $HOME/.kindly/completion/)
      --OutManDir string          Default man pages output directory (default is $HOME/.kindly/man/)
      --Source string             Source of packages (default is TODO ??? (default "https://b3o-test-bucket.s3.ca-central-1.amazonaws.com/")
      --completion string         Completion shell setting (default "bash")
      --config string             config file (default is $HOME/.kindly/.kindly.yaml)
  -h, --help                      help for kindly
  -v, --verbose                   Verbose output
      --version                   version for kindly
```
## Roadmap / TODO

- Testing
- Add more packages
- `Install` command:
	- If user installs a new version of a package that has less files or different file names than a previously installed version - remove may not properly remove all files as the package manifest will be rewritten.
	- Solution: uninstall previous version first?
- Add `Update` command
	- Updates all installed packages if new version available
- Add command to list locally installed packages

## Go Package

TODO

## Development

TODO

## Credit

TODO

## License

Apache-2.0