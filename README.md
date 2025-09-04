# GoServe

## Overview

GoServe is a simple TCP-based static file server written in Go. It allows you to serve files and directories over HTTP with minimal configuration. It is designed for quick local development, testing, or lightweight file sharing.

## Features

- Serve static files over HTTP from any directory
- Automatic MIME type detection based on file extension
- Directory listing when `index.html` is missing
- Simple CLI with configurable root directory and port
- Lightweight and easy to build/run

## Tools and Technology used

- [Go](https://golang.org/) (Golang programming language)
- Standard Go packages: `net`, `os`, `path/filepath`, `bufio`, `mime`

## Getting Started

### Platforms

This project supports the following platforms:

- Linux
- macOS
- Windows

### Requirements

- Go 1.21 or higher installed
- Basic terminal/command-line usage knowledge

### Installation

1. **Clone the repository:**
    ```sh
    git clone https://github.com/{{your-username}}/goserve.git
    cd goserve
    ```

2. **Build the project:**
    ```sh
    go build -o goserve
    ```

3. **Run the server:**
    ```sh
    ./goserve --root ./public --port 8000
    ```
    > Note: Replace `./public` with the directory you want to serve.

## Configuration

- Set the root directory to serve files from using `--root` (default is current directory)
- Set the port to listen on using `--port` (default is 8000)
- Check version using `--version`
- Help is available with `-h` or `--help`

## Usage

- Open a browser and navigate to `http://localhost:8000/` to view your files
- Example command:
    ```sh
    ./goserve --root ./myfiles --port 9000
    ```
- Logs in terminal will show which files are served to which client

## Contributing

Contributions are welcome! Please fork the repo, make changes, and submit pull requests. Open issues for bugs or feature requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
