package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
)

var (
	jm           = &jsonpb.Marshaler{EmitDefaults: true, OrigName: true}
	AESSecretKey = []byte("cprZwqj4lqvbAAf2mNz7RvNKiQ3CeH8h")
)

func MD5Encrypted(s string) string {
	alg := md5.New()
	alg.Write([]byte(s))
	encba := alg.Sum(nil)
	return fmt.Sprintf("%x", encba)
}

func PasswordEncrypted(pwd string) string {
	alg := md5.New()
	alg.Write([]byte(pwd))
	encba := alg.Sum(nil)
	return base64.StdEncoding.EncodeToString(encba)
}

func GetTimeString(value *time.Time) string {
	if value == nil || value.IsZero() {
		return ""
	}
	return value.Format("2006-01-02T15:04:05.000")
}

func ParseStringToInt32(s string) int32 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return int32(i)
}

func ParseInt32ToString(i int32) string {
	return fmt.Sprintf("%d", i)
}

func ParseStringToInt64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func ParseInt64ToString(i int64) string {
	return fmt.Sprintf("%d", i)
}

func ParseStringToFloat64(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func ParseFloat64ToString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func ParseStringToBool(s string) bool {
	b, _ := strconv.ParseBool(s)
	return b
}

func ParseStringToTime(dateFormat, str string) (*time.Time, error) {
	if str == "" {
		return nil, nil
	}

	t, err := time.Parse(dateFormat, str)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func ParseTimeToString(dateFormat string, t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(dateFormat)
}

func ParseInterfaceToString(in interface{}) (string, error) {
	if protoMsg, ok := in.(proto.Message); ok {
		result, err := jm.MarshalToString(protoMsg)
		if err != nil {
			return "", err
		}
		return string(result), nil
	} else if strMsg, ok := in.(string); ok {
		return strMsg, nil
	} else if byteMsg, ok := in.([]byte); ok {
		return string(byteMsg), nil
	} else {
		result, err := json.Marshal(in)
		if err != nil {
			return "", err
		}
		return string(result), nil
	}
}

func AESEncrypt(plaintext string) (string, error) {
	c, err := aes.NewCipher(AESSecretKey)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	return hex.EncodeToString(ciphertext), nil
}

func AESDecrypt(ct string) (string, error) {
	data, err := hex.DecodeString(ct)
	if err != nil {
		return "", err
	}

	c, err := aes.NewCipher(AESSecretKey)
	if err != nil {
		fmt.Printf("\n\n no c\n\n")
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		fmt.Printf("\n\n no gcm\n\n")
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("The length of ciphertext cannot smaller than nonceSize")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		fmt.Printf("\n\n no plaintext\n\n")
		return "", err
	}

	return string(plaintext), nil
}
