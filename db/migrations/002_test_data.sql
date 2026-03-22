-- ============================================================
-- 1. Профили пользователей
-- ============================================================
INSERT INTO profiles (phone, first_name, last_name) VALUES
    ('+7 495 123 45 67', 'Иван', 'Петров'),
    ('+7 495 987 65 43', 'Анна', 'Соколова'),
    ('+7 812 111 22 33', 'Сергей', 'Кузьмин'),
    ('+7 812 222 33 44', 'Ольга', 'Морозова'),
    ('+7 843 444 55 66', 'Дмитрий', 'Волков'),
    ('+7 843 555 66 77', 'Мария', 'Новикова');

-- ============================================================
-- 2. Пользователи
-- ============================================================
INSERT INTO users (email, hashed_password, provider, profile_id, salt) VALUES
    ('ivan.petrov@mail.ru', '$2b$12$abcdefgh1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMN', NULL, 1, 'salt_ivan_001'),
    ('anna.sokolova@yandex.ru', '$2b$12$zyxwvutsrqponmlkjihgfedcbaZYXWVUTSRQPONMLKJIHGFEDCBA0123', NULL, 2, 'salt_anna_002'),
    ('sergey.kuzmin@gmail.com', NULL, 'google', 3, 'salt_sergey_003'),
    ('olga.morozova@mail.ru', '$2b$12$1234567890abcdefgh1234567890abcdefghijklmnopqrstuvwxyzABC', NULL, 4, 'salt_olga_004'),
    ('dmitry.volkov@yandex.ru', NULL, 'vk', 5, 'salt_dmitry_005'),
    ('maria.novikova@gmail.com', '$2b$12$qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM12', NULL, 6, 'salt_maria_006');

-- ============================================================
-- 3. Города
-- ============================================================
INSERT INTO cities (city_name) VALUES
    ('Москва'),
    ('Санкт-Петербург'),
    ('Казань'),
    ('Новосибирск'),
    ('Екатеринбург');

-- ============================================================
-- 4. Станции метро
-- ============================================================
INSERT INTO metro_stations (station_name, geo) VALUES
    ('Арбатская', ST_GeogFromText('SRID=4326;POINT(37.5806 55.7495)')),
    ('Тверская', ST_GeogFromText('SRID=4326;POINT(37.6173 55.7558)')),
    ('Белорусская', ST_GeogFromText('SRID=4326;POINT(37.5880 55.7750)')),
    ('Киевская', ST_GeogFromText('SRID=4326;POINT(37.5614 55.7438)')),
    ('Лубянка', ST_GeogFromText('SRID=4326;POINT(37.6296 55.7626)')),
    ('Невский проспект', ST_GeogFromText('SRID=4326;POINT(30.3294 59.9343)')),
    ('Площадь Восстания', ST_GeogFromText('SRID=4326;POINT(30.3600 59.9340)')),
    ('Сенная площадь', ST_GeogFromText('SRID=4326;POINT(30.3150 59.9250)')),
    ('Казань Центральная', ST_GeogFromText('SRID=4326;POINT(49.1221 55.7887)')),
    ('Площадь Ленина', ST_GeogFromText('SRID=4326;POINT(49.1250 55.7910)'));

-- ============================================================
-- 5. Категории недвижимости
-- ============================================================
INSERT INTO property_categories (name) VALUES
    ('Квартиры'),
    ('Дома'),
    ('Апартаменты');

