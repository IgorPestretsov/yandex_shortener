package storage

import (
	"encoding/gob"
	"os"
)

type producer struct {
	file    *os.File
	encoder *gob.Encoder
}

func NewProducer(filename string) (*producer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}

	return &producer{file: file, encoder: gob.NewEncoder(file)}, nil
}

func (p *producer) WriteEvent(data map[string]string) error {
	return p.encoder.Encode(data)
}

func (p *producer) Close() error {
	return p.file.Close()
}
