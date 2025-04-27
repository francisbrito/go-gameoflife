package game

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"log"
)

const (
	ScreenWidth  = 640
	ScreenHeight = 480
	MinCellSize  = 5
	maxColumns   = ScreenWidth / MinCellSize
	maxRows      = ScreenHeight / MinCellSize
)

type State int

func (s State) String() string {
	if s == Paused {
		return "Paused"
	}
	if s == Running {
		return "Running"
	}
	return "Unknown"
}

const (
	Paused State = iota
	Running
)

type Game struct {
	grid               [maxColumns][maxRows]bool
	cellSize           int
	columns            int
	rows               int
	ticks              int
	generation         int
	ticksPerGeneration int
	state              State
	selectedThemeID    ThemeID
	darkTheme          *Theme
	lightTheme         *Theme
}

type Options struct {
	CellSize int
}

func NewFromOptions(options Options) *Game {
	if options.CellSize < MinCellSize {
		log.Fatalf("cell size must be greater than or equal to %d, got %d", MinCellSize, options.CellSize)
	}
	columns := ScreenWidth / options.CellSize
	rows := ScreenHeight / options.CellSize
	darkTheme, lightTheme := NewDarkTheme(), NewLightTheme()
	return &Game{
		cellSize:           options.CellSize,
		columns:            columns,
		rows:               rows,
		state:              Paused,
		ticksPerGeneration: ebiten.TPS() / 8,
		darkTheme:          darkTheme,
		lightTheme:         lightTheme,
		selectedThemeID:    darkTheme.ID,
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
		fps, tps, maxTps, g.ticksPerGeneration, g.generation, g.state, g.theme())
	ebitenutil.DebugPrintAt(screen, msg, 16, 16)
}

func (g *Game) drawGrid(screen *ebiten.Image) {
	theme := g.theme()
	for i := 0; i < g.columns; i++ {
		x := float32(g.cellSize * i)
		vector.StrokeLine(screen, x, 0, x, ScreenHeight, 1.0, theme.GridColor, true)
	}
	for j := 0; j < g.rows; j++ {
		y := float32(g.cellSize * j)
		vector.StrokeLine(screen, 0, y, ScreenWidth, y, 1.0, theme.GridColor, true)
	}
	for i := 0; i < g.columns; i++ {
		for j := 0; j < g.rows; j++ {
			isAlive := g.grid[i][j]
			x, y := float32(i*g.cellSize), float32(j*g.cellSize)
			size := float32(g.cellSize)
			if isAlive {
				vector.DrawFilledRect(screen, x, y, size, size, theme.CellColor, true)
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
	if g.selectedThemeID == Dark {
		g.selectedThemeID = Light
	} else {
		g.selectedThemeID = Dark
	}
}

func (g *Game) drawBackground(screen *ebiten.Image) {
	screen.Fill(g.theme().BackgroundColor)
}

func (g *Game) theme() *Theme {
	if g.selectedThemeID == Light {
		return g.lightTheme
	} else {
		return g.darkTheme
	}
}
