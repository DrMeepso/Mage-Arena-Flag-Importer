# Flag Importer

A Windows utility that imports images as flags into MageArena.

## How to Use

1. Create or find a image that's 100x66 pixels
2. Either drag it onto the exe file or run the command
3. In the flag editor, click "Load From Disk"
4. Done!


## Usage

### Import a Flag Image

Run from command line:
```
flagimporter.exe path/to/your/flag.png
```

### Export Current Flag

Save your current in-game flag to a PNG file:
```
flagimporter.exe --save
```

This will create an `output.png` file with your current flag.

## Drag and Drop

You can also drag and drop a PNG file directly onto the `flagimporter.exe` executable to import it as your flag.

## Requirements

- Windows computer
- MageArena game installed

## Image Requirements

- **Format**: PNG & JPEG supported
- **Size**: Must be exactly 100x66 pixels
- **Color**: Any colors work, will be converted to the game's color palette