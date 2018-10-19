package bot

import "github.com/binarycode/trewoga/pkg/model"

func help(user model.User) {
	send(user, `You can control me by sending following commands:

/version - get Trewoga server version
/status - get status of all subscribed services
/subscribe SERVICE\_TOKEN - subscribe to a service (SERVICE\_TOKEN can be obtained from Trewoga server administrator)
/unsubscribe SERVICE\_NAME - unsubscribe from a service`)
}
