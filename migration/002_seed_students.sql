INSERT INTO students (
    nim,
    name,
    grade,
    is_active
)
VALUES
    ('434241118', 'Abdullah Azzam', 85, TRUE),
    ('434241119', 'Diaul Haq', 78, TRUE),
    ('434241200', 'Anshari Shidqi', 92, FALSE)
ON CONFLICT DO NOTHING;