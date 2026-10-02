package notify

import "fmt"

type topic struct{}

var Topic topic

func (t topic) UserRole(clientID string, role string) string {
	return fmt.Sprintf("_%s.user_role.%s", clientID, role)
}
