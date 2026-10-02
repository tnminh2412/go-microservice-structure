package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
)

func hashTo32Bytes(input string) []byte {
	data := sha256.Sum256([]byte(input))
	return data[0:]
}

func encryptAES(key, data []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// create two 'windows' in to the output slice.
	output := make([]byte, aes.BlockSize+len(data))
	iv := output[:aes.BlockSize]
	encrypted := output[aes.BlockSize:]

	// populate the IV slice with random data.
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)

	// note that encrypted is still a window in to the output slice
	stream.XORKeyStream(encrypted, data)
	return output, nil
}

func decryptAES(key, data []byte) ([]byte, error) {
	// split the input up in to the IV seed and then the actual encrypted data.
	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(data, data)
	return data, nil
}

type Encrypter struct {
	secret string
}

func New(secret string) Encrypter {
	return Encrypter{secret: secret}
}

func (e Encrypter) GenerateToken(data string) (string, error) {
	hashKey := hashTo32Bytes(e.secret)
	encrypted, err := encryptAES(hashKey, []byte(data))
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(encrypted), nil
}

func (e Encrypter) DecryptToken(token string) (string, error) {
	hashKey := hashTo32Bytes(e.secret)
	encrypted, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return "", err
	}

	decrypted, err := decryptAES(hashKey, encrypted)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}
