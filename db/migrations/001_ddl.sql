-- Города
CREATE TABLE IF NOT EXISTS cities (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    city_name TEXT NOT NULL UNIQUE,

    CONSTRAINT city_name_length_check CHECK ( LENGTH(city_name) < 40 )
);


-- Станции метро
CREATE TABLE IF NOT EXISTS metro_stations (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    station_name TEXT NOT NULL UNIQUE,

    CONSTRAINT station_name_length_check CHECK ( LENGTH(station_name) < 40 )
);


-- ЖК Компании
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE IF NOT EXISTS utility_companies (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    company_name TEXT NOT NULL,
    contacts TEXT CHECK (
        contacts ~ '^\+?[\d\s\-\(\)]{10,20}$|^$'
    ),
    geo GEOGRAPHY(POINT, 4326),
    address TEXT CHECK (
        address ~ '^[а-яА-ЯёЁ\s\-\,\.\d\/]+$' 
        AND LENGTH(address) >= 5 
        AND LENGTH(address) < 150
    ),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    city_id BIGINT,
    metro_station_id BIGINT,

    CONSTRAINT fk_city FOREIGN KEY (city_id) REFERENCES cities(id),
    CONSTRAINT fk_metro_station FOREIGN KEY (metro_station_id) REFERENCES metro_stations(id),
    CONSTRAINT company_name_length_check CHECK ( LENGTH(company_name) < 40 ),
    CONSTRAINT contacts_length_check CHECK ( LENGTH(contacts) < 50 )
);



-- Пользователи
CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    email TEXT NOT NULL,
    hashed_password TEXT,
    provider TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    salt TEXT NOT NULL,
    company_id BIGINT,

    CONSTRAINT fk_company FOREIGN KEY (company_id) REFERENCES utility_companies(id),
    CONSTRAINT email_check CHECK ( email ~ '^[^@]+@[^@]+\.[^@]+$' AND LENGTH(email)<150 ),
    CONSTRAINT auth_check CHECK ( hashed_password IS NOT NULL OR provider IS NOT NULL)
);
--Рефреш токены пользователей
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    token_id UUID  NOT NULL ,
    user_id BIGINT REFERENCES users(id) NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL
);


-- Дома
CREATE TABLE IF NOT EXISTS buildings (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    geo GEOGRAPHY(POINT, 4326) NOT NULL,
    address TEXT NOT NULL,
    district TEXT,
    company_id BIGINT,
    city_id BIGINT,
    metro_station_id BIGINT,

    CONSTRAINT fk_city FOREIGN KEY (city_id) REFERENCES cities(id),
    CONSTRAINT fk_metro_station FOREIGN KEY (metro_station_id) REFERENCES metro_stations(id),
    CONSTRAINT fk_company FOREIGN KEY (company_id) REFERENCES utility_companies(id),
    CONSTRAINT address_length_check CHECK ( LENGTH(address) < 150 ),
    CONSTRAINT district_length_check CHECK ( LENGTH(district) < 30 )
);


-- Тип помещения
CREATE TABLE IF NOT EXISTS apartment_categories (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name TEXT NOT NULL,
    description TEXT,

    CONSTRAINT name_length_check CHECK ( LENGTH(name) < 30 ),
    CONSTRAINT description_length_check CHECK ( LENGTH(description) < 500 )
);


-- Помощение
CREATE TABLE IF NOT EXISTS apartments (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    floor SMALLINT NOT NULL,
    number SMALLINT NOT NULL,
    building_id BIGINT NOT NULL,
    category_id BIGINT,

    CONSTRAINT fk_building FOREIGN KEY (building_id) REFERENCES buildings(id),
    CONSTRAINT fk_category FOREIGN KEY (category_id) REFERENCES apartment_categories(id),
    CONSTRAINT floor_check CHECK (floor > 0),
    CONSTRAINT number_check CHECK (number > 0)
);


-- Объявления
CREATE TABLE IF NOT EXISTS posters (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    title TEXT NOT NULL,
    price NUMERIC(10,2) NOT NULL,
    avatar_url TEXT,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    user_id BIGINT NOT NULL,
    apartment_id BIGINT NOT NULL,

    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_apartment FOREIGN KEY (apartment_id) REFERENCES apartments(id),
    CONSTRAINT title_length_check CHECK ( LENGTH(title) < 100 ),
    CONSTRAINT price_check CHECK ( price > 0 ),
    CONSTRAINT description_length_check CHECK ( LENGTH(description) < 500 )
);


-- История изменения цены
CREATE TABLE IF NOT EXISTS price_history (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    price NUMERIC(10,2) NOT NULL,
    changed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    poster_id BIGINT NOT NULL,

    CONSTRAINT price_check CHECK ( price > 0 ),
    CONSTRAINT fk_poster FOREIGN KEY (poster_id) REFERENCES posters(id)
);


-- Лайки
CREATE TABLE IF NOT EXISTS likes (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    user_id BIGINT NOT NULL,
    poster_id BIGINT NOT NULL,

    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_poster FOREIGN KEY (poster_id) REFERENCES posters(id)
);


-- Просмотры
CREATE TABLE IF NOT EXISTS views (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    user_id BIGINT NOT NULL,
    poster_id BIGINT NOT NULL,

    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_poster FOREIGN KEY (poster_id) REFERENCES posters(id)
);


-- Объявление-Фото
CREATE TABLE IF NOT EXISTS poster_photos (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    img_url TEXT,
    sequence_order SMALLINT,
    poster_id BIGINT NOT NULL,

    CONSTRAINT fk_poster FOREIGN KEY (poster_id) REFERENCES posters(id),
    CONSTRAINT sequence_order_check CHECK (sequence_order > 0)
);