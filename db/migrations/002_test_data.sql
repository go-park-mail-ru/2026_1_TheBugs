INSERT INTO profiles (phone, first_name, last_name)
VALUES
('+7 900 111 22 33', 'Иван', 'Иванов'),
('+7 901 222 33 44', 'Мария', 'Петрова'),
('+7 902 333 44 55', 'Сергей', 'Сидоров');

INSERT INTO users (email, hashed_password, provider, profile_id, salt)
VALUES
('ivan@example.com', 'hashed_pass_1', NULL, 1, 'salt1'),
('maria@example.com', 'hashed_pass_2', NULL, 2, 'salt2'),
('sergey@example.com', NULL, 'google', 3, NULL);

INSERT INTO cities (city_name)
VALUES
('Москва'),
('Санкт-Петербург'),
('Новосибирск');


INSERT INTO metro_stations (station_name, geo)
VALUES
('Киевская', ST_GeogFromText('SRID=4326;POINT(37.561432 55.743831)')),
('Пушкинская', ST_GeogFromText('SRID=4326;POINT(37.605678 55.764321)')),
('Невский проспект', ST_GeogFromText('SRID=4326;POINT(30.329439 59.934280)'));


INSERT INTO property_categories (name)
VALUES
('Квартира'),
('Дом'),
('Апартаменты');


INSERT INTO utility_companies (company_name, phone, geo, address, avatar_url, alias)
VALUES
('Энергокомфорт', '+7 999 123 45 67', ST_GeogFromText('SRID=4326;POINT(37.618423 55.751244)'), 'г. Москва, ул. Тверская, д. 7', NULL, 'energocomfort'),
('Водоканал-Сервис', '+7 495 987 65 43', ST_GeogFromText('SRID=4326;POINT(37.611987 55.758102)'), 'г. Москва, ул. Лесная, д. 15', NULL, 'vodokanal_service'),
('КомфортЭко', '+7 926 555 12 34', ST_GeogFromText('SRID=4326;POINT(37.603123 55.752345)'), 'г. Москва, пр-т Садовый, д. 10', NULL, 'komfort_eco');


INSERT INTO buildings (address, geo, city_id, metro_station_id, district, floor_count, company_id)
VALUES
('ул. Арбат, д. 12', ST_GeogFromText('SRID=4326;POINT(37.586123 55.752987)'), 1, 1, 'Арбат', 10, 1),
('пр. Невский, д. 20', ST_GeogFromText('SRID=4326;POINT(30.330123 59.934123)'), 2, 3, 'Центр', 15, 2);


INSERT INTO property (category_id, building_id, area)
VALUES
(1, 1, 55.5),
(2, 2, 30.0);


INSERT INTO posters (price, avatar_url, description, user_id, property_id, alias)
VALUES
(5000000, NULL, 'Просторная квартира с видом на парк', 1, 1, 'poster1'),
(3000000, NULL, 'Уютная студия в центре города', 2, 2, 'poster2');


INSERT INTO poster_photos (img_url, sequence_order, poster_id)
VALUES
('https://example.com/photo1.jpg', 1, 1),
('https://example.com/photo2.jpg', 2, 1),
('https://example.com/photo3.jpg', 1, 2);


INSERT INTO flat_categories (name)
VALUES
('Апартаменты'),
('Студия'),
('1 комнатная'),
('2 комнатная'),
('3 комнатная'),
('4 комнатная');

INSERT INTO flat (property_id, floor, number, category_id)
VALUES
(1, 5, 12, 2),
(2, 7, 34, 1);


INSERT INTO utility_companies_photos (img_url, sequence_order, utility_company_id)
VALUES
('https://example.com/company1.jpg', 1, 1),
('https://example.com/company2.jpg', 2, 1),
('https://example.com/company3.jpg', 1, 2);
