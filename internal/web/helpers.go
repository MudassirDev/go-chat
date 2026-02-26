package web

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/MudassirDev/go-chat/db/database"
	"github.com/google/uuid"
)

func (c *APIConfig) saveAudio(data []byte) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	id := uuid.New().String()
	fileName := fmt.Sprintf("%v.webm", id)
	filepath := path.Join("files", fileName)
	fullPath := path.Join(cwd, filepath)

	if err := os.MkdirAll("files", 0755); err != nil {
		return "", err
	}

	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return "", err
	}

	return filepath, nil
}

func getUserFromContext(r *http.Request) (*database.GetUserWithIDRow, error) {
	err := errors.New("no user")
	rawUser := r.Context().Value(AUTH_KEY)
	if rawUser == nil {
		return nil, err
	}

	user, ok := rawUser.(database.GetUserWithIDRow)
	if !ok {
		return nil, err
	}

	return &user, nil
}
