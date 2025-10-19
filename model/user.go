package model

import "encoding/xml"

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type GetUserByIDRequest struct {
	XMLName xml.Name `xml:"urn:user-service GetUserByID"`
	ID      int      `xml:"id"`
}

type GetUserByIDResponse struct {
	XMLName xml.Name `xml:"urn:user-service GetUserByIDResponse"`
	User    User     `xml:"User"`
}

// CreateUser Operation
type CreateUserRequest struct {
	XMLName xml.Name `xml:"urn:user-service CreateUser"`
	Name    string   `xml:"name"`
	Email   string   `xml:"email"`
}

type CreateUserResponse struct {
	XMLName xml.Name `xml:"urn:user-service CreateUserResponse"`
	User    User     `xml:"User"`
}

// UpdateUser Operation
type UpdateUserRequest struct {
	XMLName xml.Name `xml:"urn:user-service UpdateUser"`
	ID      int      `xml:"id"`
	Name    string   `xml:"name"`
	Email   string   `xml:"email"`
}

type UpdateUserResponse struct {
	XMLName xml.Name `xml:"urn:user-service UpdateUserResponse"`
	User    User     `xml:"User"`
}

// DeleteUser Operation
type DeleteUserRequest struct {
	XMLName xml.Name `xml:"urn:user-service DeleteUser"`
	ID      int      `xml:"id"`
}

type DeleteUserResponse struct {
	XMLName xml.Name `xml:"urn:user-service DeleteUserResponse"`
	Success bool     `xml:"success"`
	Message string   `xml:"message"`
}
