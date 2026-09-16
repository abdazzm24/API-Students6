package service

import (
	"strings"

	"api-students/app/model"
	"api-students/helper"
)

// CanAccessUser menentukan apakah current user boleh
// mengakses data user target.
//
// Ada dua jalur:
//
// 1. User mengakses dirinya sendiri.
// 2. Role memiliki permission anyPermission.
func CanAccessUser(
	current model.AuthUser,
	targetID int,
	perms *helper.PermissionSet,
	anyPermission string,
) bool {

	if current.UserID == targetID {
		return true
	}

	return perms.Can(
		current.Role,
		anyPermission,
	)
}

// ValidateAssignRole memeriksa aturan perubahan role.
//
// Role harus dikenal oleh PermissionSet dan user tidak boleh
// mengubah role dirinya sendiri.
func ValidateAssignRole(
	current model.AuthUser,
	targetID int,
	req model.AssignRoleRequest,
	perms *helper.PermissionSet,
) map[string]string {

	errs := map[string]string{}

	role := strings.TrimSpace(req.Role)

	if role == "" {
		errs["role"] = "wajib diisi"
		return errs
	}

	if !perms.IsKnownRole(role) {
		errs["role"] =
			"role tidak dikenal, pilih salah satu dari: " +
				strings.Join(perms.KnownRoles(), ", ")
	}

	if current.UserID == targetID {
		errs["role"] =
			"tidak boleh mengubah role diri sendiri"
	}

	return errs
}