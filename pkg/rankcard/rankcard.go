package rankcard

import (
	"image/color"
)

type Background struct {
	Type     string // color or image
	ImageURL string
	Color    string
}

type Track struct {
	Color color.Color
}

type Bar struct {
	Type  string // color or gradient
	Color color.Color
	Grad  []color.Color
}

type ProgressBar struct {
	Rounded   bool
	X         float64
	Y         float64
	Height    float64
	Width     float64
	Track     Track
	Bar       Bar
	Direction string // horizontal or vertical
}

type Overlay struct {
	Display      bool
	Level        float64
	Color        string
	ColorOpacity float64
}

type Avatar struct {
	Source interface{} // URL to user's avatar on Discord
	X      float64
	Y      float64
	Height float64
	Width  float64
}

type Status struct {
	Width  float64
	Type   string // online, offline or idle
	Color  string
	Circle bool
}

type Rank struct {
	Display     bool
	Data        int
	TextColor   string
	Color       string
	DisplayText string
}

type Level struct {
	Display     bool
	Data        int
	TextColor   string
	Color       string
	DisplayText string
}

type CurrentXP struct {
	Data  int
	Color string
}

type RequiredXP struct {
	Data  int
	Color string
}

type Discriminator struct {
	Discrim string
	Color   string
}

type UserName struct {
	Name  string
	Color string
}

type RankCard struct {
	Width         float64
	Height        float64
	Background    Background
	ProgressBar   ProgressBar
	Overlay       Overlay
	Avatar        Avatar
	Status        Status
	Rank          Rank
	Level         Level
	CurrentXP     CurrentXP
	RequiredXP    RequiredXP
	Discriminator Discriminator
	UserName      UserName
	RenderEmojis  bool
}

// NewRankCard creates a new RankCard instance
func NewRankCard() *RankCard {
	return &RankCard{
		Width:         500,
		Height:        200,
		Background:    Background{},
		ProgressBar:   ProgressBar{},
		Overlay:       Overlay{},
		Avatar:        Avatar{},
		Status:        Status{},
		Rank:          Rank{},
		Level:         Level{},
		CurrentXP:     CurrentXP{},
		RequiredXP:    RequiredXP{},
		Discriminator: Discriminator{},
		UserName:      UserName{},
		RenderEmojis:  true,
	}
}
