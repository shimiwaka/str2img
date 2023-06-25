package main

import (
    "bytes"
	"fmt"

	"github.com/shimiwaka/str2img"
)

func main() {
	generater := &str2img.Generater{
		ImageHeight: 630,
		ImageWidth: 1200,
		FontSize: 40.0,
		FontFile: "Koruri-Regular.ttf",
		ImageBytes: &bytes.Buffer{},
	}

	err := generater.Generate("テストてすとです\nほげほげほげ太郎\nほげ太郎")
	if err != nil {
		fmt.Printf("%v", err)
	}

	err = generater.OutputImageFile("test.png")
	if err != nil {
		fmt.Printf("%v", err)
	}
}