-- +goose Up
CREATE TABLE metadata(
  date TEXT PRIMARY KEY,
  explanation TEXT,
  hdurl TEXT,
  media_type TEXT,
  service_version TEXT,
  title TEXT,
  url TEXT
);


--+goose Down
DROP TABLE metadata;