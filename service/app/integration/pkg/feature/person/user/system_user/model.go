package system_user

import (
	"proto-repo/generated/pkg/nullable"
	"time"
)

type TaskRole struct {
	Role    string
	Tasks   []string
	Locales []*Locale
}

type Locale struct {
	Language    string
	Description string
}
type Info struct {
	Username         string
	FirstName        string
	LastName         string
	DisplayName      string
	Roles            []*TaskRole
	Password         string
	Level            int32
	IsSystem         bool
	Status           string
	LastLoginCountry string
	LastLoginIP      string
	LastLoginTime    time.Time

	IsRequireChangePassword bool
	Data                    map[string]string
	Email                   string
	Chat                    *Chat

	CreatedBy   string
	CreatedTime time.Time
	UpdatedBy   string
	UpdatedTime time.Time
}

type CreateRequest struct {
	FirstName   string
	LastName    string
	Roles       []string
	DisplayName string
	Password    string
	Language    string
	CreatedBy   string
	CreatedTime time.Time
	Email       string
	AccessLines []string
}

type UpdateRequest struct {
	FirstName   *nullable.NullString
	LastName    *nullable.NullString
	DisplayName *nullable.NullString
	Roles       []string
	Data        map[string]string
	UpdatedBy   string
	UpdatedTime time.Time
	Email       *nullable.NullString
}

type ChangePasswordRequest struct {
	OldPassword string
	NewPassword string
	Language    string
	UpdatedBy   string
	UpdatedTime time.Time
}

type SetPasswordRequest struct {
	Password    string
	Language    string
	UpdatedBy   string
	UpdatedTime time.Time
}

type SetStatusRequest struct {
	Status      string
	UpdatedBy   string
	UpdatedTime time.Time
}

type LoginRequest struct {
	Password         string
	LastLoginCountry string
	LastLoginIP      string
	LastLoginTime    time.Time
	UpdatedBy        string
	UpdatedTime      time.Time
	Otp              string
}

type LoginResult struct {
	Token       string
	IsOtpEnable bool
	IsOtpSkip   bool
}

type TokenInfo struct {
	Token       string
	SubUsername string
	Username    string
	ClientID    string
	Roles       []string
	Tasks       []string
	IsSystem    bool
	IsMain      bool
}

type AccessCheckRequest struct {
	AccessUsername string
	TargetUsername string
}

type OtpGenerateParams struct {
	UpdatedBy   string
	UpdatedTime time.Time
}
type OtpGenerateResult struct {
	Url       string
	Activated bool
	Secret    string
}
type OtpAcceptParams struct {
	Otp         string
	Time        time.Time
	UpdatedBy   string
	UpdatedTime time.Time
	Token       string
}
type OtpResetParams struct {
	UpdatedBy   string
	UpdatedTime time.Time
}

type OtpSkipParams struct {
	UpdatedBy   string
	UpdatedTime time.Time
	Value       bool
}

type internalCreateRequest struct {
	ClientID    string
	Username    string
	FirstName   string
	LastName    string
	DisplayName string
	Password    string
	CreatedBy   string
	CreatedTime time.Time
	Language    string
	Roles       []string
	Email       string
	AccessLines []string
}

type internalUpdateRequest struct {
	FirstName   *nullable.NullString
	LastName    *nullable.NullString
	DisplayName *nullable.NullString
	Roles       []string
	NewRoles    []string
	UpdatedBy   string
	UpdatedTime time.Time
	Data        map[string]string
	Email       *nullable.NullString
}

type FreshWorkApp struct {
	ClientID     string
	LoginPageUri string
	Name         string
	Description  string
}

type FreshWork struct {
	PlayerApp *App
	AgentID   string
	Enable    bool
}

type Chat struct {
	FreshWork *FreshWork
}

type App struct {
	Ticket      *FreshWorkApp
	ChatSupport *FreshWorkApp
}
