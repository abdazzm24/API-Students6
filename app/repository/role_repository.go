package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RoleRepository bertanggung jawab mengambil data role
// dan permission dari database.
type RoleRepository interface {
	LoadPermissions(
		ctx context.Context,
	) (map[string][]string, error)
}

// roleRepository adalah implementasi RoleRepository.
type roleRepository struct {
	pool *pgxpool.Pool
}

// NewRoleRepository membuat repository role baru.
func NewRoleRepository(
	pool *pgxpool.Pool,
) RoleRepository {
	return &roleRepository{
		pool: pool,
	}
}

// LoadPermissions mengambil seluruh role dan permission.
//
// Role yang belum memiliki permission tetap dimasukkan
// ke dalam hasil map dengan array permission kosong.
func (r *roleRepository) LoadPermissions(
	ctx context.Context,
) (map[string][]string, error) {
	const query = `
		SELECT
			r.name AS role_name,
			rp.permission_name
		FROM roles r
		LEFT JOIN role_permissions rp
			ON rp.role_name = r.name
		ORDER BY r.name, rp.permission_name
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf(
			"load role permissions query: %w",
			err,
		)
	}
	defer rows.Close()

	result := make(map[string][]string)

	for rows.Next() {
		var (
			roleName       string
			permissionName *string
		)

		err := rows.Scan(
			&roleName,
			&permissionName,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"scan role permissions: %w",
				err,
			)
		}

		if _, exists := result[roleName]; !exists {
			result[roleName] = []string{}
		}

		if permissionName != nil {
			result[roleName] = append(
				result[roleName],
				*permissionName,
			)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"iterate role permissions: %w",
			err,
		)
	}

	return result, nil
}