CREATE TABLE users (
  id               SERIAL PRIMARY KEY,
  twitter_token    VARCHAR(255) UNIQUE,
  facebook_token   VARCHAR(255) UNIQUE,
  google_token     VARCHAR(255) UNIQUE,

  name             VARCHAR(255) DEFAULT '',
  username         VARCHAR(255) UNIQUE,
  email            VARCHAR(255) UNIQUE DEFAULT '',
  password         VARCHAR(512), -- <algorithm>$<iterations>$<salt>$<hash>
  created_at       TIMESTAMP NOT NULL DEFAULT now(),
  last_login       TIMESTAMP,
  active           BOOLEAN NOT NULL DEFAULT TRUE,

  CONSTRAINT proper_email CHECK (email ~* '^[A-Za-z0-9._%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$')
);

CREATE INDEX users_idx ON users (id);
CREATE INDEX users_twitter_token_idx ON users (twitter_token);
CREATE INDEX users_facebook_token_idx ON users (facebook_token);
CREATE INDEX users_google_token_idx ON users (google_token);
CREATE INDEX users_username_idx ON users (username);
CREATE INDEX users_email_idx ON users (email);



CREATE TABLE followers (
  follower   INTEGER REFERENCES users (id) ON DELETE CASCADE,
  following  INTEGER REFERENCES users (id) ON DELETE CASCADE,

  PRIMARY KEY (follower, following)
);

CREATE INDEX followers_follower_idx ON followers (follower);
CREATE INDEX followers_following_idx ON followers (following);
CREATE INDEX followers_idx ON followers (follower, following);



CREATE TABLE tweets (
  id          SERIAL PRIMARY KEY,
  user_id     INTEGER REFERENCES users (id) NOT NULL ON DELETE CASCADE,
  created_at  TIMESTAMP NOT NULL DEFAULT now(),
  content     VARCHAR(150) NOT NULL
);

CREATE INDEX posts_idx ON tweets (id);
CREATE INDEX posts_fulltext_idx ON tweets USING GIN (to_tsvector(content));



CREATE TABLE tags (
  id        SERIAL PRIMARY KEY,
  name      VARCHAR(150) UNIQUE NOT NULL
);

CREATE UNIQUE INDEX tags_lowercase_name_idx ON tags ((lower(name)));



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
