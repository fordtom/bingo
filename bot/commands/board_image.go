package commands

import (
	"bytes"
	"fmt"
	"image/color"
	"image/png"
	"strings"

	"github.com/fogleman/gg"
	"github.com/fordtom/bingo/db"
)

const (
	cellSize    = 150  // pixels per cell
	padding     = 10   // padding around board
	fontSize    = 14.0 // base font size
	lineSpacing = 1.3  // line height multiplier
	maxLines    = 3    // max lines of text per cell
	maxChars    = 47   // truncate with "..." above this
)

var (
	colorOpen      = color.RGBA{180, 180, 180, 255} // grey for open cells
	colorCompleted = color.RGBA{76, 175, 80, 255}   // green for completed cells
	colorBorder    = color.RGBA{60, 60, 60, 255}    // dark grey border
	colorText      = color.RGBA{33, 33, 33, 255}    // dark text
)

// GenerateBoardImage creates a PNG image of the bingo board in memory
func GenerateBoardImage(grid [][]db.BoardSquareWithEvent, gridSize int) ([]byte, error) {
	// Calculate canvas size
	width := gridSize*cellSize + 2*padding
	height := gridSize*cellSize + 2*padding

	// Create drawing context
	dc := gg.NewContext(width, height)

	// Set background
	dc.SetColor(color.RGBA{245, 245, 245, 255})
	dc.Clear()

	// Draw each cell
	for row := 0; row < gridSize; row++ {
		for col := 0; col < gridSize; col++ {
			sq := grid[row][col]
			drawCell(dc, row, col, sq)
		}
	}

	// Encode to PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, dc.Image()); err != nil {
		return nil, fmt.Errorf("failed to encode PNG: %w", err)
	}

	return buf.Bytes(), nil
}

// drawCell renders a single cell with background color and text
func drawCell(dc *gg.Context, row, col int, sq db.BoardSquareWithEvent) {
	x := float64(col*cellSize + padding)
	y := float64(row*cellSize + padding)

	// Draw background
	if sq.EventStatus == "CLOSED" {
		dc.SetColor(colorCompleted)
	} else {
		dc.SetColor(colorOpen)
	}
	dc.DrawRectangle(x, y, cellSize, cellSize)
	dc.Fill()

	// Draw border
	dc.SetColor(colorBorder)
	dc.SetLineWidth(2)
	dc.DrawRectangle(x, y, cellSize, cellSize)
	dc.Stroke()

	// Draw text
	dc.SetColor(colorText)
	if err := dc.LoadFontFace("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", fontSize); err != nil {
		// Fallback if font not found - try common alternatives
		dc.LoadFontFace("/usr/share/fonts/truetype/liberation/LiberationSans-Regular.ttf", fontSize)
	}

	// Wrap text into lines
	lines := wrapText(dc, sq.EventDescription, cellSize-20) // 20px margin

	// Center text vertically
	textHeight := float64(len(lines)) * fontSize * lineSpacing
	startY := y + (cellSize-textHeight)/2 + fontSize

	// Draw each line centered
	for i, line := range lines {
		lineY := startY + float64(i)*fontSize*lineSpacing
		lineWidth, _ := dc.MeasureString(line)
		lineX := x + (cellSize-lineWidth)/2
		dc.DrawString(line, lineX, lineY)
	}
}

// wrapText breaks text into multiple lines, avoiding word splits and truncating with "..."
func wrapText(dc *gg.Context, text string, maxWidth float64) []string {
	// Truncate if too long
	if len(text) > maxChars {
		text = text[:maxChars] + "…"
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}

	var lines []string
	var currentLine strings.Builder

	for _, word := range words {
		// Try adding word to current line
		testLine := currentLine.String()
		if testLine != "" {
			testLine += " " + word
		} else {
			testLine = word
		}

		width, _ := dc.MeasureString(testLine)

		if width <= maxWidth {
			// Word fits, add it
			if currentLine.Len() > 0 {
				currentLine.WriteString(" ")
			}
			currentLine.WriteString(word)
		} else {
			// Word doesn't fit
			if currentLine.Len() == 0 {
				// Single word is too long, break it with hyphen
				currentLine.WriteString(word)
			} else {
				// Start new line
				lines = append(lines, currentLine.String())
				currentLine.Reset()
				currentLine.WriteString(word)
			}
		}

		// Stop if we hit max lines
		if len(lines) >= maxLines {
			break
		}
	}

	// Add remaining text
	if currentLine.Len() > 0 && len(lines) < maxLines {
		lines = append(lines, currentLine.String())
	}

	// If we hit max lines but have more words, truncate last line with "..."
	if len(words) > 0 && len(lines) == maxLines {
		// Check if we've processed all words
		processedWords := 0
		for _, line := range lines {
			processedWords += len(strings.Fields(line))
		}
		if processedWords < len(words) {
			lastLine := lines[len(lines)-1]
			if len(lastLine) > 3 {
				lines[len(lines)-1] = lastLine[:len(lastLine)-3] + "…"
			}
		}
	}

	return lines
}