-- ============================================================
-- Застройщики (developers)
-- ============================================================
INSERT INTO developers (developer_name, avatar_url) VALUES
    ('ГК СтройГрупп Девелопмент', 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRw0YaHp30Xwau65hiGgdHeglHZVI9tFZDzoQ&s'),
    ('ПремиумДом Девелопмент', 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRw0YaHp30Xwau65hiGgdHeglHZVI9tFZDzoQ&s'),
    ('НордСтрой Девелопмент', 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRw0YaHp30Xwau65hiGgdHeglHZVI9tFZDzoQ&s'),
    ('КазаньИнвест Девелопмент', 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRw0YaHp30Xwau65hiGgdHeglHZVI9tFZDzoQ&s'),
    ('УралСтройКом Девелопмент', 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRw0YaHp30Xwau65hiGgdHeglHZVI9tFZDzoQ&s');

-- ============================================================
-- 6. ЖК Компании (utility_companies)
-- ============================================================
INSERT INTO utility_companies 
(company_name, phone, geo, address, avatar_url, alias, description, developer_id) 
VALUES
    ('СтройГрупп', '+7 495 123 48 77',
     ST_GeogFromText('SRID=4326;POINT(37.6173 55.7558)'),
     'г. Москва, ул. Тверская, д. 10, офис 5',
     'https://logotab.ru/storage/logotypes/1194/logotip-zhk-1083.jpg.jpg',
     'stroigroup',
     'Современный жилой комплекс бизнес-класса в центре Москвы с развитой инфраструктурой и подземным паркингом.',
     (SELECT id FROM developers WHERE developer_name = 'ГК СтройГрупп Девелопмент')),

    ('ПремиумДом', '+7 495 987 65 43',
     ST_GeogFromText('SRID=4326;POINT(37.5806 55.7495)'),
     'г. Москва, ул. Арбат, д. 20',
     'https://profi-storage.storage.yandexcloud.net/iblock/6c3/ymsd8l4okdnq64hrej0l6mnjkunh3ym4/logo-_11_.svg',
     'premiumdom',
     'Элитный жилой комплекс с дизайнерской архитектурой, закрытой территорией и круглосуточной охраной.',
     (SELECT id FROM developers WHERE developer_name = 'ПремиумДом Девелопмент')),

    ('НордСтрой', '+7 812 111 22 33',
     ST_GeogFromText('SRID=4326;POINT(30.3141 59.9311)'),
     'г. Санкт-Петербург, Невский пр., д. 50',
     'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRH_W7hzypC9OPgWn77SgqQ2OOZHYzGY1ZYXw&s',
     'nordstroy',
     'ЖК в историческом центре Санкт-Петербурга с видом на Неву и удобной транспортной доступностью.',
     (SELECT id FROM developers WHERE developer_name = 'НордСтрой Девелопмент')),

    ('КазаньИнвест', '+7 843 444 55 66',
     ST_GeogFromText('SRID=4326;POINT(49.1221 55.7887)'),
     'г. Казань, ул. Баумана, д. 15',
     'https://mir-s3-cdn-cf.behance.net/projects/404/171d7853318135.Y3JvcCwxMDIyLDgwMCwxODcsMA.jpg',
     'kazaninvest',
     'Комфортный жилой комплекс в центре Казани с благоустроенными дворами и развитой инфраструктурой.',
     (SELECT id FROM developers WHERE developer_name = 'КазаньИнвест Девелопмент')),

    ('УралСтройКом', '+7 343 777 88 99',
     ST_GeogFromText('SRID=4326;POINT(60.6122 56.8519)'),
     'г. Екатеринбург, ул. Ленина, д. 30',
     'https://sh.agency/upload/iblock/76b/76b329d4d06d8a87939c571a4601aa60.jpg',
     'uralstroy',
     'Современный ЖК в Екатеринбурге с просторными квартирами и удобным доступом к деловому центру города.',
     (SELECT id FROM developers WHERE developer_name = 'УралСтройКом Девелопмент'));



INSERT INTO utility_companies_photos (img_url, sequence_order, utility_company_id) VALUES
    ('https://dizayn-interera.moscow/images/blog/111/0_ta0g-5m.jpg', 1, (SELECT id FROM utility_companies WHERE alias = 'stroigroup')),
    ('https://salon.ru/storage/thumbs/gallery/272/271492/835_3500_s927.jpg', 2, (SELECT id FROM utility_companies WHERE alias = 'stroigroup')),
     ('https://n1s1.hsmedia.ru/0c/2e/40/0c2e4035e8da10aafba72e6f8b35b889/1000x750_0xac120003_8249795801571942265.jpg', 3, (SELECT id FROM utility_companies WHERE alias = 'stroigroup')),
    ('https://dizayn-interera.moscow/images/blog/111/0_ta0g-5m.jpg', 1, (SELECT id FROM utility_companies WHERE alias = 'premiumdom')),
    ('https://dizayn-interera.moscow/images/blog/111/0_ta0g-5m.jpg', 1, (SELECT id FROM utility_companies WHERE alias = 'nordstroy'));

-- ============================================================
-- 7. Дома (buildings)
-- ============================================================
INSERT INTO buildings (address, geo, city_id, metro_station_id, district, floor_count, company_id) VALUES
    -- Москва
    ('ул. Тверская, д. 25', ST_GeogFromText('SRID=4326;POINT(37.6155 55.7520)'), 1, 2, 'Центральный', 12, (SELECT id FROM utility_companies WHERE alias = 'stroigroup')),
    ('ул. Арбат, д. 36', ST_GeogFromText('SRID=4326;POINT(37.5870 55.7490)'), 1, 1, 'Арбат', 15, (SELECT id FROM utility_companies WHERE alias = 'premiumdom')),
    ('Смоленская пл., д. 3', ST_GeogFromText('SRID=4326;POINT(37.5990 55.7510)'), 1, 4, 'Арбат', 8, (SELECT id FROM utility_companies WHERE alias = 'premiumdom')),
    -- Санкт-Петербург
    ('Невский пр., д. 88', ST_GeogFromText('SRID=4326;POINT(30.3200 59.9340)'), 2, 7, 'Центральный', 18, (SELECT id FROM utility_companies WHERE alias = 'nordstroy')),
    ('ул. Садовая, д. 14', ST_GeogFromText('SRID=4326;POINT(30.2980 59.9280)'), 2, 8, 'Адмиралтейский', 10, (SELECT id FROM utility_companies WHERE alias = 'nordstroy')),
    -- Казань
    ('ул. Баумана, д. 42', ST_GeogFromText('SRID=4326;POINT(49.1250 55.7910)'), 3, 9, 'Вахитовский', 14, (SELECT id FROM utility_companies WHERE alias = 'kazaninvest'));

-- ============================================================
-- 8. Категории квартир (flat_categories)
-- ============================================================
INSERT INTO flat_categories (name) VALUES
    ('Студия'),
    ('1 комнатная'),
    ('2 комнатная'),
    ('3 комнатная'),
    ('Апартаменты'),
    ('Пентхаус');

-- ============================================================
-- 9. Объекты недвижимости (property)
-- ============================================================
INSERT INTO property (category_id, building_id, area) VALUES
    (1, 1, 20),  -- Тверская 25, студия
    (1, 1, 43),  -- Тверская 25, 2-к
    (1, 2, 93),  -- Арбат 36, пентхаус
    (1, 2, 34),  -- Арбат 36, 1-к
    (1, 3, 12),  -- Смоленская, студия
    (1, 3, 110), -- Смоленская, 2-к
    (1, 4, 223), -- Невский 88, 1-к
    (1, 4, 123), -- Невский 88, апартаменты
    (1, 5, 222), -- Садовая, 2-к
    (1, 6, 12);  -- Баумана, 3-к

-- ============================================================
-- 10. Квартиры (flat)
-- ============================================================
INSERT INTO flat (property_id, floor, number, category_id) VALUES
    (1, 3, 12, (SELECT id FROM flat_categories WHERE name = 'Студия')),
    (2, 7, 54, (SELECT id FROM flat_categories WHERE name = '2 комнатная')),
    (3, 10, 99, (SELECT id FROM flat_categories WHERE name = 'Пентхаус')),
    (4, 2, 5, (SELECT id FROM flat_categories WHERE name = '1 комнатная')),
    (5, 1, 2, (SELECT id FROM flat_categories WHERE name = 'Студия')),
    (6, 6, 45, (SELECT id FROM flat_categories WHERE name = '2 комнатная')),
    (7, 2, 8, (SELECT id FROM flat_categories WHERE name = '1 комнатная')),
    (8, 9, 77, (SELECT id FROM flat_categories WHERE name = 'Апартаменты')),
    (9, 3, 22, (SELECT id FROM flat_categories WHERE name = '2 комнатная')),
    (10, 4, 35, (SELECT id FROM flat_categories WHERE name = '3 комнатная'));

-- ============================================================
-- 11. Объявления (posters)
-- ============================================================
INSERT INTO posters (price, avatar_url, description, user_id, property_id, alias, created_at) VALUES
    (65000.00, 'https://dizayn-interera.moscow/images/blog/111/0_ta0g-5m.jpg', 'Уютная студия после капитального ремонта. Новая кухня, встроенные шкафы. Рядом метро.', 1, 1, 'studio-tverskaya', NOW() - INTERVAL '65 days'),
    (120000.00, 'https://salon.ru/storage/thumbs/gallery/272/271492/835_3500_s927.jpg', 'Светлая квартира с панорамным видом. Паркинг в подарок.', 1, 2, '2room-tverskaya', NOW() - INTERVAL '50 days'),
    (350000.00, 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTj_cN2apQuTB2vu5_v4J3FnyrhHD6Y5x_BXA&s', 'Уникальный пентхаус с открытой террасой 80 кв.м. Консьерж, закрытая территория.', 2, 3, 'penthouse-arbatskaya', NOW() - INTERVAL '95 days'),
    (85000.00, 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTxc4pnGQ858I3MeioxaDuJavns23B_bbJ_pw&s', 'Отличное место для жизни. Арбат в шаговой доступности. Мебель остаётся.', 2, 4, '1room-arbatskaya', NOW() - INTERVAL '35 days'),
    (55000.00, 'https://n1s1.hsmedia.ru/0c/2e/40/0c2e4035e8da10aafba72e6f8b35b889/1000x750_0xac120003_8249795801571942265.jpg', 'Компактная студия для одного или пары. Первый этаж, высокие потолки.', 3, 5, 'studio-smolenskaya', NOW() - INTERVAL '25 days'),
    (75000.00, 'https://inminecraft.ru/_ph/7/937019994.png', 'Квартира в историческом центре Петербурга. Дизайнерский ремонт 2024 года.', 4, 7, '1room-nevskiy', NOW() - INTERVAL '40 days'),
    (200000.00, 'https://st.dg-home.ru/upload/blog_editor/18b/2hwc7stcrv1vq9vx40o0k7bt3skuc1xn/11_divan.jpg', 'Элитные апартаменты, 9 этаж. Потрясающий вид на Неву. Подземный паркинг.', 4, 8, 'apartments-neva-view', NOW() - INTERVAL '60 days'),
    (90000.00, 'https://cs14.pikabu.ru/post_img/big/2024/01/12/11/1705087246125258124.jpg', 'Тихий двор, развитая инфраструктура. Школа и детсад в 5 минутах.', 5, 9, '2room-sadovaya', NOW() - INTERVAL '30 days'),
    (95000.00, 'https://garagetek.ru/uploads/images/GarageTek_Etush01.jpg', 'Большая семейная квартира. Все комнаты изолированы. Лоджия 8 кв.м.', 6, 10, '3room-kazan', NOW() - INTERVAL '35 days');

-- ============================================================
-- 12. Фотографии объявлений
-- ============================================================
-- Фотографии для объявлений
INSERT INTO poster_photos (img_url, sequence_order, poster_id) VALUES
    ('https://dizayn-interera.moscow/images/blog/111/0_ta0g-5m.jpg', 1, (SELECT id FROM posters WHERE alias = 'studio-tverskaya')),
    ('https://salon.ru/storage/thumbs/gallery/272/271492/835_3500_s927.jpg', 2, (SELECT id FROM posters WHERE alias = 'studio-tverskaya')),
    ('https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTj_cN2apQuTB2vu5_v4J3FnyrhHD6Y5x_BXA&s', 3, (SELECT id FROM posters WHERE alias = 'studio-tverskaya')),
    ('https://salon.ru/storage/thumbs/gallery/272/271492/835_3500_s927.jpg', 4, (SELECT id FROM posters WHERE alias = 'studio-tverskaya')),
    ('https://salon.ru/storage/thumbs/gallery/272/271492/835_3500_s927.jpg', 1, (SELECT id FROM posters WHERE alias = '2room-tverskaya')),
    ('https://dizayn-interera.moscow/images/blog/111/0_ta0g-5m.jpg', 2, (SELECT id FROM posters WHERE alias = '2room-tverskaya')),
    ('https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTj_cN2apQuTB2vu5_v4J3FnyrhHD6Y5x_BXA&s', 3, (SELECT id FROM posters WHERE alias = '2room-tverskaya')),
    ('https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTj_cN2apQuTB2vu5_v4J3FnyrhHD6Y5x_BXA&s', 1, (SELECT id FROM posters WHERE alias = 'penthouse-arbatskaya')),
    ('https://salon.ru/storage/thumbs/gallery/272/271492/835_3500_s927.jpg', 2, (SELECT id FROM posters WHERE alias = 'penthouse-arbatskaya')),
    ('https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTxc4pnGQ858I3MeioxaDuJavns23B_bbJ_pw&s', 1, (SELECT id FROM posters WHERE alias = '1room-arbatskaya')),
    ('https://n1s1.hsmedia.ru/0c/2e/40/0c2e4035e8da10aafba72e6f8b35b889/1000x750_0xac120003_8249795801571942265.jpg', 1, (SELECT id FROM posters WHERE alias = 'studio-smolenskaya')),
    ('https://salon.ru/storage/thumbs/gallery/272/271492/835_3500_s927.jpg', 2, (SELECT id FROM posters WHERE alias = 'studio-smolenskaya')),
    ('https://inminecraft.ru/_ph/7/937019994.png', 1, (SELECT id FROM posters WHERE alias = '1room-nevskiy')),
    ('https://salon.ru/storage/thumbs/gallery/272/271492/835_3500_s927.jpg', 2, (SELECT id FROM posters WHERE alias = '1room-nevskiy')),
    ('https://st.dg-home.ru/upload/blog_editor/18b/2hwc7stcrv1vq9vx40o0k7bt3skuc1xn/11_divan.jpg', 1, (SELECT id FROM posters WHERE alias = 'apartments-neva-view')),
    ('https://salon.ru/storage/thumbs/gallery/272/271492/835_3500_s927.jpg', 2, (SELECT id FROM posters WHERE alias = 'apartments-neva-view')),
    ('https://cs14.pikabu.ru/post_img/big/2024/01/12/11/1705087246125258124.jpg', 1, (SELECT id FROM posters WHERE alias = '2room-sadovaya')),
    ('https://garagetek.ru/uploads/images/GarageTek_Etush01.jpg', 1, (SELECT id FROM posters WHERE alias = '3room-kazan'));


-- ============================================================
-- 13. Лайки 
-- ============================================================
INSERT INTO likes (user_id, poster_id) VALUES
    (3, (SELECT id FROM posters WHERE alias = 'studio-tverskaya')),
    (5, (SELECT id FROM posters WHERE alias = 'studio-tverskaya')),
    (6, (SELECT id FROM posters WHERE alias = 'penthouse-arbatskaya')),
    (4, (SELECT id FROM posters WHERE alias = 'penthouse-arbatskaya')),
    (1, (SELECT id FROM posters WHERE alias = 'apartments-neva-view')),
    (2, (SELECT id FROM posters WHERE alias = 'apartments-neva-view'));
