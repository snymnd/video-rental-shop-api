package entity

import "time"

type Users struct {
	ID        string
	Name      string
	Email     string
	Password  string
	Role      int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type Login struct {
	ID       string
	Name     string
	Email    string
	Role     int
	Password string
	Token    string
}
