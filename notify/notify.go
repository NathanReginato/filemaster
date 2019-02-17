package notify

import (
	"github.com/0xAX/notificator"
)

const (
	appName         = "Organizer"
	finishedMessage = "File Organization Complete"
	iconPath        = "/home/user/icon.png"
)

// Finished notify that the process is finished
func Finished() {
	notify := notificator.New(notificator.Options{
		DefaultIcon: "icon/default.png",
		AppName:     "Organizer",
	})

	notify.Push(appName,
		finishedMessage,
		iconPath,
		notificator.UR_CRITICAL)
}
