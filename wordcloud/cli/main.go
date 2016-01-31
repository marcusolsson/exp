package main

import (
	"bufio"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/freetype/truetype"
	"github.com/marcusolsson/wordcloud"
)

var (
	defaultWidth  = 1024
	defaultHeight = 1024
	defaultFont   = "/usr/share/fonts/truetype/roboto/Roboto-Regular.ttf"
	defaultOutput = "out.png"
)

func main() {
	words := make(map[string]int)

	// Read lines from stdin and count the occurrences.
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		w := scanner.Text()
		if _, ok := words[w]; !ok {
			words[w] = 0
		}

		words[w] = words[w] + 1
	}

	// Read the default font.
	f, err := readFont(defaultFont)
	if err != nil {
		log.Fatal(err)
	}

	// Generate the cloud.
	c := wordcloud.Cloud{
		Width:     defaultWidth,
		Height:    defaultHeight,
		Font:      f,
		Generator: wordcloud.NewSpiralGenerator(),
	}

	im := c.Generate(words)

	// Save to file.
	if err := writeToFile(defaultOutput, im); err != nil {
		log.Fatal(err)
	}
}

func writeToFile(str string, im image.Image) error {
	f, err := os.Create(str)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	if err = png.Encode(w, im); err != nil {
		return err
	}

	if err = w.Flush(); err != nil {
		return err
	}

	return nil
}

func readFont(s string) (*truetype.Font, error) {
	b, err := ioutil.ReadFile(s)
	if err != nil {
		return nil, err
	}

	f, err := truetype.Parse(b)
	if err != nil {
		return nil, err
	}

	return f, nil
}
