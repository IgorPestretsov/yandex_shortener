package storage

import (
	"encoding/gob"
	"os"
)

type consumer struct {
	file    *os.File
	decoder *gob.Decoder
}

func NewConsumer(filename string) (*consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0777)
	if err != nil {
		return nil, err
	}
	return &consumer{file: file, decoder: gob.NewDecoder(file)}, nil
}

func (c *consumer) ReadData() (map[string]string, error) {
	data := make(map[string]string)
	if err := c.decoder.Decode(&data); err != nil {
		return data, err
	}
	return data, nil
}

func (c *consumer) Close() error {
	return c.file.Close()
}
