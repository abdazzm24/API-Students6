-- MIGRATION 004
-- ROLE BASED ACCESS CONTROL 

BEGIN;
 
-- 1. TABEL ROLES

CREATE TABLE IF NOT EXISTS roles (
    name VARCHAR(50) PRIMARY KEY,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Data awal role
INSERT INTO roles (name, description)
VALUES
    ('admin', 'Administrator dengan akses penuh'),
    ('staff', 'Petugas dengan akses terbatas'),
    ('user', 'Pengguna biasa')
ON CONFLICT (name) DO NOTHING;

-- 2. TABEL PERMISSIONS

CREATE TABLE IF NOT EXISTS permissions (
    name VARCHAR(100) PRIMARY KEY,
    description TEXT
);

-- Permission untuk user
INSERT INTO permissions (name, description)
VALUES
    ('user:list', 'Melihat daftar user'),
    ('user:read:any', 'Melihat data user siapa pun'),
    ('user:update:any', 'Mengubah data user siapa pun'),
    ('user:delete', 'Menghapus user'),
    ('role:assign', 'Mengubah role user')
ON CONFLICT (name) DO NOTHING;

-- Permission untuk student
INSERT INTO permissions (name, description)
VALUES
    ('student:list', 'Melihat daftar student'),
    ('student:read:any', 'Melihat data student siapa pun'),
    ('student:create', 'Membuat student baru'),
    ('student:update:any', 'Mengubah data student siapa pun'),
    ('student:delete', 'Menghapus student')
ON CONFLICT (name) DO NOTHING;

-- 3. TABEL RELASI ROLE DAN PERMISSION

CREATE TABLE IF NOT EXISTS role_permissions (
    role_name VARCHAR(50) NOT NULL,
    permission_name VARCHAR(100) NOT NULL,

    PRIMARY KEY (role_name, permission_name),

    CONSTRAINT fk_role_permissions_role
        FOREIGN KEY (role_name)
        REFERENCES roles(name)
        ON DELETE CASCADE,

    CONSTRAINT fk_role_permissions_permission
        FOREIGN KEY (permission_name)
        REFERENCES permissions(name)
        ON DELETE CASCADE
);

-- 4. PERMISSION UNTUK ADMIN

INSERT INTO role_permissions (role_name, permission_name)
VALUES
    ('admin', 'user:list'),
    ('admin', 'user:read:any'),
    ('admin', 'user:update:any'),
    ('admin', 'user:delete'),
    ('admin', 'role:assign'),

    ('admin', 'student:list'),
    ('admin', 'student:read:any'),
    ('admin', 'student:create'),
    ('admin', 'student:update:any'),
    ('admin', 'student:delete')
ON CONFLICT (role_name, permission_name) DO NOTHING;

-- 5. PERMISSION UNTUK STAFF

INSERT INTO role_permissions (role_name, permission_name)
VALUES
    ('staff', 'user:list'),
    ('staff', 'user:read:any'),

    ('staff', 'student:list'),
    ('staff', 'student:read:any'),
    ('staff', 'student:create')
ON CONFLICT (role_name, permission_name) DO NOTHING;

-- 6. ROLE USER TIDAK DIBERI PERMISSION ADMIN

-- Role user sengaja tidak diberi permission khusus.

-- 7. NORMALISASI ROLE USER

-- Jika ada data user dengan role kosong atau tidak dikenal,
-- ubah terlebih dahulu menjadi user.

UPDATE users
SET role = 'user'
WHERE role IS NULL
   OR role NOT IN ('admin', 'staff', 'user');

-- 8. TAMBAHKAN FOREIGN KEY ROLE PADA USERS

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_users_role'
          AND table_name = 'users'
    ) THEN
        ALTER TABLE users
        ADD CONSTRAINT fk_users_role
        FOREIGN KEY (role)
        REFERENCES roles(name);
    END IF;
END $$;

COMMIT;