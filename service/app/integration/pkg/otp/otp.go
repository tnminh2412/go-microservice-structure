package otp

const (
	OtpEnable   = "otpEnable"
	OtpSkip     = "otpSkip"
	OtpSecret   = "otpSecret"
	OtpProvider = "otpProvider"
)

const (
	GoogleOtpProvider = "google"
)

func IsOtpSkip(data map[string]string) bool {
	v, ok := data[OtpSkip]
	if !ok {
		return true
	}
	if v == "true" {
		return true
	}
	return false
}

func IsOtpEnable(data map[string]string) bool {
	if IsOtpSkip(data) {
		return false
	}

	v, ok := data[OtpEnable]
	if !ok {
		return false
	}
	if v == "true" {
		return true
	}
	return false
}
