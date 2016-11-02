CREATE TABLE users (
  id               SERIAL PRIMARY KEY,
  twitter_token    VARCHAR(255) UNIQUE,
  facebook_token   VARCHAR(255) UNIQUE,
  google_token     VARCHAR(255) UNIQUE,

  name             VARCHAR(255),
  username         VARCHAR(255) UNIQUE NOT NULL,
  email            VARCHAR(255) UNIQUE NOT NULL, -- should be not null? (we get email from twitter?)
  password         VARCHAR(512) NOT NULL, -- <algorithm>$<iterations>$<salt>$<hash>
  created_on       TIMESTAMP NOT NULL,
  last_login       TIMESTAMP,

  CONSTRAINT proper_email CHECK (email ~* '^[A-Za-z0-9._%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$')
);

CREATE INDEX users_idx ON users (id);


CREATE TABLE followers (
  first_user   INTEGER REFERENCES users (id) ON DELETE CASCADE, -- should be named (first_user_id)?
  second_user  INTEGER REFERENCES users (id) ON DELETE CASCADE, -- should be named (second_user_id)?

  PRIMARY KEY (first_user, second_user),
  CONSTRAINT avoid_duplication CHECK (first_user < second_user)
);

CREATE INDEX followers_first_user_idx ON followers (first_user);
CREATE INDEX followers_second_user_idx ON followers (second_user);
CREATE INDEX followers_idx ON followers (first_user, second_user);


CREATE TABLE posts ( -- should be named (tweets) or (chirps)?
  id       SERIAL PRIMARY KEY,
  user_id  INTEGER REFERENCES users (id) ON DELETE CASCADE,
  post     VARCHAR(150) NOT NULL -- should be named (content) ?
);

CREATE INDEX posts_idx ON posts (id);
CREATE INDEX posts_fulltext_idx ON posts USING GIN (to_tsvector('english', post));


CREATE TABLE tags (
  id        SERIAL PRIMARY KEY,
  name      VARCHAR(150) NOT NULL
);


CREATE TABLE posts_tags (
  post_id  INTEGER REFERENCES posts (id),
  tag_id  INTEGER REFERENCES tags (id),

  PRIMARY KEY (post_id, tag_id)
);

CREATE INDEX posts_tags_posts_idx ON posts_tags (post_id);
CREATE INDEX posts_tags_tags_idx ON posts_tags (tag_id);
CREATE INDEX posts_tags_idx ON posts_tags (post_id, tag_id);


/*
--- GENERAL IDEA IN NOSQL DATABASE ---

TABLE likes (
  posts.id, user.id, 1
)

TABLE shares (
  posts.id, user.id, 1
)

*/
