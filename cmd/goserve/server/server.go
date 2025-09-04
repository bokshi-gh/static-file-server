package server

import (
	"bufio"
	"fmt"
	"mime"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func StartServer(rootPath string, port int, version string) {
	addr := fmt.Sprintf(":%d", port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	fmt.Println("GoServe version", version, "listening on port", port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go HandleClient(conn, rootPath)
	}
}

func HandleClient(conn net.Conn, rootPath string) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	line, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading from", conn.RemoteAddr(), ":", err)
		return
	}

	line = strings.TrimSpace(line)
	parts := strings.Split(line, " ")
	if len(parts) < 3 {
		conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
		return
	}

	method, reqPath, version := parts[0], parts[1], parts[2]
	if method != "GET" {
		conn.Write([]byte(version + " 405 Method Not Allowed\r\n\r\n"))
		return
	}

	reqPath = filepath.Clean(reqPath)
	fullPath := filepath.Join(rootPath, reqPath)
	if !strings.HasPrefix(fullPath, rootPath) {
		conn.Write([]byte(version + " 403 Forbidden\r\n\r\n"))
		return
	}

	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			conn.Write([]byte(version + " 404 Not Found\r\n\r\n"))
		} else {
			conn.Write([]byte(version + " 500 Internal Server Error\r\n\r\n"))
		}
		return
	}

	if info.IsDir() {
		indexPath := filepath.Join(fullPath, "index.html")
		body, err := os.ReadFile(indexPath)
		if err == nil {
			SendFile(conn, version, indexPath, body)
			return
		}

		entries, err := os.ReadDir(fullPath)
		if err != nil {
			conn.Write([]byte(version + " 500 Internal Server Error\r\n\r\n"))
			return
		}

		html := "<html><body><h1>Directory listing for " + reqPath + "</h1><ul>"
		for _, entry := range entries {
			name := entry.Name()
			if entry.IsDir() {
				name += "/"
			}
			html += fmt.Sprintf(`<li><a href="%s">%s</a></li>`, name, name)
		}
		html += "</ul></body></html>"

		conn.Write([]byte(fmt.Sprintf("%s 200 OK\r\nContent-Length: %d\r\nContent-Type: text/html\r\n\r\n%s", version, len(html), html)))
		return
	}

	body, err := os.ReadFile(fullPath)
	if err != nil {
		conn.Write([]byte(version + " 500 Internal Server Error\r\n\r\n"))
		return
	}

	SendFile(conn, version, fullPath, body)
}

func SendFile(conn net.Conn, version, path string, body []byte) {
	ext := filepath.Ext(path)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	headers := fmt.Sprintf("%s 200 OK\r\nContent-Length: %d\r\nContent-Type: %s\r\n\r\n",
		version, len(body), mimeType)
	conn.Write([]byte(headers))
	conn.Write(body)

	fmt.Printf("Served %s to %s\n", path, conn.RemoteAddr())
}

