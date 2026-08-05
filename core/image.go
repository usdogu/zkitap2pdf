package core

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"io"
	"strings"

	"github.com/h2non/filetype"
	"golang.org/x/image/webp"
)

type Image struct {
	Name string
	Type string
	Data []byte
}

func convertImage(name string, data []byte) (Image, error) {
	reader := bytes.NewReader(data)

	kind, err := filetype.MatchReader(reader)
	if err != nil {
		return Image{}, err
	}

	reader.Seek(0, io.SeekStart)

	switch kind.MIME.Value {
	case "image/webp":
		img, err := webp.Decode(reader)
		if err != nil {
			return Image{}, err
		}

		var buf bytes.Buffer
		if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90}); err != nil {
			return Image{}, err
		}

		return Image{
			Name: name,
			Type: "jpg",
			Data: buf.Bytes(),
		}, nil

	case "image/jpeg", "image/png":
		return Image{
			Name: name,
			Type: strings.TrimPrefix(kind.MIME.Value, "image/"),
			Data: data,
		}, nil

	default:
		return Image{}, fmt.Errorf("unsupported image: %s", name)
	}
}
