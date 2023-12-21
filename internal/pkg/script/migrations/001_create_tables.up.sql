CREATE TYPE "role_type" AS ENUM (
  'ADMIN'
);

CREATE TABLE "users" (
  "id" uuid PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "password" varchar NOT NULL,
  "role" role_type NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "created_by" uuid,
  "updated_at" timestamp,
  "updated_by" uuid,
  "deleted_at" timestamp,
  "deleted_by" uuid
);

CREATE TABLE "persons" (
  "id" serial PRIMARY KEY,
  "full_name" varchar NOT NULL,
  "avatar" varchar,
  "info" jsonb,
  "status" bool DEFAULT true,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "created_by" uuid,
  "updated_at" timestamp,
  "updated_by" uuid,
  "deleted_at" timestamp,
  "deleted_by" uuid
);

CREATE TABLE "places" (
  "id" uuid PRIMARY KEY,
  "name" varchar NOT NULL,
  "code" varchar,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "created_by" uuid,
  "updated_at" timestamp,
  "updated_by" uuid,
  "deleted_at" timestamp,
  "deleted_by" uuid
);

CREATE TABLE "meetings" (
  "id" uuid PRIMARY KEY,
  "title" jsonb NOT NULL,
  "description" jsonb NOT NULL,
  "meeting_time" timestamp,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "created_by" uuid,
  "updated_at" timestamp,
  "updated_by" uuid,
  "deleted_at" timestamp,
  "deleted_by" uuid
);

ALTER TABLE "users" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");

ALTER TABLE "users" ADD FOREIGN KEY ("updated_by") REFERENCES "users" ("id");

ALTER TABLE "users" ADD FOREIGN KEY ("deleted_by") REFERENCES "users" ("id");

ALTER TABLE "persons" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");

ALTER TABLE "persons" ADD FOREIGN KEY ("updated_by") REFERENCES "users" ("id");

ALTER TABLE "persons" ADD FOREIGN KEY ("deleted_by") REFERENCES "users" ("id");

ALTER TABLE "places" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");

ALTER TABLE "places" ADD FOREIGN KEY ("updated_by") REFERENCES "users" ("id");

ALTER TABLE "places" ADD FOREIGN KEY ("deleted_by") REFERENCES "users" ("id");

ALTER TABLE "meetings" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");

ALTER TABLE "meetings" ADD FOREIGN KEY ("updated_by") REFERENCES "users" ("id");

ALTER TABLE "meetings" ADD FOREIGN KEY ("deleted_by") REFERENCES "users" ("id");