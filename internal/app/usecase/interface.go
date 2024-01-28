package usecase

import "otp_api/internal/app/pkg/models"

type IService interface {
	CreateOTP(request *models.Request, secretKey string, validateLimit int) (*models.Response, *models.ErrorResponse)
	GetOTP(id string, secretKey string) (*models.Response, *models.ErrorResponse)
	ValidateOTP(id string, otp string, secretKey string, retryAfter int, banDur int) (*models.Response, *models.ErrorResponse)
}

type IStorage interface {
	CreateOTP(id string, hash map[string]interface{}) error
	GetOTP(key string) (map[string]string, error)
	UpdateOTPField(hashKey string, key string, value interface{}) error
	BanAccountForTime(account string, field string, duration string) error
	GetAccount(account string) (map[string]string, error)
	RemoveAccountFromBan(account string) error
	UpdateMultipleFields(hashKey string, fields map[string]interface{}) error
}

type IAdapter interface {
	SendOTPToClient(account string, otp string) error
}
