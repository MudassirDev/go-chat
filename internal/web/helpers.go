package web

import (
	"fmt"
	"os"
	"path"

	"github.com/google/uuid"
)

func (c *apiConfig) saveAudio(data []byte) (string, error) {
	id := uuid.New().String()
	fileName := fmt.Sprintf("%v.webm", id)
	filepath := path.Join("files", fileName)
	fullPath := path.Join(c.cwd, filepath)

	if err := os.MkdirAll("files", 0755); err != nil {
		return "", err
	}

	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return "", err
	}

	return filepath, nil
}
