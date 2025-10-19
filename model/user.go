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
