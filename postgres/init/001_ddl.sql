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
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS users (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    email TEXT UNIQUE CHECK (LENGTH(email) < 255)  NOT NULL ,
    hashed_password TEXT CHECK (LENGTH(hashed_password) BETWEEN 8 AND 255), -- делаю null чтобы можно было добавить oauth
    salt TEXT CHECK (LENGTH(salt) < 255),  -- делаю null чтобы можно было добавить oauth
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now()
);


CREATE TABLE IF NOT EXISTS refresh_tokens (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    token_id UUID  NOT NULL ,
    user_id BIGINT REFERENCES users(id) NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL
);

ALTER TABLE users ADD CONSTRAINT email_check CHECK (email ~* '^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$');

CREATE OR REPLACE FUNCTION refresh_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


DROP TRIGGER IF EXISTS update_users_updated_at ON users;

DROP TRIGGER IF EXISTS update_posters_updated_at ON posters;

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION refresh_updated_at();

CREATE TRIGGER update_posters_updated_at
    BEFORE UPDATE ON posters
    FOR EACH ROW
    EXECUTE FUNCTION refresh_updated_at();