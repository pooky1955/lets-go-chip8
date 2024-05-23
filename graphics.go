package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

type Graphics interface {
	Render()
	Clear()
	Set(x, y int, value bool) error
	Toggle(x, y int) (bool, error)
	Read(x, y int) (bool, error)
	Init(width, height int) error
	CheckDimensions(x, y int) (int, int, error)
}

type DimensionError struct {
	Details string
}

func (de *DimensionError) Error() string {
	return fmt.Sprintf("Invalid dimensions received: %s", de.Details)

}

type CLIGraphics struct {
	Grid   []([]bool)
	Width  int
	Height int
}

func (cg *CLIGraphics) Init(width, height int) error {
	if width <= 0 || height <= 0 {
		return &DimensionError{Details: fmt.Sprintf("expected width and height to be strictly positive, but received %d x %d dimensions", width, height)}
	}
	grid := make([][]bool, height)
	for i := 0; i < height; i++ {
		grid[i] = make([]bool, width)
	}
	cg.Grid = grid
	cg.Width = width
	cg.Height = height
	return nil
}

func clearScreen() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("clear") //Linux example, its tested
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls") //Windows example, its tested
	default:
		fmt.Println("CLS for ", runtime.GOOS, " not implemented")
		return
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func (cg *CLIGraphics) Show() {
	var buffer bytes.Buffer
	for _, row := range cg.Grid {
		for _, value := range row {
			character := "．"
			if value {
				character = "＃"
			}
			buffer.WriteString(character)
		}
		buffer.WriteString("\n")
	}
	fmt.Print(buffer.String())
}

func (cg *CLIGraphics) Clear() {
	for i, row := range cg.Grid {
		for j := range row {
			cg.Grid[i][j] = false
		}
	}
}

func (cg *CLIGraphics) Render() {
	clearScreen()
	cg.Show()
}

// CheckDimensions verifies the dimension of the accessed row and column
func (cg *CLIGraphics) CheckDimensions(x, y int) (int, int, error) {
	newX := x
	newY := y
	if x >= cg.Width {
		newX = x - cg.Width
	} else if x < 0 {
		newX = x + cg.Width
	}

	if y >= cg.Height {
		newY = y - cg.Height
	} else if y < 0 {
		newY = y + cg.Height
	}
	INVALID := -1
	if newX < 0 || newX >= cg.Width {
		return INVALID, INVALID, &DimensionError{Details: fmt.Sprintf("expected x to be between %d and %d, but received %d", -cg.Width, 2*cg.Width, x)}
	}
	if newY < 0 || newY >= cg.Height {
		return INVALID, INVALID, &DimensionError{Details: fmt.Sprintf("expected y to be between %d and %d, but received %d", -cg.Height, 2*cg.Height, y)}

	}
	return newX, newY, nil
}

func (cg *CLIGraphics) CheckDimensionsStrict(x, y int) error {
	if x < 0 || x >= cg.Width {
		return &DimensionError{Details: fmt.Sprintf("expected x to be between %d and %d, but received %d", 0, cg.Width, x)}
	}
	if y < 0 || y >= cg.Height {
		return &DimensionError{Details: fmt.Sprintf("expected y to be between %d and %d, but received %d", 0, cg.Height, y)}

	}
	return nil

}

// Set sets the pixel value at a specific x-y coordinate
func (cg *CLIGraphics) Set(x, y int, value bool) error {
	err := cg.CheckDimensionsStrict(x, y)
	if err != nil {
		return err
	}

	cg.Grid[y][x] = value
	return nil
}

// Toggle toggles the value at a specific x-y coordinate. It returns 1 if a pixel was erased
func (cg *CLIGraphics) Toggle(x, y int) (bool, error) {
	newX, newY, err := cg.CheckDimensions(x, y)
	if err != nil {
		return false, err
	}
	cg.Grid[newY][newX] = !cg.Grid[newY][newX]
	return !cg.Grid[newY][newX], nil
}

// Read reads the pixel value at the specified x-y coordinate
func (cg *CLIGraphics) Read(x, y int) (bool, error) {
	err := cg.CheckDimensionsStrict(x, y)
	if err != nil {
		return false, err
	}

	return cg.Grid[y][x], nil
}
