# ccode

**ccode** is a lightweight command-line tool for combining text files into a single output file.  
It supports `.ccodeignore` rules (similar to `.gitignore`) to exclude files or directories,  
and provides options for safe overwriting, dry-run previews, and stripping empty lines.

---

## Features

- 📂 Recursively collect and combine text files from a directory.
- 📝 Supports `.ccodeignore` files for flexible file exclusion.
- 🚀 CLI flags for customization:
  - `--force` – overwrite output if it already exists.
  - `--dry-run` – preview which files would be combined without writing output.
  - `--strip-empty` – remove empty lines from file contents.
  - `--output, -o` – specify output file path (default: `combined.txt`).
- 🔍 Marks empty files with `[empty file]`.
- 🧪 Tested with unit and integration tests.
- 📦 Prebuilt binaries for Linux, macOS, and Windows (via GitHub Releases).

---

## Installation

### Using Go

```bash
go install github.com/x45iq/ccode@latest
```

### From Source

```bash
git clone https://github.com/x45iq/ccode.git
cd ccode
make build
./bin/ccode --help
```

### Prebuilt Binaries

Download the latest release from the [Releases page](https://github.com/x45iq/ccode/releases).

---

## Usage

Combine all text files in the current directory into `combined.txt`:

```bash
ccode
```

Specify a custom output file:

```bash
ccode --output myfiles.txt
```

Preview without writing output:

```bash
ccode --dry-run
```

Force overwrite if output already exists:

```bash
ccode --force --output result.txt
```

Remove empty lines from files:

```bash
ccode --strip-empty
```

---

## Ignore Rules (`.ccodeignore`)

You can exclude files and directories using `.ccodeignore` files placed anywhere in your project.
The syntax follows `.gitignore` conventions:

```
# Ignore all log files
*.log

# Ignore a specific file
secret.txt

# Keep a file explicitly
!keep.txt

# Ignore files in a subdirectory
subdir/*.tmp
```

---

## Development

Run linting, tests, and build:

```bash
make lint
make test
make build
```

Clean build artifacts:

```bash
make clean
```

---

## License

Unlicense license – see [LICENSE](LICENSE).
