CREATE TABLE IF NOT EXISTS TODOS(
  id bigserial PRIMARY KEY,
  item text NOT NULL,
  description text NOT NULL,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  status BOOLEAN NOT NULL
)