package bot

import "github.com/binarycode/trewoga/pkg/model"

func start(user model.User) {
	send(user, "Hi! I can alert you about service failures.")
	help(user)
}
