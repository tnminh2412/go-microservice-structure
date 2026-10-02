package notify

import "fmt"

type id struct{}

var ID id

func (t id) Player(clientID string, username string) string {
	return fmt.Sprintf("_%s.player.%s", clientID, username)
}

func (t id) Anonymous(clientID string, token string) string {
	return fmt.Sprintf("_%s.anon.%s", clientID, token)
}

func (t id) User(clientID string, username string) string {
	return fmt.Sprintf("_%s.user.%s", clientID, username)
}

var UserType = struct {
	User   string
	Player string
}{
	User:   "user",
	Player: "player",
}
