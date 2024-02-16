package server

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func handle(conn net.Conn, dir string) error {
	b := make([]byte, 1024)
	if _, err := conn.Read(b); err != nil {
		return fmt.Errorf("read from connection: %w", err)
	}

	req, err := ParseRequest(b)
	if err != nil {
		log.Printf("parse request: %v", err)
		return fmt.Errorf("parse request: %w", err)
	}

	resp := Response{}
	resp.Body = make([]byte, 1024)
	resp.Headers = make(map[string]string)
	switch {
	case req.Target == "/":
		resp.StatusCode = 200
		resp.Version = req.Version

	case strings.HasPrefix(req.Target, "/echo"):
		body := strings.TrimPrefix(req.Target, "/echo/")

		resp.StatusCode = 200
		resp.Version = req.Version
		resp.Headers["Content-Type"] = "text/plain"
		resp.Headers["Content-Length"] = strconv.Itoa(len(body))
		resp.Body = []byte(body)

	case strings.HasPrefix(req.Target, "/files"):
		filename := strings.TrimPrefix(req.Target, "/files/")
		filePath := filepath.Join(dir, filename)

		// if method post, then create a file
		if req.Method == POST {
			body := bytes.Trim(req.Body, "\x00")
			if err := createFile(filePath, body); err != nil {
				return fmt.Errorf("create file: %w", err)
			}

			resp.StatusCode = 201
			resp.Version = req.Version
		} else if req.Method == GET {
			body, err := readFile(filePath)
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					resp.StatusCode = 404
					resp.Version = req.Version
					break
				}
				return fmt.Errorf("read file: %w", err)
			}

			resp.StatusCode = 200
			resp.Version = req.Version
			resp.Headers["Content-Type"] = "application/octet-stream"
			resp.Headers["Content-Length"] = strconv.Itoa(len(body))
			resp.Body = body
		}

	case req.Target == "/user-agent":
		body := req.Headers["User-Agent"]

		resp.StatusCode = 200
		resp.Version = req.Version
		resp.Headers["Content-Type"] = "text/plain"
		resp.Headers["Content-Length"] = strconv.Itoa(len(body))
		resp.Body = []byte(body)

	default:
		resp.StatusCode = 404
		resp.Version = req.Version

	}

	if err := WriteResponse(conn, resp); err != nil {
		log.Printf("write response: %v", err)
		return fmt.Errorf("write response: %w", err)
	}

	return nil
}
