package helper

import "sync"

// PermissionSet menyimpan daftar permission berdasarkan role.
//
// Contoh:
//
// admin:
//   student:list
//   student:create
//   student:delete
//
// staff:
//   student:list
//   student:create
type PermissionSet struct {
	mu     sync.RWMutex
	byRole map[string]map[string]struct{}
}

// NewPermissionSet membuat PermissionSet baru dari data role dan permission.
func NewPermissionSet(
	permissionsByRole map[string][]string,
) *PermissionSet {
	result := &PermissionSet{
		byRole: make(map[string]map[string]struct{}),
	}

	for role, permissions := range permissionsByRole {
		result.byRole[role] = make(map[string]struct{})

		for _, permission := range permissions {
			result.byRole[role][permission] = struct{}{}
		}
	}

	return result
}

// Can memeriksa apakah suatu role memiliki permission tertentu.
//
// Sistem menggunakan prinsip fail closed.
// Artinya, jika role atau permission tidak ditemukan,
// hasilnya false.
func (p *PermissionSet) Can(
	role string,
	permission string,
) bool {
	if p == nil {
		return false
	}

	p.mu.RLock()
	defer p.mu.RUnlock()

	rolePermissions, exists := p.byRole[role]
	if !exists {
		return false
	}

	_, exists = rolePermissions[permission]

	return exists
}

// PermissionsOf mengembalikan seluruh permission dari sebuah role.
func (p *PermissionSet) PermissionsOf(
	role string,
) []string {
	if p == nil {
		return []string{}
	}

	p.mu.RLock()
	defer p.mu.RUnlock()

	rolePermissions, exists := p.byRole[role]
	if !exists {
		return []string{}
	}

	result := make([]string, 0, len(rolePermissions))

	for permission := range rolePermissions {
		result = append(result, permission)
	}

	return result
}

// KnownRoles mengembalikan daftar role yang diketahui sistem.
func (p *PermissionSet) KnownRoles() []string {
	if p == nil {
		return []string{}
	}

	p.mu.RLock()
	defer p.mu.RUnlock()

	result := make([]string, 0, len(p.byRole))

	for role := range p.byRole {
		result = append(result, role)
	}

	return result
}

// IsKnownRole memeriksa apakah role terdaftar dalam sistem.
func (p *PermissionSet) IsKnownRole(role string) bool {
	if p == nil {
		return false
	}

	p.mu.RLock()
	defer p.mu.RUnlock()

	_, exists := p.byRole[role]

	return exists
}