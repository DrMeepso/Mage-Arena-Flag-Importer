package main

import (
	"fmt"
	"image"
	"strconv"
	"strings"
)

func readImage(imageString string) (image.Image, error) {

	value := imageString

	pixels := make([]Vector2, 0)

	// split values by comma
	parts := strings.Split(string(value), ",")
	for _, part := range parts {
		// split by colon
		subParts := strings.Split(part, ":")
		if len(subParts) == 2 {
			u, v := strings.TrimSpace(subParts[0]), strings.TrimSpace(subParts[1])

			// Remove null bytes and other non-printable characters
			u = strings.Trim(u, "\x00")
			v = strings.Trim(v, "\x00")

			// convert to float
			x, err := strconv.ParseFloat(u, 64)
			if err != nil {
				fmt.Printf("Error converting '%s' to float: %v\n", u, err)
				continue
			}
			y, err := strconv.ParseFloat(v, 64)
			if err != nil {
				fmt.Printf("Error converting '%s' to float: %v\n", v, err)
				continue
			}
			// append to pixels
			pixels = append(pixels, Vector2{X: x, Y: y})
		}
	}

	// read the colorpicker.png image for the pallete
	file, _ := colorPickerImage.Open("colorpicker.png")
	colorPalette, _, err := image.Decode(file)
	if err != nil {
		fmt.Printf("Error loading color palette: %v\n", err)
		return nil, err
	}

	// make a bitmap image to draw the pixels on
	img := image.NewRGBA(image.Rect(0, 0, int(FLAGSIZE.X), int(FLAGSIZE.Y)))
	for i, p := range pixels {

		X := i / 66
		Y := 65 - (i % 66)

		// use the vector as u, v coordinates
		u := (p.X) * float64(colorPalette.Bounds().Dx())
		v := (1 - p.Y) * float64(colorPalette.Bounds().Dy())
		// get the color from the color palette
		color := colorPalette.At(int(u)-1, int(v)-1)
		// set the pixel in the image
		img.Set(X, Y, color)
	}

	return img, nil

}
