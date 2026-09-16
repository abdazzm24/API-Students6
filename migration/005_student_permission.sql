-- MIGRATION 005
-- STUDENT AUTHORIZATION

BEGIN;

-- 1. TAMBAHKAN OWNER_ID PADA TABEL STUDENTS

ALTER TABLE students
ADD COLUMN IF NOT EXISTS owner_id INTEGER;

-- 2. ISI OWNER_ID UNTUK DATA LAMA

-- Jika tabel students sudah memiliki data lama,
-- owner_id perlu diisi terlebih dahulu.
--
-- Untuk sementara, data lama diarahkan kepada user admin.
--
-- Pastikan user dengan role admin sudah tersedia.
UPDATE students
SET owner_id = (
    SELECT id
    FROM users
    WHERE role = 'admin'
    ORDER BY id
    LIMIT 1
)
WHERE owner_id IS NULL;

-- 3. TAMBAHKAN FOREIGN KEY OWNER_ID

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_students_owner'
          AND table_name = 'students'
    ) THEN
        ALTER TABLE students
        ADD CONSTRAINT fk_students_owner
        FOREIGN KEY (owner_id)
        REFERENCES users(id)
        ON DELETE RESTRICT;
    END IF;
END $$;

-- 4. INDEX OWNER_ID

CREATE INDEX IF NOT EXISTS idx_students_owner_id
ON students(owner_id);

COMMIT;