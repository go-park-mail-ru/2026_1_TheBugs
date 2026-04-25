
CREATE EXTENSION IF NOT EXISTS postgis;

 
CREATE TABLE IF NOT EXISTS profiles ( 
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY, 
    phone TEXT NOT NULL, 
    first_name TEXT NOT NULL, 
    last_name TEXT NOT NULL, 
    avatar_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), 
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), 
 
    CONSTRAINT phone_check CHECK (phone ~ '^(\+7|8)\s\d{3}\s\d{3}\s\d{2}\s\d{2}$'), 
    CONSTRAINT first_name_length_check CHECK ( LENGTH(first_name) <= 40 ),
    CONSTRAINT last_name_length_check CHECK ( LENGTH(last_name) <= 40 ) 
); 
COMMENT ON TABLE profiles IS 'Профиль';

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
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,

    CONSTRAINT fk_profile FOREIGN KEY (profile_id) REFERENCES profiles(id), 
    CONSTRAINT email_check CHECK ( email ~ '^[^@]+@[^@]+\.[^@]+$' AND LENGTH(email)<= 255 ), 
    CONSTRAINT auth_check CHECK ( hashed_password IS NOT NULL OR provider IS NOT NULL) 
); 

COMMENT ON TABLE users IS 'Пользователи';

CREATE TABLE IF NOT EXISTS refresh_tokens ( 
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY, 
    token_id UUID UNIQUE NOT NULL,
    user_id BIGINT REFERENCES users(id) NOT NULL, 
    expires_at TIMESTAMPTZ NOT NULL, 
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL 
); 
COMMENT ON TABLE refresh_tokens IS 'Рефреш токены пользователей';
 

CREATE TABLE IF NOT EXISTS cities ( 
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY, 
    city_name TEXT NOT NULL UNIQUE, 
 
    CONSTRAINT city_name_length_check CHECK ( LENGTH(city_name) <= 40 ) 
); 
COMMENT ON TABLE cities IS 'Города';
 
 

CREATE TABLE IF NOT EXISTS metro_stations ( 
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY, 
    station_name TEXT NOT NULL, 
    geo GEOGRAPHY(POINT, 4326) NOT NULL, 

    CONSTRAINT station_name_length_check CHECK ( LENGTH(station_name) <= 40 ) 
); 

COMMENT ON TABLE metro_stations IS 'Станции метро';
 
 
CREATE TABLE IF NOT EXISTS property_categories ( 
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY, 
    name TEXT NOT NULL, 
    alias TEXT UNIQUE NOT NULL,
 
    CONSTRAINT name_length_check CHECK ( LENGTH(name) <= 30 ),
    CONSTRAINT alias_length_check CHECK ( LENGTH(alias) <= 50 )
); 
COMMENT ON TABLE property_categories IS 'Тип помещения';

CREATE TABLE IF NOT EXISTS facilities(
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY, 
    name TEXT NOT NULL, 
    alias TEXT UNIQUE NOT NULL,

    CONSTRAINT name_length_check CHECK ( LENGTH(name) <= 30 ),
    CONSTRAINT alias_length_check CHECK ( LENGTH(alias) <= 50 ) 
);

COMMENT ON TABLE facilities IS 'Удобства';



CREATE TABLE IF NOT EXISTS developers (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    developer_name TEXT NOT NULL,
    avatar_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT developer_name_length_check CHECK (LENGTH(developer_name) <= 100)
);

COMMENT ON TABLE developers IS 'Застройщики';

 
CREATE TABLE IF NOT EXISTS utility_companies(
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY, 
    company_name TEXT NOT NULL,
    phone TEXT NOT NULL,
    geo GEOGRAPHY(POINT, 4326) NOT NULL, 
    address TEXT NOT NULL,
    alias TEXT UNIQUE NOT NULL,
    description TEXT,
    avatar_url TEXT,
    developer_id BIGINT,

    CONSTRAINT phone_check CHECK (phone ~ '^(\+7|8)\s\d{3}\s\d{3}\s\d{2}\s\d{2}$'), 
    CONSTRAINT address_check CHECK ( LENGTH(address) >= 5 AND LENGTH(address) <= 500 ),
    CONSTRAINT fk_developer FOREIGN KEY (developer_id) REFERENCES developers(id),
    CONSTRAINT description_length_check CHECK ( LENGTH(description) <= 3000 )
);
COMMENT ON TABLE utility_companies IS 'Компании услуг/ЖК';

