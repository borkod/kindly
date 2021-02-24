<head>
<link rel="preconnect" href="https://fonts.gstatic.com">
<link href="https://fonts.googleapis.com/css2?family=Molle:ital@1&display=swap" rel="stylesheet">
</head>
<div align="center">
	<br>
	<p style="font-family:Molle; font-size: 40px">Kindly</p> 
	<br>
</div>

# Kindly

[![Release](https://img.shields.io/github/v/release/borkod/kindly?sort=semver)](https://github.com/borkod/kindly/releases/latest)
![Build Status](https://github.com/borkod/kindly/workflows/build/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/borkod/kindly?style=flat-square)](https://goreportcard.com/report/github.com/borkod/kindly)
![GitHub](https://img.shields.io/github/license/borkod/kindly)
[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg?style=flat-square)](https://www.paypal.me/borkodj)
[![Buy me a coffee](https://img.shields.io/badge/buy%20me-a%20coffee-orange.svg?style=flat-square)](https://www.buymeacoffee.com/borkod)

[Kindly](https://kindly.sh/) install tool.

## Documentation

All documentation at [kindly.sh](https://kindly.sh).

## Install

### Download

1. Download a prebuilt executable binary for your operating system from the [GitHub releases page](https://github.com/borko/kindly/releases).
2. Unzip the archive and place the executable binary wherever you would like to run it from. Additionally consider adding the location directory in the `PATH` variable if you would like the `kindly` command to be available everywhere.

### Homebrew

```sh
brew install kindly
```

### Compile

**Clone**

```sh
git clone https://github.com/borkod/kindly
```

**Build using make**

```sh
make build
```

**Build using go**

```sh
cd cmd/kindly
go build .
```

## Usage

```
Usage:
  kindly [command]

Available Commands:
  check       A brief description of your command
  help        Help about any command
  install     A brief description of your command

Flags:
      --OutBinDir string          Default binary file output directory (default is $HOME/.local/bin/)
      --OutCompletionDir string   Default Completions file output directory (default is $HOME/.local/completion/)
      --OutManDir string          Default Man Pages output directory (default is $HOME/.local/man/)
      --completion string         Completion shell setting (default "bash")
      --config string             config file (default is $HOME/.kindly/.kindly.yaml)
  -h, --help                      help for kindly
      --unique-directory          write files into unique directory (default is false)
  -v, --verbose                   Verbose output
      --version                   version for kindly
```

## Go Package

TODO

## Development

TODO

## Credit

TODO

## License

Apache-2.0