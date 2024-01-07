package models

type LoginResponse struct {
	Status int
	Body   LoginBodyResponse
}
type LoginBodyResponse struct {
	Token string
}

type RegisterResponse struct {
	Status int
	Body   RegisterBodyResponse
}
type RegisterBodyResponse struct {
	UserID int64
}
