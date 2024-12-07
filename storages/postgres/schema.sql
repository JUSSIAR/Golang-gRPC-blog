CREATE TYPE "entity_type" AS ENUM (
  'post',
  'comment'
);

CREATE TABLE "user" (
  "user_id" varchar PRIMARY KEY
);

CREATE TABLE "post" (
  "post_id" uuid PRIMARY KEY,
  "author_id" varchar,
  "title" varchar,
  "text" text,
  "timestamp" timestamp,
  "like_count" int
);

CREATE TABLE "comment" (
  "comment_id" uuid PRIMARY KEY,
  "author_id" varchar,
  "post_id" uuid,
  "text" text,
  "timestamp" timestamp,
  "like_count" int
);

CREATE TABLE "relation" (
  "id" int PRIMARY KEY,
  "author_id" varchar,
  "post_id" uuid,
  "comment_id" uuid,
  "entity_type" entity_type,
  "value" bool
);

ALTER TABLE "post" ADD FOREIGN KEY ("author_id") REFERENCES "user" ("user_id");

ALTER TABLE "comment" ADD FOREIGN KEY ("author_id") REFERENCES "user" ("user_id");

ALTER TABLE "relation" ADD FOREIGN KEY ("author_id") REFERENCES "user" ("user_id");

ALTER TABLE "comment" ADD FOREIGN KEY ("post_id") REFERENCES "post" ("post_id");

ALTER TABLE "relation" ADD FOREIGN KEY ("post_id") REFERENCES "post" ("post_id");

ALTER TABLE "relation" ADD FOREIGN KEY ("comment_id") REFERENCES "comment" ("comment_id");
