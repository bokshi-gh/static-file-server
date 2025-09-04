package main

import (
    "flag"
    "path/filepath"
    "os"
    "log"
    "strings"
    "mime"
    "bufio"
    "fmt"
    "net"
)

func main() {
    rootPath := flag.String("root", ".", "Root path to serve files from")
    port := flag.Int("port", 8000, "Port to run the server on")
    flag.Parse()

    absPath, err := filepath.Abs(*rootPath)
    if err != nil {
        log.Fatalf("Error resolving path: %v", err)
    }

    info, err := os.Stat(absPath)
    if os.IsNotExist(err) {
        log.Fatalf("Path does not exist: %s", absPath)
    }
    if !info.IsDir() {
        log.Fatalf("Path is not a directory: %s", absPath)
    }

    addr := fmt.Sprintf(":%d", *port)

    ln, err := net.Listen("tcp", addr)
    if err != nil {
        panic(err)
    }
    defer ln.Close()
    fmt.Println("Static File Server listening on", addr)

    for {
        conn, err := ln.Accept()
        if err != nil {
            fmt.Println("Error accepting connection:", err)
            continue
        }
        go handleConnection(conn, absPath)
    }
}

func handleConnection(conn net.Conn, absPath string) {
    defer conn.Close()
    reader := bufio.NewReader(conn)

    line, err := reader.ReadString('\n')
    if err != nil {
        fmt.Println("Error reading:", err)
        return
    }

    line = strings.TrimSpace(line)

    parts := strings.Split(line, " ")
    if len(parts) < 3 {
        conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
        return
    }

    filePath := parts[1]
    filePath = filepath.Clean(filePath)
    fullPath := filepath.Join(absPath, filePath)
    
    body, err := os.ReadFile(fullPath)
    if err != nil {
        if os.IsNotExist(err) {
			resp := fmt.Sprintf("HTTP/1.1 404 Not Found\r\n\r\n")
			conn.Write([]byte(resp))
        } else {
			
			resp := fmt.Sprintf("HTTP/1.1 500 Internal Server Error\r\n\r\n")
			conn.Write([]byte(resp))
        }
        return
    }
    
    ext := filepath.Ext(filePath)
    mimeType := mime.TypeByExtension(ext)
    if mimeType == "" {
        mimeType = "application/octet-stream"
    }

    headers := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Length: %d\r\nContent-Type: %s\r\n\r\n", len(body), mimeType)
    conn.Write([]byte(headers))
    conn.Write(body)

    fmt.Printf("Served %s to %s\n", fullPath, conn.RemoteAddr())
}

