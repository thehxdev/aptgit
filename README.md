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

## Todo!
- [ ] Implement `init`, `install`, `uninstall`, `upgrade`, `list-versions`, `latest-version`, `list-installed`, `switch`, `cleanup` sub-commands
- [X] Install and set custom version of a program
- [ ] Override aptgit and package parameters from command line
- [ ] Ensure all required directories exist before any processing
- [ ] Structured logging
- [ ] Better error handling
- [ ] Process multiple packages concurrently
- [ ] `aptgit.lock` file to keep metadata about installed packages (if needed - probably required for upgrading packages)
- [ ] Better naming, coding style and cleaner architecture
