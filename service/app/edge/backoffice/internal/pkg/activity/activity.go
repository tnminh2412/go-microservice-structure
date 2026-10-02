package activity

import "time"

var ID = struct {
	Login     string
	Logout    string
	SetStatus string
}{
	Login:     "user_login",
	Logout:    "user_logout",
	SetStatus: "user_set_status",
}

var Data = struct {
	SystemUser string
	Status     string
}{
	SystemUser: "system_user",
	Status:     "status",
}

type Activity struct {
	Username    string
	ActivityID  string
	Time        time.Time
	Data        map[string]string
	Country     string
	IP          string
	CreatedBy   string
	CreatedTime time.Time
}
