package main

import (
	"fmt"
	"image"
	"os"

	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
)

type converter struct {
	Options *encoder.Options
}

func (conv converter) decode(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return img, nil
}

func (conv converter) encode(path string, img image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}
	defer file.Close()

	err = webp.Encode(file, img, conv.Options)
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}
	return nil
}

func (conv converter) Do(path string) error {
	img, err := conv.decode(path)
	if err != nil {
		return err
	}

	opath := webpPath(path)
	err = conv.encode(opath, img)
	if err != nil {
		return err
	}

	return nil
}
