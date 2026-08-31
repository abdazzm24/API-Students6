package main

import "api-students/repository"

type Student = repository.Student
type CreateStudentRequest = repository.CreateStudentRequest
type ReplaceStudentRequest = repository.ReplaceStudentRequest
type PatchStudentRequest = repository.PatchStudentRequest
type StudentRepository = repository.StudentRepository

type WebResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Meta    *Meta  `json:"meta,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

type Meta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type ListQuery struct {
	Page     int
	Limit    int
	Search   string
	Sort     string
	Order    string
	IsActive *bool
}
