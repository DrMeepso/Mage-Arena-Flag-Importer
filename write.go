package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
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
			x := float64(x) / float64(palette.Bounds().Dx())
			y := float64(y) / float64(palette.Bounds().Dy())

			// round x,y to 2 decimal places
			x = float64(int(x*100)) / 100
			y = float64(int(y*100)) / 100

			// if the x is 0 add 0.01 if its 1 remove 0.01
			if x == 0 {
				x += 0.01
			} else if x == 1 {
				x -= 0.01
			}
			if y == 0 {
				y += 0.01
			} else if y == 1 {
				y -= 0.01
			}
			colorMap[Color{R: uint8(r), G: uint8(g), B: uint8(b)}] = Vector2{X: x, Y: y}
		}
	}

	println("Quantizing image...")

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

// Helper function to convert RGB to OKLab
func rgbToOKLab(c Color) (L, a, b float64) {
	// Normalize RGB values to [0, 1]
	r := float64(c.R) / 255.0
	g := float64(c.G) / 255.0
	ba := float64(c.B) / 255.0

	// Linearize RGB values
	r = linearize(r)
	g = linearize(g)
	ba = linearize(ba)

	// Convert to LMS space
	l := 0.4122214708*r + 0.5363325363*g + 0.0514459929*ba
	m := 0.2119034982*r + 0.6806995451*g + 0.1073969566*ba
	s := 0.0883024619*r + 0.2817188376*g + 0.6299787005*ba

	// Apply cube root
	l = cubeRoot(l)
	m = cubeRoot(m)
	s = cubeRoot(s)

	// Convert to OKLab
	L = 0.2104542553*l + 0.7936177850*m - 0.0040720468*s
	a = 1.9779984951*l - 2.4285922050*m + 0.4505937099*s
	b = 0.0259040371*l + 0.7827717662*m - 0.8086757660*s

	return
}

// Linearize RGB values
func linearize(value float64) float64 {
	if value <= 0.04045 {
		return value / 12.92
	}
	return math.Pow((value+0.055)/1.055, 2.4)
}

// Cube root function with handling for small values
func cubeRoot(value float64) float64 {
	if value > 0 {
		return math.Pow(value, 1.0/3.0)
	}
	return -math.Pow(-value, 1.0/3.0)
}

// Function to calculate CIEDE2000 color difference
func cieDE2000(L1, a1, b1, L2, a2, b2 float64) float64 {
	// Constants
	kL, kC, kH := 1.0, 1.0, 1.0

	// Calculate CIEDE2000 components
	deltaL := L2 - L1
	LBar := (L1 + L2) / 2.0

	C1 := math.Sqrt(a1*a1 + b1*b1)
	C2 := math.Sqrt(a2*a2 + b2*b2)
	CBar := (C1 + C2) / 2.0

	aPrime1 := a1 + a1/2.0*(1.0-math.Sqrt(math.Pow(CBar, 7.0)/(math.Pow(CBar, 7.0)+math.Pow(25.0, 7.0))))
	aPrime2 := a2 + a2/2.0*(1.0-math.Sqrt(math.Pow(CBar, 7.0)/(math.Pow(CBar, 7.0)+math.Pow(25.0, 7.0))))

	CPrime1 := math.Sqrt(aPrime1*aPrime1 + b1*b1)
	CPrime2 := math.Sqrt(aPrime2*aPrime2 + b2*b2)
	CBarPrime := (CPrime1 + CPrime2) / 2.0

	deltaCPrime := CPrime2 - CPrime1

	hPrime1 := math.Atan2(b1, aPrime1)
	if hPrime1 < 0 {
		hPrime1 += 2.0 * math.Pi
	}
	hPrime2 := math.Atan2(b2, aPrime2)
	if hPrime2 < 0 {
		hPrime2 += 2.0 * math.Pi
	}

	deltaHPrime := hPrime2 - hPrime1
	if math.Abs(deltaHPrime) > math.Pi {
		if hPrime2 <= hPrime1 {
			deltaHPrime += 2.0 * math.Pi
		} else {
			deltaHPrime -= 2.0 * math.Pi
		}
	}
	deltaHPrime = 2.0 * math.Sqrt(CPrime1*CPrime2) * math.Sin(deltaHPrime/2.0)

	HBarPrime := (hPrime1 + hPrime2) / 2.0
	if math.Abs(hPrime1-hPrime2) > math.Pi {
		if hPrime1+hPrime2 < 2.0*math.Pi {
			HBarPrime += math.Pi
		} else {
			HBarPrime -= math.Pi
		}
	}

	T := 1.0 - 0.17*math.Cos(HBarPrime-math.Pi/6.0) + 0.24*math.Cos(2.0*HBarPrime) +
		0.32*math.Cos(3.0*HBarPrime+math.Pi/30.0) - 0.20*math.Cos(4.0*HBarPrime-63.0*math.Pi/180.0)

	deltaTheta := 30.0 * math.Pi / 180.0 * math.Exp(-math.Pow((180.0/math.Pi*HBarPrime-275.0)/25.0, 2.0))
	RC := 2.0 * math.Sqrt(math.Pow(CBarPrime, 7.0)/(math.Pow(CBarPrime, 7.0)+math.Pow(25.0, 7.0)))
	SL := 1.0 + (0.015*math.Pow(LBar-50.0, 2.0))/math.Sqrt(20.0+math.Pow(LBar-50.0, 2.0))
	SC := 1.0 + 0.045*CBarPrime
	SH := 1.0 + 0.015*CBarPrime*T
	RT := -math.Sin(2.0*deltaTheta) * RC

	// Final CIEDE2000 formula
	return math.Sqrt(math.Pow(deltaL/(kL*SL), 2.0) +
		math.Pow(deltaCPrime/(kC*SC), 2.0) +
		math.Pow(deltaHPrime/(kH*SH), 2.0) +
		RT*(deltaCPrime/(kC*SC))*(deltaHPrime/(kH*SH)))
}

func findClosestColor(target Color, colorMap map[Color]Vector2) Color {
	var closestColor Color
	minDistance := float64(1<<63 - 1) // Start with a large number

	var colorDistance = func(c1, c2 Color) float64 {
		L1, a1, b1 := rgbToOKLab(c1)
		L2, a2, b2 := rgbToOKLab(c2)
		return cieDE2000(L1, a1, b1, L2, a2, b2)
	}

	/*
		var colorDistance = func(c1, c2 Color) float64 {
			return float64((int(c1.R)-int(c2.R))*(int(c1.R)-int(c2.R)) +
				(int(c1.G)-int(c2.G))*(int(c1.G)-int(c2.G)) +
				(int(c1.B)-int(c2.B))*(int(c1.B)-int(c2.B)))
		}
	*/

	for color := range colorMap {
		distance := colorDistance(target, color)
		if distance < minDistance {
			minDistance = distance
			closestColor = color
		}
	}

	return closestColor
}
