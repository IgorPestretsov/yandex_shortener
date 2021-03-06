package filestorage

import (
	"encoding/gob"
	"os"
)

type reader struct {
	file    *os.File
	decoder *gob.Decoder
}

func NewReader(filename string) (*reader, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0444)
	if err != nil {
		return nil, err
	}
	return &reader{file: file, decoder: gob.NewDecoder(file)}, nil
}

func (c *reader) ReadData() (NestedMap, error) {
	data := make(map[string]map[string]string)
	if err := c.decoder.Decode(&data); err != nil {
		return data, err
	}
	return data, nil
}

func (c *reader) Close() error {
	return c.file.Close()
}
