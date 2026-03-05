CREATE OR REPLACE FUNCTION refresh_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS posters (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    price NUMERIC(12, 3) NOT NULL CHECK(price > 0),
    image_url TEXT CHECK (LENGTH(image_url) < 2048),
    address TEXT NOT NULL CHECK (LENGTH(address) < 255),
    metro TEXT CHECK (LENGTH(metro) < 255),
    area NUMERIC(10, 4) NOT NULL CHECK(area > 0),
    floor SMALLINT NOT NULL CHECK(floor > 0),
    type TEXT NOT NULL CHECK(LENGTH(type) < 255),
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL
);

DROP TRIGGER IF EXISTS update_posters_updated_at ON posters;

CREATE TRIGGER update_posters_updated_at
    BEFORE UPDATE ON posters
    FOR EACH ROW
    EXECUTE FUNCTION refresh_updated_at();