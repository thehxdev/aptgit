# aptgit
aptgit is a package and version manager for Github releases.

> [!WARNING]
> This project is in alpha stage.

## Build
First install [GNU make](https://www.gnu.org/software/make/) and [Go compiler](https://go.dev). Then build aptgit:
```bash
make
```

## Usage

> [!NOTE]
> Since this project is a prototype, use it only for testing or in a testing environment.

> [!WARNING]
> The packages defined in [gpkgs](gpkgs) directory are only tested on Linux Mint 22 environment (other Linux systems must be OK) but macOS is not tested.

Provide a config file and a package's definition to `aptgit` executable to download the specified file:
```bash
# Make required directories
mkdir -p ~/.aptgit/{downloads,installs,bin}

# Copy package definitions to aptgit home
cp -r ./gpkgs ~/.aptgit

# Get help message
aptgit help
```

## Todo!
- [ ] Implement `init`, `install`, `uninstall`, `upgrade`, `list-versions`, `latest-version`, `list-installed`, `switch`, `cleanup` sub-commands
- [X] Install and set custom version of a program
- [ ] Override aptgit and package parameters from command line
- [x] Ensure all required directories exist before any processing
- [ ] Structured logging
- [ ] Better error handling
- [ ] Process multiple packages concurrently
- [X] `aptgit.lock` file to keep metadata about installed packages (if needed - probably required for upgrading packages)
- [ ] Better naming, coding style and cleaner architecture
- [ ] Cleanup the source code and seperate modules cause it's a mess
