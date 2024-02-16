package server

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Method string

const (
	GET  Method = "GET"
	POST Method = "POST"
)

type Request struct {
	Body    []byte
	Headers map[string]string
	Method  Method
	Target  string
	Version string
}

// ParseRequest parses an HTTP request
func ParseRequest(b []byte) (*Request, error) {
	// parse request line
	reader := bufio.NewReader(bytes.NewReader(b))
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("read request line: %w", err)
	}

	fields := strings.Fields(requestLine)
	if len(fields) < 3 {
		return nil, fmt.Errorf("invalid request line: %s", requestLine)
	}
	method := fields[0]
	target := fields[1]
	version := fields[2]

	// parse headers
	headers := make(map[string]string)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("read request headers: %w", err)
		}
		if line == "\r\n" {
			break
		}

		splits := strings.SplitN(line, ":", 2)
		headers[splits[0]] = formatHeader(splits[1])
	}

	// parse body
	var contentLength int
	contentLengthHeader := headers["Content-Length"]
	if contentLengthHeader != "" {
		contentLength, err = strconv.Atoi(headers["Content-Length"])
		if err != nil {
			return nil, fmt.Errorf("parse content length: %w", err)
		}
	}
	body := make([]byte, contentLength)
	if _, err := io.ReadFull(reader, body); err != nil && err != io.EOF {
		return nil, fmt.Errorf("read request body: %w", err)
	}

	req := &Request{
		Body:    body,
		Headers: headers,
		Method:  Method(method),
		Target:  target,
		Version: version,
	}

	return req, nil
}

func formatHeader(val string) string {
	val = strings.TrimSpace(val)
	val = strings.Trim(val, "\r\n")
	return val
}
