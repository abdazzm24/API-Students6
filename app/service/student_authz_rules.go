package service

import "api-students/helper"

// StudentAuthUser adalah identitas user yang digunakan
// untuk pemeriksaan authorization.
//
// Struktur ini dibuat terpisah agar aturan authorization
// tidak bergantung pada Fiber atau repository.
type StudentAuthUser struct {
	ID   int
	Role string
}

// CanAccessStudent menentukan apakah user boleh mengakses
// sebuah data student.
//
// Aturan:
//
// 1. Jika user memiliki permission anyPermission,
//    user boleh mengakses student siapa pun.
//
// 2. Jika user tidak memiliki permission tersebut,
//    user hanya boleh mengakses student yang owner_id-nya
//    sama dengan ID user.
//
// 3. Jika tidak memenuhi kedua aturan tersebut,
//    akses ditolak.
func CanAccessStudent(
	current StudentAuthUser,
	ownerID int,
	permissions *helper.PermissionSet,
	anyPermission string,
) bool {
	// Admin atau role lain yang memiliki permission
	// untuk mengakses semua student diperbolehkan.
	if permissions != nil &&
		permissions.Can(current.Role, anyPermission) {
		return true
	}

	// Jika tidak memiliki akses semua data,
	// user hanya boleh mengakses data miliknya sendiri.
	return current.ID == ownerID
}