package service

import (
	"api-students/app/model"
	"api-students/helper"
)

// CanAccessStudent menentukan apakah user boleh mengakses
// data student tertentu.
//
// User boleh mengakses jika:
//
// 1. Student tersebut adalah miliknya sendiri.
// 2. Role user mempunyai permission anyPermission.
//
// Jika keduanya tidak terpenuhi, akses ditolak.
func CanAccessStudent(
	current model.AuthUser,
	ownerID int,
	perms *helper.PermissionSet,
	anyPermission string,
) bool {

	// Ownership selalu memberikan akses.
	if current.UserID == ownerID {
		return true
	}

	// Jika bukan pemilik, harus memiliki permission khusus.
	return perms.Can(
		current.Role,
		anyPermission,
	)
}