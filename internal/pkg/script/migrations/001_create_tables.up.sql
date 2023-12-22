CREATE TYPE "role_type" AS ENUM (
  'ADMIN'
);

CREATE TABLE "users" (
  "id" serial PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "password" varchar NOT NULL,
  "role" role_type NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "created_by" int,
  "updated_at" timestamp,
  "updated_by" int,
  "deleted_at" timestamp,
  "deleted_by" int
);

CREATE TABLE "persons" (
  "id" serial PRIMARY KEY,
  "full_name" varchar NOT NULL,
  "avatar" varchar,
  "info" jsonb,
  "status" bool DEFAULT true,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "created_by" int,
  "updated_at" timestamp,
  "updated_by" int,
  "deleted_at" timestamp,
  "deleted_by" int
);

CREATE TABLE "places" (
  "id" serial PRIMARY KEY,
  "name" varchar NOT NULL,
  "code" varchar,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "created_by" int,
  "updated_at" timestamp,
  "updated_by" int,
  "deleted_at" timestamp,
  "deleted_by" int
);

CREATE TABLE "meetings" (
  "id" serial PRIMARY KEY,
  "title" jsonb NOT NULL,
  "description" jsonb NOT NULL,
  "meeting_time" timestamp,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "created_by" int,
  "updated_at" timestamp,
  "updated_by" int,
  "deleted_at" timestamp,
  "deleted_by" int
);

CREATE TABLE "meeting_places" (
  "id" serial PRIMARY KEY,
  "meeting_id" int,
  "person_id" int,
  "place_id" int,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "created_by" int,
  "updated_at" timestamp,
  "updated_by" int,
  "deleted_at" timestamp,
  "deleted_by" int
);

CREATE TABLE "materials" (
  "id" serial PRIMARY KEY,
  "index" int NOT NULL,
  "title" jsonb NOT NULL,
  "content" jsonb,
  "file" jsonb,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "created_by" int,
  "updated_at" timestamp,
  "updated_by" int,
  "deleted_at" timestamp,
  "deleted_by" int
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

ALTER TABLE "meeting_places" ADD FOREIGN KEY ("meeting_id") REFERENCES "meetings" ("id");

ALTER TABLE "meeting_places" ADD FOREIGN KEY ("person_id") REFERENCES "persons" ("id");

ALTER TABLE "meeting_places" ADD FOREIGN KEY ("place_id") REFERENCES "places" ("id");

ALTER TABLE "meeting_places" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");

ALTER TABLE "meeting_places" ADD FOREIGN KEY ("updated_by") REFERENCES "users" ("id");

ALTER TABLE "meeting_places" ADD FOREIGN KEY ("deleted_by") REFERENCES "users" ("id");

ALTER TABLE "materials" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");

ALTER TABLE "materials" ADD FOREIGN KEY ("updated_by") REFERENCES "users" ("id");

ALTER TABLE "materials" ADD FOREIGN KEY ("deleted_by") REFERENCES "users" ("id");