package midd

import (
	"fmt"

	"github.com/teoit/gosctx"
	"github.com/teoit/gosctx/component/discord"
)

func RecoveryMiddCustom(logger gosctx.Logger, discordSvc discord.SendMessageDiscordSVC) {
	defer func() {
		if err := recover(); err != nil {
			// Recovered from panic
			errMsg := fmt.Sprintf("Recovered from panic: %v", err)
			logger.Error(errMsg)
			discordSvc.SendMessageDev(errMsg)
		}
	}()
}
