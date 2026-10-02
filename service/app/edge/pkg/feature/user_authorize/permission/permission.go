package permission

var Task = struct {
	All              string
	None             string
	ViewSystemUser   string
	CreateSystemUser string
	UpdateSystemUser string

	SetSystemUserPassword string
	SetSystemUserStatus   string
	ResetSystemUserOtp    string

	ViewRole string

	GetSystemUser string
}{
	All:                   "*",
	None:                  "none",
	ViewSystemUser:        "view_system_user",
	CreateSystemUser:      "create_system_user",
	UpdateSystemUser:      "update_system_user",
	SetSystemUserPassword: "set_system_user_password",
	SetSystemUserStatus:   "set_system_user_status",
	ResetSystemUserOtp:    "reset_system_user_otp",

	ViewRole:      "view_role",
	GetSystemUser: "get_system_user",
}
