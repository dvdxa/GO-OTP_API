package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strconv"
)

func GenerateOTP(message string, secretKey string) (int, error) {

	secretKeyBytes := []byte(secretKey)

	hash := hmac.New(sha256.New, secretKeyBytes)
	hash.Write([]byte(message))
	hashBytes := hash.Sum(nil)

	hexDigest := hex.EncodeToString(hashBytes)

	sixDigitCode, err := strconv.ParseInt(hexDigest[:5], 16, 0)
	if err != nil {
		return 0, err
	}

	return int(sixDigitCode), nil
}

// Convert struct to map
func StructToMap(data interface{}) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var resultMap map[string]interface{}
	err = json.Unmarshal(jsonData, &resultMap)
	if err != nil {
		return nil, err
	}

	return resultMap, nil
}

// Convert map to struct
func MapToStruct(data map[string]string, result interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonData, result)
	if err != nil {
		return err
	}

	return nil
}
