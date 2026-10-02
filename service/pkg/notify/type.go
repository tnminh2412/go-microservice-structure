package notify

import (
	"encoding/json"
	"fmt"
	"go-pkg/helper/id_helper"
	"go-pkg/ocode"
	"time"
)

var Type = struct {
	PlayerSignInAnotherDevice string
	PlayerSignOut             string
	SetPlayerStatus           string
	Announcement              string
}{
	PlayerSignInAnotherDevice: "player_signin_another_device",
	PlayerSignOut:             "player_sign_out",
	SetPlayerStatus:           "player_set_status",
	Announcement:              "announcement",
}

type NotifyData struct {
	ID   string                 `json:"id"`
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
	Time string                 `json:"time"`
}

func SignoutAnotherDeviceData() string {
	data := map[string]interface{}{}

	notifyData := &NotifyData{
		ID:   ocode.String(id_helper.New()),
		Type: Type.PlayerSignInAnotherDevice,
		Data: data,
		Time: fmt.Sprint(time.Now().UnixMilli()),
	}

	jsonData, err := json.Marshal(notifyData)
	if err != nil {
		return ""
	}

	return string(jsonData)
}

func SignoutData(reason string) string {
	data := map[string]interface{}{
		"reason": reason,
	}

	notifyData := &NotifyData{
		ID:   ocode.String(id_helper.New()),
		Type: Type.PlayerSignOut,
		Data: data,
		Time: fmt.Sprint(time.Now().UnixMilli()),
	}

	jsonData, err := json.Marshal(notifyData)
	if err != nil {
		return ""
	}

	return string(jsonData)
}

func SetPlayerStatusData(status string) string {
	data := map[string]interface{}{
		"status": status,
	}

	notifyData := &NotifyData{
		ID:   ocode.String(id_helper.New()),
		Type: Type.SetPlayerStatus,
		Data: data,
		Time: fmt.Sprint(time.Now().UnixMilli()),
	}

	jsonData, err := json.Marshal(notifyData)
	if err != nil {
		return ""
	}

	return string(jsonData)
}
