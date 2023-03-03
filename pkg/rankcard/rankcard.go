package rankcard

import (
	"fmt"
	"image/color"
	"log"

	"github.com/fogleman/gg"
	"github.com/lucasb-eyer/go-colorful"
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
	Source string
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

// SetAvatar sets the avatar URL.
func (rc *RankCard) SetAvatar(url string) {
	rc.Avatar.Source = url
}

// SetUsername sets the username and username color.
func (rc *RankCard) SetUsername(username, color string) {
	rc.UserName.Name = username
	rc.UserName.Color = color
}

// SetDiscriminator sets the discriminator.
func (rc *RankCard) SetDiscriminator(discriminator interface{}, color string) {
	discStr, ok := discriminator.(string)
	if !ok {
		discInt, ok := discriminator.(int)
		if !ok || len(fmt.Sprint(discInt)) != 4 {
			rc.Discriminator.Discrim = ""
		} else {
			discStr = fmt.Sprint(discInt)
		}
	}
	rc.Discriminator.Discrim = discStr
	rc.Discriminator.Color = color
}

// SetProgressBar sets the progress bar style
// SetProgressBar sets the progress bar style.
func (rc *RankCard) SetProgressBar(color interface{}, fillType string, rounded bool) error {
	switch fillType {
	case "COLOR":
		c, err := parseColor(color)
		if err != nil {
			return err
		}
		rc.ProgressBar.Bar.Type = "color"
		rc.ProgressBar.Bar.Color = c
	case "GRADIENT":
		colors, err := parseGradientColors(color)
		if err != nil {
			return err
		}
		rc.ProgressBar.Bar.Type = "gradient"
		rc.ProgressBar.Bar.Grad = colors
	default:
		return fmt.Errorf("unsupported progressbar type %q", fillType)
	}
	rc.ProgressBar.Rounded = rounded
	return nil
}

// parseColor parses a color from a string.
func parseColor(c interface{}) (color.Color, error) {
	switch c := c.(type) {
	case string:
		if c[0] == '#' {
			return gg.ParseHexColor(c)
		}
		return gg.ParseColor(c)
	case color.Color:
		return c, nil
	default:
		return nil, fmt.Errorf("invalid color type %T", c)
	}
}

// parseGradientColors parses gradient colors from an interface{}.
func parseGradientColors(colors interface{}) ([]color.Color, error) {
	switch colors := colors.(type) {
	case []string:
		result := make([]color.Color, len(colors))
		for i, c := range colors {
			if c[0] == '#' {
				col, err := colorful.ParseHexColor(c)
				if err != nil {
					return nil, err
				}
				result[i] = col
			} else {
				col, err := colorful.ParseColor(c)
				if err != nil {
					return nil, err
				}
				result[i] = col
			}
		}
		return result, nil
	case []color.Color:
		return colors, nil
	default:
		return nil, fmt.Errorf("invalid gradient color type %T", colors)
	}
}

// SetProgressBarTrack sets the progress bar track color
func (rc *RankCard) SetProgressBarTrack(color string) {
	c, err := colorful.Hex(color)
	if err != nil {
		log.Fatalf("Error parsing color: %v", err)
	}
	rc.ProgressBar.Track.Color = c
}

// SetStatus sets the user's status.
func (rc *RankCard) SetStatus(status string, circle bool, width interface{}) {
	switch status {
	case "online":
		rc.Status.Type = status
		rc.Status.Color = "#43B581"
		break
	case "idle":
		rc.Status.Type = status
		rc.Status.Color = "#FAA61A"
		break
	case "dnd":
		rc.Status.Type = status
		rc.Status.Color = "#F04747"
		break
	case "offline":
		rc.Status.Type = status
		rc.Status.Color = "#747F8E"
		break
	case "streaming":
		rc.Status.Type = status
		rc.Status.Color = "#593595"
		break
	default:
		log.Fatalf("Invalid status '%s'", status)
	}

	if width != false {
		switch width.(type) {
		case int:
			rc.Status.Width = float64(width.(int))
		case float64:
			rc.Status.Width = width.(float64)
		}
	} else {
		rc.Status.Width = 5
	}

	rc.Status.Circle = circle
}

// SetRank sets the rank of the user.
func (rc *RankCard) SetRank(data int, displayText string, isDisplay bool) {
	rc.Rank.Data = data
	rc.Rank.DisplayText = displayText
	rc.Rank.Display = isDisplay
}

// SetLevel sets the level of the user.
func (rc *RankCard) SetLevel(data int, displayText string, isDisplay bool) {
	rc.Level.Data = data
	rc.Level.DisplayText = displayText
	rc.Level.Display = isDisplay
}

// SetCurrentXP sets the user's current experience (XP).
func (rc *RankCard) SetCurrentXP(xp int) {
	rc.CurrentXP.Data = xp
}

// SetRequiredXP sets the required XP to get to the next level.
func (rc *RankCard) SetRequiredXP(xp int) {
	rc.RequiredXP.Data = xp
}

// SetBackground sets the background type.
func (rc *RankCard) SetBackground(typ string, data interface{}) {
	switch typ {
	case "color":
		rc.Background.Type = typ
		rc.Background.Color = data.(string)
	case "image":
		rc.Background.Type = typ
		rc.Background.ImageURL = data.(string)
	default:
		log.Fatalf("Unsupported background type %q", typ)
	}
}
