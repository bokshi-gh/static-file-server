package main

import (
    "bufio"
    "flag"
    "fmt"
    "log"
    "mime"
    "net"
    "os"
    "path/filepath"
    "strings"
)

func main() {
    rootPath := flag.String("root", ".", "Root directory to serve files from")
    port := flag.Int("port", 8000, "Port to run the server on")
    flag.Parse()

    absRoot, err := filepath.Abs(*rootPath)
    if err != nil {
        log.Fatalf("Error resolving path: %v", err)
    }

    info, err := os.Stat(absRoot)
    if os.IsNotExist(err) || !info.IsDir() {
        log.Fatalf("Root path does not exist or is not a directory: %s", absRoot)
    }

    addr := fmt.Sprintf(":%d", *port)
    ln, err := net.Listen("tcp", addr)
    if err != nil {
        log.Fatalf("Error listening on %s: %v", addr, err)
    }
    defer ln.Close()

    fmt.Println("Static File Server listening on", addr)

    for {
        conn, err := ln.Accept()
        if err != nil {
            fmt.Println("Error accepting connection:", err)
            continue
        }
        go handleConnection(conn, absRoot)
    }
}

func handleConnection(conn net.Conn, absRoot string) {
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
        fmt.Printf("Bad request from %s\n", conn.RemoteAddr())
        return
    }

    method, reqPath, version := parts[0], parts[1], parts[2]

    if method != "GET" {
        conn.Write([]byte(version + " 405 Method Not Allowed\r\n\r\n"))
        fmt.Printf("Method not allowed: %s from %s\n", method, conn.RemoteAddr())
        return
    }

    reqPath = filepath.Clean(reqPath)
    fullPath := filepath.Join(absRoot, reqPath)

    if !strings.HasPrefix(fullPath, absRoot) {
        conn.Write([]byte(version + " 403 Forbidden\r\n\r\n"))
        fmt.Printf("Forbidden access attempt: %s from %s\n", fullPath, conn.RemoteAddr())
        return
    }

    info, err := os.Stat(fullPath)
    if err != nil {
        if os.IsNotExist(err) {
            conn.Write([]byte(version + " 404 Not Found\r\n\r\n"))
            fmt.Printf("File not found: %s requested by %s\n", fullPath, conn.RemoteAddr())
        } else {
            conn.Write([]byte(version + " 500 Internal Server Error\r\n\r\n"))
            fmt.Printf("Internal error serving %s to %s\n", fullPath, conn.RemoteAddr())
        }
        return
    }

    if info.IsDir() {
        indexPath := filepath.Join(fullPath, "index.html")
        body, err := os.ReadFile(indexPath)
        if err == nil {
            sendFile(conn, version, indexPath, body)
            fmt.Printf("Served %s to %s\n", indexPath, conn.RemoteAddr())
            return
        }

        entries, err := os.ReadDir(fullPath)
        if err != nil {
            conn.Write([]byte(version + " 500 Internal Server Error\r\n\r\n"))
            fmt.Printf("Internal error reading directory %s to %s\n", fullPath, conn.RemoteAddr())
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
        conn.Write([]byte(fmt.Sprintf("%s 200 OK\r\nContent-Length: %d\r\nContent-Type: text/html\r\n\r\n%s",
            version, len(html), html)))
        fmt.Printf("Served directory listing %s to %s\n", fullPath, conn.RemoteAddr())
        return
    }

    body, err := os.ReadFile(fullPath)
    if err != nil {
        conn.Write([]byte(version + " 500 Internal Server Error\r\n\r\n"))
        fmt.Printf("Internal error serving %s to %s\n", fullPath, conn.RemoteAddr())
        return
    }

    sendFile(conn, version, fullPath, body)
    fmt.Printf("Served %s to %s\n", fullPath, conn.RemoteAddr())
}

func sendFile(conn net.Conn, version, path string, body []byte) {
    ext := filepath.Ext(path)
    mimeType := mime.TypeByExtension(ext)
    if mimeType == "" {
        mimeType = "application/octet-stream"
    }

    headers := fmt.Sprintf("%s 200 OK\r\nContent-Length: %d\r\nContent-Type: %s\r\n\r\n",
        version, len(body), mimeType)
    conn.Write([]byte(headers))
    conn.Write(body)
}
