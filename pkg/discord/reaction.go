package discord

import (
	"fmt"

	"github.com/augustine0890/dapp-bot/pkg/logging"
	"github.com/bwmarrin/discordgo"
)

// HandleRemoveReaction removes specific reactions from a message in response to a reaction add event.
func HandleRemoveReaction(s *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
	channelID := reaction.ChannelID
	messageID := reaction.MessageID

	// Define an array of reactions to be removed
	removeReactions := []string{
		"🖕🏻", "🖕", "🖕🏽",
	}

	// Check if bot has permission to remove reactions
	perms, err := s.UserChannelPermissions(s.State.User.ID, channelID)
	if err != nil {
		logging.Error("Failed to get bot permissions: %v", err)
		return
	}
	if perms&discordgo.PermissionManageMessages == 0 {
		msg := fmt.Sprintf("Bot does not have permission to manage messages in channel %s", channelID)
		logging.Warn(msg)
		return
	}

	// Check if the reaction emoji is in the removeReactions array
	for _, emoji := range removeReactions {
		if emoji == reaction.Emoji.Name {
			err := s.MessageReactionsRemoveEmoji(channelID, messageID, emoji)
			if err != nil {
				logging.Error("Failed to remove reaction:", err)
			}
			break
		}
	}
}