CREATE TABLE IF NOT EXISTS buildings ( 
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY, 
    address TEXT NOT NULL, 
    geo GEOGRAPHY(POINT, 4326) NOT NULL, 
    city_id BIGINT NOT NULL, 
    --city TEXT,
    metro_station_id BIGINT, 
    district TEXT, 
    floor_count SMALLINT NOT NULL,  
    company_id BIGINT, 
 
    CONSTRAINT fk_city FOREIGN KEY (city_id) REFERENCES cities(id), 
    CONSTRAINT fk_metro_station FOREIGN KEY (metro_station_id) REFERENCES metro_stations(id), 
    --CONSTRAINT city_check LENGTH(city) < 30 , 
    CONSTRAINT fk_company FOREIGN KEY (company_id) REFERENCES utility_companies(id), 
    CONSTRAINT address_check CHECK ( LENGTH(address) >= 5 AND LENGTH(address) <= 500 ), 
    CONSTRAINT district_length_check CHECK ( LENGTH(district) <= 100 ), 
    CONSTRAINT floor_count_length_check CHECK (floor_count > 0 AND floor_count <= 99999 ) 
); 
 
COMMENT ON TABLE buildings IS 'Дома';

CREATE TABLE IF NOT EXISTS property ( 
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY, 
    category_id BIGINT NOT NULL, 
    building_id BIGINT NOT NULL UNIQUE, 
    area NUMERIC(10,2) NOT NULL, 
 
    CONSTRAINT fk_category FOREIGN KEY (category_id) REFERENCES property_categories(id), 
    CONSTRAINT area_check CHECK ( area > 0 ),
    CONSTRAINT fk_building FOREIGN KEY (building_id) REFERENCES buildings(id) 
); 
COMMENT ON TABLE property IS 'Объект недвижимости';

CREATE TABLE IF NOT EXISTS facility_property(
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY, 
    property_id BIGINT,
    facility_id BIGINT,
    CONSTRAINT fk_property_id FOREIGN KEY (property_id) REFERENCES property(id),
    CONSTRAINT fk_facility_id FOREIGN KEY (facility_id) REFERENCES facilities(id),
    CONSTRAINT unique_ids UNIQUE (property_id, facility_id)
);

COMMENT ON TABLE facilities IS 'Удобства-Недвижемость';

 

CREATE TABLE IF NOT EXISTS posters ( 
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY, 
    price NUMERIC(10,2) NOT NULL, 
    avatar_url TEXT, 
    description TEXT, 
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), 
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), 
    deleted_at TIMESTAMPTZ, 
    user_id BIGINT NOT NULL, 
    property_id BIGINT, 
    alias TEXT UNIQUE NOT NULL, 
    
 
   
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id), 
    CONSTRAINT fk_property FOREIGN KEY (property_id) REFERENCES property(id), 
    CONSTRAINT price_check CHECK ( price > 0 ), 
    CONSTRAINT description_length_check CHECK ( LENGTH(description) <= 3000 ) 
); 
COMMENT ON TABLE posters IS 'Объявления';
 

CREATE TABLE IF NOT EXISTS poster_photos ( 
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY, 
    img_url TEXT NOT NULL, 
    sequence_order SMALLINT, 
    poster_id BIGINT NOT NULL, 
 
    CONSTRAINT fk_poster FOREIGN KEY (poster_id) REFERENCES posters(id), 
    CONSTRAINT sequence_order_check CHECK (sequence_order >= 0 AND sequence_order <= 12) 
); 

COMMENT ON TABLE poster_photos IS 'Фото объявления';
 
 
CREATE TABLE IF NOT EXISTS flat_categories ( 
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY, 
    name TEXT NOT NULL,
    room_count SMALLINT,
    CONSTRAINT room_count_check CHECK (room_count BETWEEN 0 AND 6) 
); 

COMMENT ON TABLE flat_categories IS 'Категории квартиры';
 
 
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

COMMENT ON TABLE flat IS 'Квартира';


CREATE TABLE IF NOT EXISTS utility_companies_photos ( 
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY, 
    img_url TEXT NOT NULL, 
    sequence_order SMALLINT, 
    utility_company_id BIGINT NOT NULL, 
 
    CONSTRAINT fk_utility_company  FOREIGN KEY (utility_company_id) REFERENCES utility_companies(id), 
    CONSTRAINT sequence_order_check CHECK (sequence_order > 0 AND sequence_order < 16) 
); 

