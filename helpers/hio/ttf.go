package hio

import (
	_ "embed"
	"io"
	"os"

	"github.com/golang/freetype/truetype"
)

//go:embed JBMono_bold.ttf
var FontJetBrainsMonoBold []byte

func ParseTTFFile(name string) (*truetype.Font, error) {
	fps, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	b, err := io.ReadAll(fps)
	if err != nil {
		return nil, err
	}
	return truetype.Parse(b)
}
