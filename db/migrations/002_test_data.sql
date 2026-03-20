-- profiles
INSERT INTO profiles (phone, first_name, last_name) VALUES
('8 999 111 22 33', 'Иван', 'Иванов');

-- users
INSERT INTO users (email, hashed_password, salt, profile_id) VALUES
('ivan@example.com', 'hashed_password_stub', 'salt_stub', 1);

-- cities
INSERT INTO cities (city_name) VALUES
('Москва');

-- metro_stations
INSERT INTO metro_stations (station_name, geo) VALUES
('Арбатская', ST_GeogFromText('SRID=4326;POINT(37.6045 55.7522)'));

-- property_categories
INSERT INTO property_categories (name) VALUES
('Квартира');

-- utility_companies
INSERT INTO utility_companies (company_name, phone, geo, address, avatar_url, alias) VALUES
(
    'ПИК',
    '8 495 123 45 67',
    ST_GeogFromText('SRID=4326;POINT(37.6000 55.7500)'),
    'Москва, ул. Арбат, 1',
    'https://example.com/company.png',
    'pik'
);

-- buildings
INSERT INTO buildings (address, geo, city_id, metro_station_id, district, floor_count, company_id) VALUES
(
    'Москва, ул. Арбат, 10',
    ST_GeogFromText('SRID=4326;POINT(37.6050 55.7525)'),
    1,
    1,
    'Арбат',
    12,
    1
);

-- property
INSERT INTO property (category_id, building_id, area) VALUES
(1, 1, 45.50);

-- flat_categories
INSERT INTO flat_categories (name) VALUES
('1-комнатная');

-- flat
INSERT INTO flat (property_id, floor, number, category_id) VALUES
(1, 5, 25, 1);

-- posters
INSERT INTO posters (price, avatar_url, description, user_id, property_id, alias) VALUES
(
    12500000.00,
    'https://example.com/poster-main.jpg',
    'Светлая квартира рядом с метро',
    1,
    1,
    'kvartira-na-arbate'
);

-- poster_photos
INSERT INTO poster_photos (img_url, sequence_order, poster_id) VALUES
('https://example.com/poster-1.jpg', 1, 1),
('https://example.com/poster-2.jpg', 2, 1),
('https://example.com/poster-3.jpg', 3, 1);