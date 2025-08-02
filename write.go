package main

import (
	"fmt"
	"image"
	"image/color"
)

type Color struct {
	R, G, B uint8
}

func write(img image.Image) (string, image.Image, error) {

	file, _ := colorPickerImage.Open("colorpicker.png")
	palette, _, err := image.Decode(file)
	if err != nil {
		fmt.Printf("Error loading color palette: %v\n", err)
		return "", nil, err
	}

	// read the palette's pixels and create a map of colors
	colorMap := make(map[Color]Vector2)
	for y := 0; y < palette.Bounds().Dy(); y++ {
		for x := 0; x < palette.Bounds().Dx(); x++ {
			c := palette.At(x, y)
			r, g, b, _ := c.RGBA()
			colorMap[Color{R: uint8(r), G: uint8(g), B: uint8(b)}] = Vector2{X: float64(x) / float64(palette.Bounds().Dx()), Y: float64(y) / float64(palette.Bounds().Dy())}
		}
	}

	println("Found", len(colorMap), "colors in the palette")

	// find the closest color for each pixel in the image
	remappedImage := image.NewRGBA(img.Bounds())
	uvImage := make([][]Vector2, img.Bounds().Dy())
	for i := range uvImage {
		uvImage[i] = make([]Vector2, img.Bounds().Dx())
	}
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			c := img.At(x, y)
			r, g, b, _ := c.RGBA()
			closestColor := findClosestColor(Color{R: uint8(r), G: uint8(g), B: uint8(b)}, colorMap)
			uvImage[y][x] = colorMap[closestColor]
			remappedImage.Set(x, y, color.RGBA{R: closestColor.R, G: closestColor.G, B: closestColor.B, A: 255})
		}
	}

	/*
		// save the remapped image
		outFile, err := os.Create("remapped.png")
		if err != nil {
			fmt.Printf("Error creating output file: %v\n", err)
			return
		}
		defer outFile.Close()
		err = png.Encode(outFile, remappedImage)
		if err != nil {
			fmt.Printf("Error encoding image: %v\n", err)
			return
		}
	*/

	// create a image string from the UV coordinates
	var uvString string
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := img.Bounds().Dy() - 1; y >= 0; y-- {
			uv := uvImage[y][x]
			// flip the Y coordinate
			uv.Y = 1 - uv.Y
			uvString += fmt.Sprintf("%.2f:%.2f,", uv.X, uv.Y)
		}
	}

	return uvString, remappedImage, nil

}

func findClosestColor(target Color, colorMap map[Color]Vector2) Color {
	var closestColor Color
	minDistance := float64(1<<63 - 1) // Start with a large number

	var colorDistance = func(c1, c2 Color) float64 {
		return float64((int(c1.R)-int(c2.R))*(int(c1.R)-int(c2.R)) +
			(int(c1.G)-int(c2.G))*(int(c1.G)-int(c2.G)) +
			(int(c1.B)-int(c2.B))*(int(c1.B)-int(c2.B)))
	}

	for color := range colorMap {
		distance := colorDistance(target, color)
		if distance < minDistance {
			minDistance = distance
			closestColor = color
		}
	}

	return closestColor
}
