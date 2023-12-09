package models

import "time"

type ErrorBody struct {
	Err string `json:"error"`
}

type Pets struct {
	Pets []struct {
		PetId    uint64 `json:"PetId"`
		Name     string `json:"Name"`
		Type     string `json:"Type"`
		Age      uint64 `json:"Age"`
		Health   uint64 `json:"Health"`
		ClientId uint64 `json:"ClientId"`
	} `json:"pets"`
}

type Records struct {
	Records []struct {
		RecordId      uint64    `json:"RecordId"`
		PetId         uint64    `json:"PetId"`
		ClientId      uint64    `json:"ClientId"`
		DoctorId      uint64    `json:"DoctorId"`
		DatetimeStart time.Time `json:"DatetimeStart"`
		DatetimeEnd   time.Time `json:"DatetimeEnd"`
	} `json:"records"`
}

type Client struct {
	ClientId uint64 `json:"ClientId"`
	Login    string `json:"Login"`
	Password string
	Token    string `json:"Token"`
	Email    string `json:"Email"`
	OTP      string `json:"OTP"`
}

type Pet struct {
	PetId    uint64 `json:"PetId"`
	Name     string `json:"Name"`
	Type     string `json:"Type"`
	Age      uint64 `json:"Age"`
	Health   uint64 `json:"Health"`
	ClientId uint64 `json:"ClientId"`
}
