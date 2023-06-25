package str2img

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"os"
	"strings"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type Generater struct {
	ImageHeight int
	ImageWidth  int
	FontSize    float64
	FontFile    string
	ImageBytes  *bytes.Buffer
}

func fillRect(img *image.RGBA, col color.Color) {
	rect := img.Rect

	for h := rect.Min.Y; h < rect.Max.Y; h++ {
		for v := rect.Min.X; v < rect.Max.X; v++ {
			img.Set(v, h, col)
		}
	}
}

func textSplitter(dr *font.Drawer, fontSize float64, imageHeight int, imageWidth int, texts []string) []string {
	newTexts := []string{}

	for _, v := range texts {
		textWidth := dr.MeasureString(v).Round()
		if textWidth > imageWidth {
			splitLen := int(float64(imageWidth) / float64(textWidth) * float64(len(v)))
			if splitLen < 1 {
				splitLen = 1
			}
			for i := 0; i < len(v); i += splitLen {
				if i+splitLen < len(v) {
					newTexts = append(newTexts, v[i:(i+splitLen)])
				} else {
					newTexts = append(newTexts, v[i:])
				}
			}
		} else {
			newTexts = append(newTexts, v)
		}
		if ((len(newTexts) + 2) * (int(fontSize) + 10)) > imageHeight {
			newTexts = append(newTexts, "...")
			return newTexts
		}
	}
	return newTexts
}

func (g *Generater) Generate(rawText string) error {
	ftBinary, err := ioutil.ReadFile(g.FontFile)
	if err != nil {
		return err
	}

	ft, err := truetype.Parse(ftBinary)
	if err != nil {
		return err
	}
	fontSize := g.FontSize

	opt := truetype.Options{
		Size:              fontSize,
		DPI:               0,
		Hinting:           0,
		GlyphCacheEntries: 0,
		SubPixelsX:        0,
		SubPixelsY:        0,
	}

	imageWidth := g.ImageWidth
	imageHeight := g.ImageHeight

	texts := strings.Split(rawText, "\n")

	img := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))
	fillRect(img, color.RGBA{255, 255, 255, 255})

	face := truetype.NewFace(ft, &opt)

	dr := &font.Drawer{
		Dst:  img,
		Src:  image.Black,
		Face: face,
		Dot:  fixed.Point26_6{},
	}

	texts = textSplitter(dr, fontSize, imageHeight, imageWidth, texts)
	textTopMargin := (imageHeight - ((len(texts) - 1) * (int(fontSize) + 10))) / 2

	for i, v := range texts {
		dr.Dot.X = (fixed.I(imageWidth) - dr.MeasureString(v)) / 2
		dr.Dot.Y = fixed.I(textTopMargin + i*(int(fontSize)+10))

		dr.DrawString(v)
	}

	g.ImageBytes = &bytes.Buffer{}
	err = png.Encode(g.ImageBytes, img)

	if err != nil {
		return err
	}
	return nil
}

func (g *Generater) OutputImageFile(fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	file.Write(g.ImageBytes.Bytes())
	return nil
}
