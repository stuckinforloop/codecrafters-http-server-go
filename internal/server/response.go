package server

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
)

type StatusCode int

const (
	StatusOK       StatusCode = 200
	StatusCreated  StatusCode = 201
	StatusNotFound StatusCode = 400
)

type Response struct {
	Body       []byte
	Headers    map[string]string
	StatusCode StatusCode
	Version    string
}

func WriteResponse(conn net.Conn, resp Response) error {
	defer conn.Close()

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	// write status line
	statusCode := StatusCode(resp.StatusCode)
	if _, err := writer.WriteString(fmt.Sprintf("%s %s\r\n", resp.Version, statusCode.String())); err != nil {
		return fmt.Errorf("write status line: %w", err)
	}

	// write headers
	for k, v := range resp.Headers {
		if _, err := writer.WriteString(fmt.Sprintf("%s: %s\r\n", k, v)); err != nil {
			return fmt.Errorf("write header: %w", err)
		}
	}

	// separate headers and body
	if _, err := writer.WriteString("\r\n"); err != nil {
		return fmt.Errorf("break header from body: %w", err)
	}

	// write body
	if resp.Body != nil {
		if _, err := writer.Write(resp.Body); err != nil {
			return fmt.Errorf("write body: %w", err)
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("flush writer: %w", err)
	}

	if _, err := conn.Write(buf.Bytes()); err != nil {
		return fmt.Errorf("write response: %w", err)
	}

	return nil
}

func (sc StatusCode) String() string {
	switch sc {
	case 200:
		return "200 OK"
	case 201:
		return "201 Created"
	case 404:
		return "404 Not Found"
	default:
		return ""
	}
}
