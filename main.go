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
	minCellSize  = 10
	maxColumns   = screenWidth / minCellSize
	maxRows      = screenHeight / minCellSize
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
	Paused GameState = iota
	Running
)

type GameTheme int

const (
	GameThemeDark GameTheme = iota
	GameThemeLight
)

func (g GameTheme) String() string {
	switch g {
	case GameThemeLight:
		return "Light"
	case GameThemeDark:
		return "Dark"
	}
	return fmt.Sprintf("Unknown (%d)", g)
}

type Game struct {
	backgroundColor    color.Color
	gridColor          color.Color
	cellColor          color.Color
	grid               [maxColumns][maxRows]bool
	cellSize           int
	columns            int
	rows               int
	ticks              int
	generation         int
	ticksPerGeneration int
	state              GameState
	theme              GameTheme
}

type NewGameOptions struct {
	Theme    GameTheme
	CellSize int
}

func NewGameFromOptions(options NewGameOptions) *Game {
	if options.CellSize < 10 {
		log.Fatalf("cell size must be greater than 10, got %d", options.CellSize)
	}
	columns := screenWidth / options.CellSize
	rows := screenHeight / options.CellSize
	var gridColor, cellColor, backgroundColor color.Color
	if options.Theme == GameThemeLight {
		backgroundColor = color.White
		gridColor = color.Gray{Y: 127}
		cellColor = color.Gray{Y: 255}
	} else if options.Theme == GameThemeDark {
		backgroundColor = color.Black
		gridColor = color.Gray{Y: 31}
		cellColor = color.White
	} else {
		log.Fatalf("invalid theme: %s", options.Theme)
	}
	return &Game{
		backgroundColor:    backgroundColor,
		gridColor:          gridColor,
		cellColor:          cellColor,
		cellSize:           options.CellSize,
		columns:            columns,
		rows:               rows,
		state:              Paused,
		ticksPerGeneration: ebiten.TPS() / 8,
		theme:              options.Theme,
	}
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
		cellX, cellY := x/g.cellSize, y/g.cellSize
		g.grid[cellX][cellY] = !g.grid[cellX][cellY]
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.toggleState()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.reset()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyT) {
		g.switchTheme()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawBackground(screen)
	g.drawGrid(screen)
	g.drawDebugInfo(screen)
}

func (g *Game) drawDebugInfo(screen *ebiten.Image) {
	fps := ebiten.ActualFPS()
	tps := ebiten.ActualTPS()
	maxTps := ebiten.TPS()
	msg := fmt.Sprintf(
		"FPS: %.2f\nTPS: %.2f (%d)\nTPG: %d\nGeneration: %d\nGame State: %s\nTheme: %s\nPress R to restart\nPress Space to pause\nPress T to switch themes",
		fps, tps, maxTps, g.ticksPerGeneration, g.generation, g.state, g.theme)
	ebitenutil.DebugPrintAt(screen, msg, 16, 16)
}

func (g *Game) drawGrid(screen *ebiten.Image) {
	for i := 0; i < g.columns; i++ {
		x := float32(g.cellSize * i)
		vector.StrokeLine(screen, x, 0, x, screenHeight, 1.0, g.gridColor, true)
	}
	for j := 0; j < g.rows; j++ {
		y := float32(g.cellSize * j)
		vector.StrokeLine(screen, 0, y, screenWidth, y, 1.0, g.gridColor, true)
	}
	for i := 0; i < g.columns; i++ {
		for j := 0; j < g.rows; j++ {
			isAlive := g.grid[i][j]
			x, y := float32(i*g.cellSize), float32(j*g.cellSize)
			size := float32(g.cellSize)
			if isAlive {
				vector.DrawFilledRect(screen, x, y, size, size, g.cellColor, true)
			}
		}
	}
}

func (g *Game) cycle() {
	// Create a new grid for the next generation
	var newGrid [maxColumns][maxRows]bool

	for i := 0; i < g.columns; i++ {
		for j := 0; j < g.rows; j++ {
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
			if nx >= 0 && nx < g.columns && ny >= 0 && ny < g.rows {
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
	for i := 0; i < g.columns; i++ {
		for j := 0; j < g.rows; j++ {
			g.grid[i][j] = false
			g.generation = 0
		}
	}
}

func (g *Game) switchTheme() {
	if g.theme == GameThemeDark {
		g.theme = GameThemeLight
	} else {
		g.theme = GameThemeDark
	}
}

func (g *Game) drawBackground(screen *ebiten.Image) {
	screen.Fill(g.backgroundColor)
}

func main() {
	ebiten.SetWindowTitle("Game of Life")
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetScreenClearedEveryFrame(true)
	g := NewGameFromOptions(NewGameOptions{
		Theme:    GameThemeDark,
		CellSize: 10,
	})
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
