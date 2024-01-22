# VanityForge
![VanityForge](https://vhs.charm.sh/vhs-328VUMdxvRha1adlp4fTx3.gif)
## Introduction
VanityForge is a powerful CLI tool designed for generating blockchain Vanity addresses with efficiency and ease. It supports multiple networks offering a wide range of customization options for address generation. (See [_Supported Chains_](#supported-chains))

## Key Features
- **Generate Bech32 Vanity Addresses**: Create personalized addresses with specific patterns.
- **Multi-Core Support**: Utilizes all CPU cores for faster generation.
- **Customizable Address Patterns**: Specify substrings for addresses to start with, end with, or contain.
- **Minimum Character Requirements**: Set required minimum letters or digits in addresses.
- **Cross-Platform Compatibility**: Binaries available for Linux, macOS, and Windows.
- **Generate Bech16 EVM Vanity Addresses**

## Getting Started

### Installation
Download the latest binary releases from the [_Releases_](https://github.com/01node/vanity-forge/releases) page. .

If you don't want to use the binaries you can also build from source:

Required go version: 1.20
```bash
git clone https://github.com/01node/vanity-forge vanity-forge
cd vanity-forge
go build .
```

### Usage
Interactive usage is available by running `./vanity-forge` without any flags.

![Demo](https://vhs.charm.sh/vhs-NrMqJNbxgzBN7szFJA1yu.gif)


## Advanced Usage
```bash
Usage of ./vanity-forge:
  -n, --accounts-number int   Amount of accounts you need
  -c, --chain string          Chain selector string
  -d, --digits int            Amount of digits (0-9) that the address must contain
  -l, --letters int           Amount of letters (a-z) that the address must contain
  -m, --mode string           Matcher mode (contains, starts-with, ends-with, regex)
  -s, --search string         Search string
  -v, --verbose               Verbose output
```
![Advanced demo](https://vhs.charm.sh/vhs-2v4VLUIOfeCaiu8Lz4OpU3.gif)

## Supported Chains
- Cosmos
- Celestia
- dYdX
- Berachain

## License
This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Acknowledgements
- https://github.com/hukkin/cosmosvanity used as inspiration for the project.
