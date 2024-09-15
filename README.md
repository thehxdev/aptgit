# aptgit
aptgit is a package manager for Github releases.

> [!WARNING]
> This project is not even in alpha stage and this is only a prototype.

## Build
First install GNU make and Go compiler. Then build aptgit:
```bash
make
```

## Usage
Provide a config file and a package's definition to `aptgit` executable to download the specified file:
```bash
# First make required directories
mkdir -p ~/.aptgit/{gpkgs,downloads,installs,bin}

# Downloaded files will be saved to `~/.aptgit/downloads` directory
aptgit -c ./src/config.json -def ./gpkgs/sing-box.json
```
