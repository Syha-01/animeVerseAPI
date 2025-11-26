INSERT INTO users_permissions
VALUES (
    (SELECT id FROM users WHERE email = 'test@example.com'),
    (SELECT id FROM permissions WHERE code = 'animes:write')
);
