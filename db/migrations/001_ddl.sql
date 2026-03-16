
CREATE EXTENSION IF NOT EXISTS postgis;


-- Профиль 
CREATE TABLE IF NOT EXISTS profiles ( 
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY, 
    phone TEXT NOT NULL, 
    first_name TEXT NOT NULL, 
    last_name TEXT NOT NULL, 
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), 
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), 
 
    CONSTRAINT phone_check CHECK (phone ~ '^(\+7|8)\s\d{3}\s\d{3}\s\d{2}\s\d{2}$'), 
    CONSTRAINT first_name_length_check CHECK ( LENGTH(first_name) < 40 ),
    CONSTRAINT last_name_length_check CHECK ( LENGTH(last_name) < 40 ) 
); 
 
 
-- Пользователи 
CREATE TABLE IF NOT EXISTS users ( 
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY, 
    email TEXT UNIQUE NOT NULL, 
    hashed_password TEXT, 
    provider TEXT, 
    provider_id TEXT UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), 
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), 
    salt TEXT, 
    profile_id BIGINT NOT NULL, 

    CONSTRAINT fk_profile FOREIGN KEY (profile_id) REFERENCES profiles(id), 
    CONSTRAINT email_check CHECK ( email ~ '^[^@]+@[^@]+\.[^@]+$' AND LENGTH(email)<150 ), 
    CONSTRAINT auth_check CHECK ( hashed_password IS NOT NULL OR provider IS NOT NULL) 
); 
--Рефреш токены пользователей 
CREATE TABLE IF NOT EXISTS refresh_tokens ( 
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY, 
    token_id UUID UNIQUE NOT NULL,
    user_id BIGINT REFERENCES users(id) NOT NULL, 
    expires_at TIMESTAMPTZ NOT NULL, 
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL 
); 
 
 
-- Города 
CREATE TABLE IF NOT EXISTS cities ( 
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY, 
    city_name TEXT NOT NULL UNIQUE, 
 
    CONSTRAINT city_name_length_check CHECK ( LENGTH(city_name) < 40 ) 
); 
 
 
-- Станции метро 
CREATE TABLE IF NOT EXISTS metro_stations ( 
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY, 
    station_name TEXT NOT NULL, 
    geo GEOGRAPHY(POINT, 4326) NOT NULL, 

    CONSTRAINT station_name_length_check CHECK ( LENGTH(station_name) < 40 ) 
); 
 
 
-- Тип помещения 
CREATE TABLE IF NOT EXISTS property_categories ( 
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY, 
    name TEXT NOT NULL, 
 
    CONSTRAINT name_length_check CHECK ( LENGTH(name) < 30 ) 
); 
 
