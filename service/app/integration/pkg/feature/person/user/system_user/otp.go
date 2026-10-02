package system_user

import (
	crypto_rand "crypto/rand"
	"encoding/base32"
	"go-pkg/helper/otp_helper"
	"proto-repo/generated/cmm/service/cmm/core/user/core_user"
	"strconv"
	"time"

	"github.com/tnminh2412/go-microservice-structure/service/app/integration/pkg/feature/person/user/token_user"
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/pkg/otp"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func validateOtp(data map[string]string, otpValue int, time time.Time) bool {
	sec, ok := data[otp.OtpSecret]
	if !ok {
		return false
	}
	o, err := otp_helper.NewGoogleOTP(sec)
	if err != nil {
		return false
	}
	isValid, err := o.ValidateTOTPWithTime(otpValue, time)
	return err == nil && isValid
}

type generateOtpResult struct {
	ProvisionURI string
	Secret       string
}

func generateOtp(username string, issuer string) (*generateOtpResult, error) {
	// length should not be changed
	// on IOS random secret length may cause invalid barcode
	var b [25]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		return nil, err
	}
	sec := base32.StdEncoding.EncodeToString(b[:])
	o, err := otp_helper.NewGoogleOTP(sec)
	if err != nil {
		return nil, err
	}
	uri := o.GetProvisionURIWithIssuer(username, issuer)
	return &generateOtpResult{
		ProvisionURI: uri,
		Secret:       sec,
	}, nil
}

func (s *SystemUser) GenerateOtp(req *OtpGenerateParams) (*OtpGenerateResult, error) {
	issuer, err := getOtpIssuer(s.ctx, s.clientID)
	if err != nil {
		return nil, err
	}
	user, err := s.Get()
	if err != nil {
		return nil, err
	}
	d := user.Data
	enabled := otp.IsOtpEnable(d)
	if enabled {
		return &OtpGenerateResult{Activated: true}, nil
	}

	o, err := generateOtp(s.username, issuer)
	if err != nil {
		return nil, err
	}
	_, err = userSrv.UpdateUser(s.ctx, &core_user.UpdateUserRequest{
		ClientID:    s.clientID,
		ID:          s.username,
		UpdatedBy:   req.UpdatedBy,
		UpdatedTime: timestamppb.New(req.UpdatedTime),
		Data: map[string]string{
			otp.OtpEnable:   "false",
			otp.OtpSecret:   o.Secret,
			otp.OtpProvider: otp.GoogleOtpProvider,
		},
	})
	if err != nil {
		return nil, err
	}
	return &OtpGenerateResult{
		Url:       o.ProvisionURI,
		Activated: false,
		Secret:    o.Secret,
	}, nil
}

func (s *SystemUser) AcceptOtp(req *OtpAcceptParams) error {
	o, err := strconv.ParseInt(req.Otp, 10, 64)
	if err != nil {
		return Error(InvalidOtp, "", err)
	}
	user, err := s.Get()
	if err != nil {
		return err
	}
	d := user.Data
	enabled := otp.IsOtpEnable(d)
	if enabled {
		return nil
	}

	valid := validateOtp(d, int(o), req.Time)
	if !valid {
		return Error(InvalidOtp, "", nil)
	}

	err = s.Update(&UpdateRequest{
		UpdatedBy:   req.UpdatedBy,
		UpdatedTime: req.UpdatedTime,
		Data: map[string]string{
			otp.OtpEnable: "true",
			otp.OtpSkip:   "false",
		},
	})
	if err != nil {
		return err
	}

	token, err := tokenUser.Get(s.ctx, s.clientID, req.Token)
	if err != nil {
		return err
	}

	tokenTTL, err := getTokenTTLConfig(s.ctx, s.clientID)
	if err != nil {
		return err
	}
	err = tokenUser.Set(s.ctx, &token_user.SetData{
		Token:       token.Token,
		SubUsername: token.SubUsername,
		Username:    token.Username,
		ClientID:    token.ClientID,
		Roles:       token.Roles,
		Tasks:       token.Tasks,
		AccessLines: token.AccessLines,
		ParentL1:    token.ParentL1,
		ParentL2:    token.ParentL2,
		ParentL3:    token.ParentL3,
		ParentL4:    token.ParentL4,
		ParentL5:    token.ParentL5,
		ParentL6:    token.ParentL6,
		ParentL7:    token.ParentL7,
		ParentL8:    token.ParentL8,
		ParentL9:    token.ParentL9,
		ParentL10:   token.ParentL10,
		ParentL11:   token.ParentL11,
		ParentL12:   token.ParentL12,
		ParentL13:   token.ParentL13,
		ParentL14:   token.ParentL14,
		ParentL15:   token.ParentL15,
		Level:       token.Level,
		IsSystem:    token.IsSystem,
		IsMain:      token.IsMain,
		RequireOtp:  false,

		SecretKey: token.SecretKey,
	}, tokenTTL)
	if err != nil {
		return err
	}

	return nil
}

func (s *SystemUser) ResetOtp(req *OtpResetParams) error {
	err := s.Update(&UpdateRequest{
		UpdatedBy:   req.UpdatedBy,
		UpdatedTime: req.UpdatedTime,
		Data: map[string]string{
			otp.OtpEnable: "false",
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *SystemUser) SkipOtp(req *OtpSkipParams) error {
	v := "false"
	if req.Value {
		v = "true"
	}
	err := s.Update(&UpdateRequest{
		UpdatedBy:   req.UpdatedBy,
		UpdatedTime: req.UpdatedTime,
		Data: map[string]string{
			otp.OtpSkip: v,
		},
	})
	if err != nil {
		return err
	}

	return nil
}
