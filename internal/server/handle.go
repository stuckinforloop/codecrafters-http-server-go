package server

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func handle(conn net.Conn, dir string) error {
	defer conn.Close()

	b := make([]byte, 1024)
	if _, err := conn.Read(b); err != nil {
		return fmt.Errorf("read from connection: %w", err)
	}
	lines := strings.Split(string(b), "\r\n")
	if len(lines) < 3 {
		return fmt.Errorf("invalid request")
	}
	path, ok := getPath(lines[0])
	if !ok {
		if path == "" {
			return fmt.Errorf("invalid path: %s", path)
		}

		msg := "HTTP/1.1 404 Not Found\r\n\r\n"
		if _, err := conn.Write([]byte(msg)); err != nil {
			return fmt.Errorf("write to connection: %w", err)
		}

		return nil
	}

	var msg string
	switch path {
	case "/":
		msg = "HTTP/1.1 200 OK\r\n\r\n"
	case "/user-agent":
		var userAgent string
		for _, l := range lines {
			if strings.HasPrefix(l, "User-Agent") {
				userAgent = strings.TrimPrefix(l, "User-Agent: ")
			}
		}

		msg = genMessage(userAgent)
	default:
		if strings.HasPrefix(path, "/files") {
			filename := strings.TrimPrefix(path, "/files/")
			msg = genFileMessage(conn, dir, filename)
		} else {
			msg = genMessage(path)
		}
	}

	if _, err := conn.Write([]byte(msg)); err != nil {
		return fmt.Errorf("write to connection: %w", err)
	}

	return nil
}

func getPath(requestLine string) (string, bool) {
	parts := strings.Split(requestLine, " ")
	if len(parts) != 3 {
		return "", false
	}

	// get the second part of the request line
	path := parts[1]

	if strings.HasPrefix(path, "/files") {
		return path, true
	}

	// check if path starts with /echo/
	str := strings.TrimPrefix(path, "/echo/")
	if str != path {
		return str, true
	}

	if path == "/" || path == "/user-agent" {
		return path, true
	}

	return path, false
}

func genMessage(body string) string {
	msg := "HTTP/1.1 200 OK\r\n"
	msg += "Content-Type: text/plain\r\n"
	msg += fmt.Sprintf("Content-Length: %d\r\n\r\n", len(body))
	msg += body

	return msg
}

func genFileMessage(conn net.Conn, dir, filename string) string {
	filePath := filepath.Join(dir, filename)
	_, err := os.Stat(filePath)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return "HTTP/1.1 404 Not Found\r\n\r\n"
	}

	body, _ := os.ReadFile(filePath)

	msg := "HTTP/1.1 200 OK\r\n"
	msg += "Content-Type: application/octet-stream\r\n"
	msg += fmt.Sprintf("Content-Length: %d\r\n\r\n", len(body))
	msg += string(body)

	return msg
}
