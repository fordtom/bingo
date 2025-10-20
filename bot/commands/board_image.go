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
	maxLines    = 4    // max lines of text per cell
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
	if sq.EventStatus == string(db.EventStatusClosed) {
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
	truncated := false
	if len(text) > maxChars {
		text = text[:maxChars]
		truncated = true
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}

	var lines []string
	var currentLine strings.Builder
	wordIdx := 0

	for wordIdx < len(words) {
		word := words[wordIdx]

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
			wordIdx++
		} else {
			// Word doesn't fit
			if currentLine.Len() == 0 {
				// Single word is too long, add it anyway (will overflow)
				currentLine.WriteString(word)
				wordIdx++
			} else {
				// Start new line with current word
				lines = append(lines, currentLine.String())
				
				// Check if we've hit max lines
				if len(lines) >= maxLines {
					// Add ellipsis to last line since we have more words
					lastLine := lines[len(lines)-1]
					// Ensure there's room for ellipsis
					for len(lastLine) > 0 {
						testWithEllipsis := lastLine + "…"
						width, _ := dc.MeasureString(testWithEllipsis)
						if width <= maxWidth {
							lines[len(lines)-1] = testWithEllipsis
							break
						}
						// Remove last character and try again
						lastLine = lastLine[:len(lastLine)-1]
					}
					break
				}
				
				currentLine.Reset()
				currentLine.WriteString(word)
				wordIdx++
			}
		}
	}

	// Add remaining text if we haven't hit max lines
	if currentLine.Len() > 0 && len(lines) < maxLines {
		line := currentLine.String()
		// If we had truncated the original text, add ellipsis
		if truncated && wordIdx >= len(words) {
			line += "…"
		}
		lines = append(lines, line)
	}

	return lines
}
