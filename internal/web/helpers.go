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
	filepath := path.Join(c.cwd, "files", fileName)

	if err := os.MkdirAll("files", 0755); err != nil {
		return "", err
	}

	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return "", err
	}

	return filepath, nil
}
