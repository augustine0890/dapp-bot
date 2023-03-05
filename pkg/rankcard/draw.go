package rankcard

import (
	"fmt"
	"image"
	"log"
	"net/http"

	"github.com/disintegration/imaging"
	// "github.com/fogleman/gg"
	"github.com/lucasb-eyer/go-colorful"
)

// SetAvatar sets the user's avatar as the source image.
// The given source must be a URL to the user's avatar on Discord.
func (rc *RankCard) SetAvatar(source string) {
	response, err := http.Get(source)
	if err != nil {
		return
	}
	defer response.Body.Close()

	// Read the response body and decode it as an image.
	img, _, err := image.Decode(response.Body)
	if err != nil {
		return
	}

	// Resize the image to the desired dimensions.
	img = imaging.Resize(img, int(rc.Avatar.Width), int(rc.Avatar.Height), imaging.Lanczos)

	// Set the resized image as the source.
	rc.Avatar.Source = img
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

/**
 * Set card overlay
 * @param {string} color Overlay color
 * @param {number} [level=0.5] Opacity level
 * @param {boolean} [display=true] IF it should display overlay
 */
func (rc *RankCard) SetOverlay(color string, level float64, display bool) {
	if color == "" {
		return
	}
	rc.Overlay.Color = color
	rc.Overlay.Display = display
	if level == 0 {
		rc.Overlay.Level = 0.5
	} else {
		rc.Overlay.Level = level
	}
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

// Calculates progress
func (rc *RankCard) calculateProgress() int {
	cx := rc.CurrentXP.Data
	rx := rc.RequiredXP.Data

	if rx <= 0 {
		return 1
	}
	if cx > rx {
		return int(rc.ProgressBar.Width)
	}

	width := (cx * 615) / rx
	if float64(width) > rc.ProgressBar.Width {
		width = int(rc.ProgressBar.Width)
	}
	return width
}
