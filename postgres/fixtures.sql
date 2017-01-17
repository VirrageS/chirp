-- users --
INSERT INTO users (username, email, password, name, created_at, last_login)
VALUES ('admin', 'admin@admin.com', '$2a$10$WuBMX8PWRTfOdTHcMRZ.Uuqp7Q10gRvN2UoJ/NYPN80RSj83BfL/q', 'admin', now(), now());
INSERT INTO users (username, email, password, name, created_at, last_login)
VALUES ('corpsegrinder', 'corpsegrinder@cannibalcorpse.com', '$2a$10$0nWJkUXR3nlxINaLxbcgl.2Lk4frwDShVnyhjfRGS5CcukiBHpQvi', 'George Fisher', now(), now());

-- tweets --
INSERT INTO tweets (author_id, created_at, content)
VALUES (1, now(), 'tweet');
