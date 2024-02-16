package server

import (
	"fmt"
	"os"
)

func readFile(filePath string) ([]byte, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	body, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	return body, nil
}

func createFile(filePath string, body []byte) error {
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	if _, err := f.Write(body); err != nil {
		return fmt.Errorf("write to file: %w", err)
	}

	return nil
}
