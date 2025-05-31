CREATE TABLE "records" (
  "id" bigserial PRIMARY KEY,
  "username" varchar NOT NULL,
  "content" jsonb NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now()
);