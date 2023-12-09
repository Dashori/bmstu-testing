package models

type Client struct {
	ClientId uint64
	Login    string
	Password string
	Email    string
	OTP      string
}
