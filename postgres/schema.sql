CREATE TABLE users (
  id               SERIAL PRIMARY KEY,
  twitter_token    VARCHAR(255) UNIQUE,
  facebook_token   VARCHAR(255) UNIQUE,
  google_token     VARCHAR(255) UNIQUE,

  username         VARCHAR(255) UNIQUE,
  email            VARCHAR(255) UNIQUE DEFAULT '',
  password         VARCHAR(512), -- <algorithm>$<iterations>$<salt>$<hash>
  created_at       TIMESTAMP NOT NULL DEFAULT now(),
  last_login       TIMESTAMP,
  active           BOOLEAN NOT NULL DEFAULT TRUE,

  name             VARCHAR(255) DEFAULT '',
  avatar_url       VARCHAR(1024) DEFAULT '',

  CONSTRAINT proper_email CHECK (email ~* '^[A-Za-z0-9._%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$')
);

CREATE INDEX users_idx ON users (id);
CREATE INDEX users_twitter_token_idx ON users (twitter_token);
CREATE INDEX users_facebook_token_idx ON users (facebook_token);
CREATE INDEX users_google_token_idx ON users (google_token);
CREATE INDEX users_username_idx ON users (username);
CREATE INDEX users_email_idx ON users (email);
CREATE INDEX users_active_idx ON users (active);


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
  author_id   INTEGER REFERENCES users (id) ON DELETE CASCADE,
  created_at  TIMESTAMP NOT NULL DEFAULT now(),
  content     VARCHAR(150) NOT NULL
);

CREATE INDEX tweets_idx ON tweets (id);


CREATE TABLE likes (
  tweet_id INTEGER REFERENCES tweets (id),
  user_id  INTEGER REFERENCES users (id),
  liked_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX likes_tweets_idx ON likes (tweet_id);
CREATE INDEX likes_users_idx ON likes (user_id);
CREATE INDEX likes_idx ON likes (tweet_id, user_id, liked_at);


CREATE TABLE retweets (
  tweet_id     INTEGER REFERENCES tweets (id),
  user_id      INTEGER REFERENCES users (id),
  retweeted_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX retweets_tweets_idx ON retweets (tweet_id);
CREATE INDEX retweets_users_idx ON retweets (user_id);
CREATE INDEX retweets_idx ON retweets (tweet_id, user_id, retweeted_at);


CREATE TABLE tags (
  id        SERIAL PRIMARY KEY,
  name      VARCHAR(150) UNIQUE NOT NULL
);

CREATE UNIQUE INDEX tags_lowercase_name_idx ON tags ((lower(name)));


CREATE TABLE tweets_tags (
  tweet_id  INTEGER REFERENCES tweets (id),
  tag_id    INTEGER REFERENCES tags (id),

  PRIMARY KEY (tweet_id, tag_id)
);

CREATE INDEX tweets_tags_tweets_idx ON tweets_tags (tweet_id);
CREATE INDEX tweets_tags_tags_idx ON tweets_tags (tag_id);
CREATE INDEX tweets_tags_idx ON tweets_tags (tweet_id, tag_id);
