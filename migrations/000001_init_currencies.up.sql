CREATE TABLE IF NOT EXISTS rates (
    unique_id     serial PRIMARY KEY,
    currency_id   int,
	name          TINYTEXT NOT NULL,
	abbreviation  TINYTEXT NOT NULL,
	rate          DOUBLE PRECISION,
	inserted_at   DATE
);