package models

type User struct {
    FirstName string `json:"firstName"`
    LastName  string `json:"lastName"`
    Email     string `json:"email"`
    Phone     string `json:"phone,omitempty"`
}
