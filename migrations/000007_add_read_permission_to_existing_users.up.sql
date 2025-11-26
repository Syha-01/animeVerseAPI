INSERT INTO users_permissions
SELECT id, (SELECT id FROM permissions WHERE code = 'animes:read')
FROM users
ON CONFLICT (user_id, permission_id) DO NOTHING;
