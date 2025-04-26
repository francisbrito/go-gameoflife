package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"log"
)

const (
	screenWidth  = 640
	screenHeight = 480
	cellSize     = 10
	columns      = screenWidth / cellSize
	rows         = screenHeight / cellSize
)

type GameState int

func (s GameState) String() string {
	if s == Paused {
		return "Paused"
	}
	if s == Running {
		return "Running"
	}
	return "Unknown"
}

const (
	Paused  GameState = 0
	Running           = 1
)

type Game struct {
	gridColor          color.Color
	cellColor          color.Color
	grid               [columns][rows]bool
	ticks              int
	generation         int
	ticksPerGeneration int
	state              GameState
}

func (g *Game) Update() error {
	if g.state == Running {
		g.ticks++
	}
	if g.ticks == g.ticksPerGeneration {
		g.ticks = 0
		g.cycle()
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if g.state == Running {
			g.state = Paused
		}
		x, y := ebiten.CursorPosition()
		cellX, cellY := x/cellSize, y/cellSize
		g.grid[cellX][cellY] = !g.grid[cellX][cellY]
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.toggleState()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.reset()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawGrid(screen)
	g.drawDebugInfo(screen)
}

func (g *Game) drawDebugInfo(screen *ebiten.Image) {
	fps := ebiten.ActualFPS()
	tps := ebiten.ActualTPS()
	maxTps := ebiten.TPS()
	msg := fmt.Sprintf("FPS: %.2f\nTPS: %.2f (%d)\nTPG: %d\nGeneration: %d\nGame State: %s",
		fps, tps, maxTps, g.ticksPerGeneration, g.generation, g.state)
	ebitenutil.DebugPrintAt(screen, msg, cellSize, cellSize)
}

func (g *Game) drawGrid(screen *ebiten.Image) {
	for i := 0; i < columns; i++ {
		x := float32(cellSize * i)
		vector.StrokeLine(screen, x, 0, x, screenHeight, 1.0, g.gridColor, true)
	}
	for j := 0; j < rows; j++ {
		y := float32(cellSize * j)
		vector.StrokeLine(screen, 0, y, screenWidth, y, 1.0, g.gridColor, true)
	}
	for i := 0; i < columns; i++ {
		for j := 0; j < rows; j++ {
			isAlive := g.grid[i][j]
			x, y := float32(i*cellSize), float32(j*cellSize)
			if isAlive {
				vector.DrawFilledRect(screen, x, y, cellSize, cellSize, g.cellColor, true)
			}
		}
	}
}

func (g *Game) cycle() {
	// Create a new grid for the next generation
	var newGrid [columns][rows]bool

	for i := 0; i < columns; i++ {
		for j := 0; j < rows; j++ {
			count := g.countLiveNeighbors(i, j)

			if g.grid[i][j] {
				// Live cell stays alive if it has 2 or 3 live neighbors
				newGrid[i][j] = count == 2 || count == 3
			} else {
				// Dead cell becomes alive if it has exactly 3 live neighbors
				newGrid[i][j] = count == 3
			}
		}
	}

	// Update the grid with the new generation
	g.grid = newGrid
	g.generation++
}

func (g *Game) countLiveNeighbors(x, y int) int {
	count := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			// Skip the cell itself
			if i == 0 && j == 0 {
				continue
			}

			// Calculate neighbor coordinates
			nx, ny := x+i, y+j

			// Check if neighbor is within bounds
			if nx >= 0 && nx < columns && ny >= 0 && ny < rows {
				// Count live neighbors
				if g.grid[nx][ny] {
					count++
				}
			}
		}
	}
	return count
}

func (g *Game) Layout(w, h int) (int, int) {
	return w, h
}

func (g *Game) toggleState() {
	if g.state == Paused {
		g.state = Running
	} else {
		g.state = Paused
	}
	g.ticks = 0
}

func (g *Game) reset() {
	for i := 0; i < columns; i++ {
		for j := 0; j < rows; j++ {
			g.grid[i][j] = false
			g.generation = 0
		}
	}
}

func main() {
	ebiten.SetWindowTitle("Game of Life")
	ebiten.SetWindowSize(screenWidth, screenHeight)
	g := &Game{
		gridColor:          color.RGBA{R: 63, G: 63, B: 63, A: 255},
		cellColor:          color.White,
		ticksPerGeneration: ebiten.TPS() * 4,
	}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
