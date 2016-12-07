-- users --
INSERT INTO users (username, email, password, name, created_at, last_login)
VALUES ('admin', 'admin@admin.com', 'admin', 'admin', now(), now());
INSERT INTO users (username, email, password, name, created_at, last_login)
VALUES ('corpsegrinder', 'corpsegrinder@cannibalcorpse.com', 'fuckthealliance', 'George Fisher', now(), now());

-- tweets --
INSERT INTO tweets (author_id, created_at, content)
VALUES (1, now(), 'tweet');
