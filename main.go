package main

import (
	"gameoflife/game"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func main() {
	ebiten.SetWindowTitle("Game of Life")
	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetScreenClearedEveryFrame(true)
	g := game.NewFromOptions(game.Options{
		CellSize: game.MinCellSize,
	})
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