CREATE TABLE IF NOT EXISTS utility_companies(
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY, 
    company_name TEXT NOT NULL,
    phone TEXT NOT NULL,
    geo GEOGRAPHY(POINT, 4326) NOT NULL, 
    address TEXT NOT NULL,
    avatar_url TEXT,
    alias TEXT UNIQUE NOT NULL,

    CONSTRAINT phone_check CHECK (phone ~ '^(\+7|8)\s\d{3}\s\d{3}\s\d{2}\s\d{2}$'), 
    CONSTRAINT address_check CHECK ( address ~ '^[а-яА-ЯёЁ\s\-\,\.\d\/]+$' AND LENGTH(address) >= 5 AND LENGTH(address) < 150 )
);
-- Дома 
CREATE TABLE IF NOT EXISTS buildings ( 
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY, 
    address TEXT NOT NULL, 
    geo GEOGRAPHY(POINT, 4326) NOT NULL, 
    city_id BIGINT NOT NULL, 
    metro_station_id BIGINT, 
    district TEXT, 
    floor_count SMALLINT NOT NULL,  
    company_id BIGINT, 
 
    CONSTRAINT fk_city FOREIGN KEY (city_id) REFERENCES cities(id), 
    CONSTRAINT fk_metro_station FOREIGN KEY (metro_station_id) REFERENCES metro_stations(id), 
    CONSTRAINT fk_company FOREIGN KEY (company_id) REFERENCES utility_companies(id), 
    CONSTRAINT address_check CHECK ( address ~ '^[а-яА-ЯёЁ\s\-\,\.\d\/]+$' AND LENGTH(address) >= 5 AND LENGTH(address) < 150 ), 
    CONSTRAINT district_length_check CHECK ( LENGTH(district) < 30 ), 
    CONSTRAINT floor_count_length_check CHECK ( floor_count < 100 ) 
); 
 
 
-- Объект недвижимости  
CREATE TABLE IF NOT EXISTS property ( 
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY, 
    category_id BIGINT, 
    building_id BIGINT, 
    area NUMERIC(10,2) NOT NULL, 
 
    CONSTRAINT fk_category FOREIGN KEY (category_id) REFERENCES property_categories(id), 
    CONSTRAINT area_check CHECK ( area > 0 ),
    CONSTRAINT fk_building FOREIGN KEY (building_id) REFERENCES buildings(id) 
); 
 
 
-- Объявления 
CREATE TABLE IF NOT EXISTS posters ( 
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY, 
    price NUMERIC(10,2) NOT NULL, 
    avatar_url TEXT, 
    description TEXT, 
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), 
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), 
    user_id BIGINT NOT NULL, 
    property_id BIGINT NOT NULL, 
    alias TEXT UNIQUE NOT NULL, 
    
 
   
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id), 
    CONSTRAINT fk_property FOREIGN KEY (property_id) REFERENCES property(id), 
    CONSTRAINT price_check CHECK ( price > 0 ), 
    CONSTRAINT description_length_check CHECK ( LENGTH(description) < 500 ) 
); 
 
 
-- Объявление-Фото
CREATE TABLE IF NOT EXISTS poster_photos ( 
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY, 
    img_url TEXT NOT NULL, 
    sequence_order SMALLINT, 
    poster_id BIGINT NOT NULL, 
 
    CONSTRAINT fk_poster FOREIGN KEY (poster_id) REFERENCES posters(id), 
    CONSTRAINT sequence_order_check CHECK (sequence_order > 0 AND sequence_order < 16) 
); 
 
 
-- Категория квартиры 
CREATE TABLE IF NOT EXISTS flat_categories ( 
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY, 
    name TEXT NOT NULL 
); 
 
 
-- Квартира 
CREATE TABLE IF NOT EXISTS flat ( 
    property_id BIGINT PRIMARY KEY, 
    floor SMALLINT, 
    number INT, 
    category_id BIGINT, 
 
    CONSTRAINT fk_property FOREIGN KEY (property_id) REFERENCES property(id), 
    CONSTRAINT fk_flat_category FOREIGN KEY (category_id) REFERENCES flat_categories(id),
    CONSTRAINT floor_check CHECK ( floor > 0 ),
    CONSTRAINT number_check CHECK ( number > 0 )
); 
 

 -- ЖК компании 


-- ЖК-Фото
CREATE TABLE IF NOT EXISTS utility_companies_photos ( 
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY, 
    img_url TEXT NOT NULL, 
    sequence_order SMALLINT, 
    utility_company_id BIGINT NOT NULL, 
 
    CONSTRAINT fk_utility_company  FOREIGN KEY (utility_company_id) REFERENCES utility_companies(id), 
    CONSTRAINT sequence_order_check CHECK (sequence_order > 0 AND sequence_order < 16) 
); 
 


 -- Основные FK индексы
CREATE INDEX idx_users_profile_id ON users(profile_id);
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_buildings_city_id ON buildings(city_id);
CREATE INDEX idx_buildings_metro_station_id ON buildings(metro_station_id);
CREATE INDEX idx_buildings_company_id ON buildings(company_id);
CREATE INDEX idx_property_category_id ON property(category_id);
CREATE INDEX idx_property_building_id ON property(building_id);
CREATE INDEX idx_posters_user_id ON posters(user_id);
CREATE INDEX idx_posters_property_id ON posters(property_id);
CREATE INDEX idx_poster_photos_poster_id ON poster_photos(poster_id);
CREATE INDEX idx_flat_property_id ON flat(property_id);
CREATE INDEX idx_flat_category_id ON flat(category_id);
CREATE INDEX idx_utility_companies_photos_company_id ON utility_companies_photos(utility_company_id);

-- Поисковые индексы
CREATE INDEX idx_posters_price ON posters(price);
CREATE INDEX idx_property_area ON property(area);
CREATE INDEX idx_posters_alias ON posters(alias);
CREATE INDEX idx_posters_created_at ON posters(created_at);

-- Гео-индексы (GiST для GEOGRAPHY)
CREATE INDEX idx_metro_stations_geo ON metro_stations USING GIST(geo);
CREATE INDEX idx_buildings_geo ON buildings USING GIST(geo);
CREATE INDEX idx_utility_companies_geo ON utility_companies USING GIST(geo);



-- Функция-триггер
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Триггеры для таблиц с updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_profiles_updated_at BEFORE UPDATE ON profiles FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_posters_updated_at BEFORE UPDATE ON posters FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