COMMENT ON TABLE utility_companies_photos IS 'Фото ЖК';


CREATE TABLE IF NOT EXISTS favorites (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    user_id BIGINT NOT NULL,
    poster_id BIGINT NOT NULL,

    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_poster FOREIGN KEY (poster_id) REFERENCES posters(id),
    CONSTRAINT unique_user_poster_favorite UNIQUE (user_id, poster_id)
);

COMMENT ON TABLE favorites IS 'Избранные объявления';

CREATE TABLE IF NOT EXISTS views (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    user_id BIGINT NOT NULL,
    poster_id BIGINT NOT NULL,

    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_poster FOREIGN KEY (poster_id) REFERENCES posters(id)
);

COMMENT ON TABLE views IS 'Просмотры';

CREATE TYPE handling_status AS ENUM ('sent', 'in_progress', 'finished');

CREATE TABLE IF NOT EXISTS handling_categories (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name TEXT NOT NULL,

    CONSTRAINT handling_category_name_length_check CHECK (LENGTH(name) <= 100)
);

COMMENT ON TABLE handling_categories IS 'Категории обращений';


CREATE TABLE IF NOT EXISTS handlings (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id BIGINT NOT NULL,
    category_id BIGINT NOT NULL,
    description TEXT NOT NULL,
    status handling_status NOT NULL DEFAULT 'sent',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    admin_id BIGINT,

    CONSTRAINT fk_handling_user
        FOREIGN KEY (user_id) REFERENCES users(id),

    CONSTRAINT fk_handling_admin
        FOREIGN KEY (admin_id) REFERENCES users(id),

    CONSTRAINT fk_handling_category
        FOREIGN KEY (category_id) REFERENCES handling_categories(id),

    CONSTRAINT handling_description_length_check
        CHECK (LENGTH(description) <= 3000)
);

COMMENT ON TABLE handling_categories IS 'Обращения';


CREATE TABLE IF NOT EXISTS handling_photos (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    img_url TEXT,
    sequence_order SMALLINT,
    handling_id BIGINT NOT NULL,

    CONSTRAINT fk_handling_photo_handling
        FOREIGN KEY (handling_id) REFERENCES handlings(id),

    CONSTRAINT handling_photo_sequence_order_check
        CHECK (sequence_order >= 0 AND sequence_order <= 12)
); 

COMMENT ON TABLE handling_categories IS 'Фото обращения';

CREATE INDEX idx_favorites_users_id ON favorites(user_id);


CREATE TABLE IF NOT EXISTS chats_users (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id BIGINT NOT NULL,
    chat_id BIGINT NOT NULL,

    CONSTRAINT fk_chats_users_user FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_chats_users_chat FOREIGN KEY (chat_id) REFERENCES chats(id)
);
 

CREATE INDEX idx_likes_users_id ON likes(user_id);
CREATE INDEX idx_views_users_id ON views(user_id);

CREATE INDEX idx_favorites_posters_id ON favorites(poster_id);
CREATE INDEX idx_views_posters_id ON views(poster_id);


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
CREATE INDEX idx_utility_companies_developer_id ON utility_companies(developer_id);

-- Поисковые индексы
CREATE INDEX idx_posters_price ON posters(price);
CREATE INDEX idx_property_area ON property(area);
CREATE INDEX idx_posters_alias ON posters(alias);
CREATE INDEX idx_posters_created_at ON posters(created_at);

CREATE INDEX IF NOT EXISTS idx_posters_deleted_at ON posters(deleted_at) WHERE deleted_at IS NOT NULL;

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

CREATE OR REPLACE FUNCTION get_explain_rows(p_sql text)
 RETURNS bigint AS $BODY$
 DECLARE
    plan_text text;
    plan_json jsonb;
    rows_count bigint := 0;
    v_actual text;
 BEGIN
    EXECUTE 'EXPLAIN (FORMAT JSON) ' || p_sql INTO plan_text;
    plan_json := plan_text::jsonb;
    v_actual := plan_json->0->'Plan'->>'Plan Rows';
    IF v_actual IS NOT NULL AND v_actual <> '' THEN
        rows_count := v_actual::bigint;
        RETURN rows_count;
    END IF;
     SELECT substring(plan_text from 'rows=(\d+)')::bigint INTO rows_count;
    RETURN COALESCE(rows_count, 0);
 END;
 $BODY$ LANGUAGE plpgsql;

