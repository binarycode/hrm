package bot

import "github.com/binarycode/trewoga/pkg/model"

func help(user model.User) {
	send(user, `You can control me by sending following commands:

/version - get Trewoga server version
/status - get status of all subscribed services
/subscribe SERVICE_TOKEN - subscribe to a service (SERVICE_TOKEN can be obtained from Trewoga server administrator)
/unsubscribe SERVICE_NAME - unsubscribe from a service`)
}
