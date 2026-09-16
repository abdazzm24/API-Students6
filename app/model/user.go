package model

import "time"

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateUserRequest digunakan oleh POST /users.
//
// Role sengaja tidak ada di sini agar client tidak bisa
// melakukan mass assignment menjadi admin.
type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ReplaceUserRequest digunakan oleh PUT /users/:id.
type ReplaceUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	IsActive bool   `json:"is_active"`
}

// PatchUserRequest digunakan oleh PATCH /users/:id.
type PatchUserRequest struct {
	Username *string `json:"username,omitempty"`
	Email    *string `json:"email,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}

// AssignRoleRequest digunakan oleh PATCH /users/:id/role.
type AssignRoleRequest struct {
	Role string `json:"role"`
}