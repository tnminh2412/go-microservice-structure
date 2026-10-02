package mask

var (
	maskChar   = byte('*')
	length     = 8
	maskLength = 4
)

func Username(username string) string {
	lastCharacters := username
	if len(username) > maskLength {
		lastCharacters = username[len(username)-maskLength:]
	}
	mLen := length - len(lastCharacters)
	m := make([]byte, mLen)
	for i := range m {
		m[i] = maskChar
	}
	return string(m) + lastCharacters
}
