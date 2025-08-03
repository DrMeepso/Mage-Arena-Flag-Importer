package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"os"

	"embed"

	"golang.org/x/sys/windows/registry"
)

//go:embed colorpicker.png
var colorPickerImage embed.FS

const FLAG_PATH = "Software\\jrsjams\\MageArena"
const FLAG_KEY = "flagGrid_h3042110417"

type Vector2 struct {
	X float64
	Y float64
}

var FLAGSIZE = Vector2{X: float64(100), Y: float64(66)} // Assuming the flag size is 100x66 pixels

func main() {

	// check if a png file was dropped on the executable
	if len(os.Args) < 2 {
		fmt.Println("Drag and drop a PNG image onto the exe to set your ingame flag.")
		println()
		println("Or run the command with a PNG file as an argument:")
		fmt.Println("Usage: flagimporter <image.png>")
		fmt.Println("Or use --save to save the current flag image")

		waitForUserInput()

		return
	}
	imagePath := os.Args[1]

	key, err := registry.OpenKey(registry.CURRENT_USER, FLAG_PATH, registry.QUERY_VALUE|registry.ENUMERATE_SUB_KEYS|registry.SET_VALUE)
	if err != nil {
		fmt.Printf("Error opening registry key: %v\n", err)
		return
	}

	// if the 2nd argument is --save, save the image from the registry
	if imagePath == "--save" {

		uvString, _, err := key.GetBinaryValue(FLAG_KEY)
		if err != nil {
			fmt.Printf("Error getting binary value from registry: %v\n", err)
			return
		}

		img, err := readImage(string(uvString))
		if err != nil {
			fmt.Printf("Error reading image: %v\n", err)
			return
		}

		// save the image to a file
		outFile, err := os.Create("output.png")
		if err != nil {
			fmt.Printf("Error creating output file: %v\n", err)
			return
		}

		defer outFile.Close()
		err = png.Encode(outFile, img)
		if err != nil {
			fmt.Printf("Error encoding image: %v\n", err)
			return
		}

		fmt.Println("Image saved to output.png")
		return

	} else {

		if _, err := os.Stat(imagePath); os.IsNotExist(err) {
			fmt.Printf("File does not exist: %s\n", imagePath)
			return
		}

		inImg, err := loadImage(imagePath)
		if err != nil {
			fmt.Printf("Error loading input image: %v\n", err)
			waitForUserInput()
			return
		}

		// check if the image is the same size as the flag
		if float64(inImg.Bounds().Size().X) != FLAGSIZE.X || float64(inImg.Bounds().Size().Y) != FLAGSIZE.Y {
			fmt.Printf("Image is not the same size as the flag (%v vs %v)\n", inImg.Bounds().Size(), FLAGSIZE)
			fmt.Println("Please use an image with the same size as the flag (100x66 pixels).")
			waitForUserInput()
			return
		}

		println("Loaded Image from", imagePath)

		uvString, _, err := write(inImg)
		if err != nil {
			fmt.Printf("Error writing image: %v\n", err)
			return
		}

		err = key.SetBinaryValue(FLAG_KEY, []byte(uvString))
		if err != nil {
			fmt.Printf("Error writing to registry: %v\n", err)
			return
		}

		println("Your ingame flag has been updated! ")
		key.Close()
	}

	// wait for user input before closing
	waitForUserInput()
}

func waitForUserInput() {
	fmt.Println("Press Enter to exit...")
	var input string
	fmt.Scanln(&input)
}

func loadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open image file: %w", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image file: %w", err)
	}
	return img, nil
}
