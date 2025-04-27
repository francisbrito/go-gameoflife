package game

import "image/color"

type ThemeID int

const (
	Dark ThemeID = iota
	Light
)

type Theme struct {
	ID              ThemeID
	BackgroundColor color.Color
	GridColor       color.Color
	CellColor       color.Color
}

func (t *Theme) String() string {
	switch t.ID {
	case Dark:
		return "Dark"
	case Light:
		return "Light"
	default:
		return "Unknown"
	}
}

func NewDarkTheme() *Theme {
	return &Theme{
		ID:              Dark,
		BackgroundColor: color.Gray{Y: 15},
		GridColor:       color.Gray{Y: 31},
		CellColor:       color.White,
	}
}

func NewLightTheme() *Theme {
	return &Theme{
		ID:              Light,
		BackgroundColor: color.White,
		GridColor:       color.Gray{Y: 127},
		CellColor:       color.Black,
	}
}
