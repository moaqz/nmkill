# nmkill

A tool to find and delete node_modules directories. It helps recover disk space by identifying and removing unused dependency folders.

## Installation

Download the latest binary for your platform from the [releases page](https://github.com/moaqz/nmkill/releases/latest)

For example, on Linux:

```bash
curl -fsSL https://github.com/moaqz/nmkill/releases/download/v0.1.0/nmkill_0.1.0_linux_amd64 -o nmkill
chmod u+x nmkill
./nmkill
```

## Usage

Run `nmkill` in the directory you want to scan:

```bash
nmkill
```

By default, it scans the current working directory. To scan a different path, use the `-d` flag:

```bash
nmkill -d /
```

### Interactive mode

After scanning, **nmkill** displays a numbered list of found node_modules folders with their size and last modification time.

You can then enter:

- A single number (e.g. `3`)
- A comma-separated list (e.g. `1,3,5`)
- The word `all` to select all
- `q` to quit

> [!WARNING]
> The tool will ask for confirmation before deleting any folder.


