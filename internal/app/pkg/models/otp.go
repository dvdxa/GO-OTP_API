package models

//incoming request
type Request struct {
	ID       string `json:"id"`
	Account  string `json:"account"`
	Value    string `json:"value"` //otp
	Lifetime *int64 `json:"lifetime"`
}

//response for incoming request
type Response struct {
	ID        string `json:"id"`
	Account   string `json:"account"`
	State     int64  `json:"state"`
	CreatedAt string `json:"created_at"`
	ExpiredAt string `json:"expired_at"`
	LifeTime  int64  `json:"lifetime"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Hash struct {
	Account       string `json:"account"`
	State         string `json:"state"`
	CreatedAt     string `json:"createdAt"`
	ExpiredAt     string `json:"expiredAt"`
	LifeTime      string `json:"lifetime"`
	ValidateLimit string `json:"validate_limit"`
	BanTime       string `json:"ban_time"`
}
