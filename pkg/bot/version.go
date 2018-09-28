package bot

import (
	"github.com/binarycode/trewoga/pkg/model"

	v "github.com/binarycode/trewoga/pkg/version"
)

func version(user model.User) {
	send(user, v.Version)
}
