package resizer

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"

	"github.com/disintegration/imaging"
)

type Resizer struct {
}

func (r *Resizer) Resize(src []byte, width, height int) ([]byte, error) {
	srcImg, _, err := image.Decode(bytes.NewReader(src))
	if err != nil {
		return nil, fmt.Errorf("error in resizing: %w", err)
	}

	resImg := imaging.Fill(srcImg, width, height, imaging.Center, imaging.Lanczos)

	var res bytes.Buffer
	err = jpeg.Encode(&res, resImg, nil)
	if err != nil {
		return nil, fmt.Errorf("error in encoding: %w", err)
	}

	return res.Bytes(), nil
}

func NewResizer() *Resizer {
	return &Resizer{}
}
