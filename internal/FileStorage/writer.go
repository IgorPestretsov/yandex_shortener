package FileStorage

import (
	"encoding/gob"
	"os"
)

type writer struct {
	file    *os.File
	encoder *gob.Encoder
}

func NewWriter(filename string) (*writer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &writer{file: file, encoder: gob.NewEncoder(file)}, nil
}

func (p *writer) WriteEvent(data map[string]map[string]string) error {
	return p.encoder.Encode(data)
}

func (p *writer) Close() error {
	return p.file.Close()
}
