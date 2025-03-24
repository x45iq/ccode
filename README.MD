# ccode

**ccode** is an open-source command-line tool written in Go that collects and consolidates code files from a specified directory into a single output file. It supports exclusion rules via a `.ccodeignore` file, enabling you to filter out files or directories you don’t want to include.

---

## Features

- **Recursive File Collection:** Walk through a target directory and collect all files.
- **Exclusion Support:** Use a `.ccodeignore` file (similar to `.gitignore`) to define patterns for files or directories to exclude.
- **Customizable Output:** Specify an output file name with the `-o` flag (default is `combined.txt`).
- **Force Overwrite:** Use the `--force` flag to overwrite an existing output file.
- **Simple and Lightweight:** Written in Go with minimal dependencies.

---

## Installation

To build and install **ccode**, ensure you have [Go](https://golang.org/) installed (version 1.24 or later). Then run:

```bash
git clone https://github.com/x45iq/ccode.git
cd ccode
go build -o ccode
```

You can also install it directly using `go install`:

```bash
go install github.com/x45iq/ccode@latest
```
---

## Usage

After installation, run **ccode** from the command line. The basic usage is:

```bash
ccode [flags] <target_directory>
```

### Options

- `-o, --output`: Specify the name of the output file.  
  _Example:_ `ccode -o mycode.txt /path/to/project`

- `--force`: Overwrite the output file if it already exists.  
  _Example:_ `ccode --force /path/to/project`

### Example

To collect code from a project directory while using a custom output file and forcing an overwrite, run:

```bash
ccode -o combined.txt --force /path/to/your/project
```

The tool will walk through the directory, apply any ignore rules defined in `.ccodeignore` files, and append the contents of each file (preceded by a header indicating the file path) into the specified output file.

---

## .ccodeignore

**ccode** uses a `.ccodeignore` file to determine which files and directories should be excluded from the output. The format is similar to a `.gitignore` file:

- **Comments:** Lines starting with `#` are ignored.
- **Negation:** Lines starting with `!` negate a pattern (i.e., include files previously excluded).
- **Wildcards:** Supports `*`, `?`, and `**` for matching patterns.

Example `.ccodeignore`:

```
# Ignore all log files
*.log

# Exclude node_modules directory
node_modules/

# But do include specific file even if it matches above rules
!node_modules/special-module/config.js
```

---

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests with enhancements or bug fixes.

1. Fork the repository.
2. Create a new branch for your feature or fix.
3. Commit your changes with clear commit messages.
4. Submit a pull request describing your changes.

---

## License

Distributed under the [MIT License](LICENSE).  

---
