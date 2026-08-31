package service

import (
	"strings"

	"api-students/app/model"
)

func ValidateCreateRequest(
	req model.CreateStudentRequest,
) map[string]string {

	errors := make(map[string]string)

	if strings.TrimSpace(req.NIM) == "" {
		errors["nim"] = "NIM wajib diisi"
	}

	if strings.TrimSpace(req.Name) == "" {
		errors["name"] = "nama wajib diisi"
	}

	if req.Grade < 0 || req.Grade > 100 {
		errors["grade"] =
			"grade harus berada pada rentang 0 sampai 100"
	}

	return errors
}

func ValidateReplaceRequest(
	req model.ReplaceStudentRequest,
) map[string]string {

	errors := make(map[string]string)

	if strings.TrimSpace(req.NIM) == "" {
		errors["nim"] = "NIM wajib diisi"
	}

	if strings.TrimSpace(req.Name) == "" {
		errors["name"] = "nama wajib diisi"
	}

	if req.Grade < 0 || req.Grade > 100 {
		errors["grade"] =
			"grade harus berada pada rentang 0 sampai 100"
	}

	return errors
}

func ValidatePatchRequest(
	req model.PatchStudentRequest,
) map[string]string {

	errors := make(map[string]string)

	if req.NIM != nil &&
		strings.TrimSpace(*req.NIM) == "" {

		errors["nim"] = "NIM tidak boleh kosong"
	}

	if req.Name != nil &&
		strings.TrimSpace(*req.Name) == "" {

		errors["name"] = "nama tidak boleh kosong"
	}

	if req.Grade != nil &&
		(*req.Grade < 0 || *req.Grade > 100) {

		errors["grade"] =
			"grade harus berada pada rentang 0 sampai 100"
	}

	return errors
}

func CalculateTotalPages(
	limit int,
	total int,
) int {

	if total == 0 {
		return 0
	}

	return (total + limit - 1) / limit
}