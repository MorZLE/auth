package models

// Response RESTAPI
type Response struct {
	Status int
	Body   interface{}
}

// LoginBodyResponse body LoginResponse
type LoginBodyResponse struct {
	Token string
}

// RegisterBodyResponse body RegisterResponse
type RegisterBodyResponse struct {
	UserID int64
}

// AddAppBodyResponse body AddAppResponse
type AddAppBodyResponse struct {
	AppID int32
}

type CreateAdminBodyResponse struct {
	AdminID int64
}

type DeleteAdminBodyResponse struct {
	Result bool
}

type IsAdminBodyResponse struct {
	Result bool
	LVL    int32
}
