-- ============================================================
-- 1. Города
-- ============================================================
INSERT INTO cities (city_name) VALUES
    ('Москва'),
    ('Санкт-Петербург'),
    ('Казань'),
    ('Новосибирск'),
    ('Екатеринбург');

-- ============================================================
-- 2. Станции метро
-- ============================================================
INSERT INTO metro_stations (station_name) VALUES
    ('Арбатская'),
    ('Тверская'),
    ('Белорусская'),
    ('Киевская'),
    ('Лубянка'),
    ('Невский проспект'),
    ('Площадь Восстания'),
    ('Сенная площадь'),
    ('Казань Центральная'),
    ('Площадь Ленина');

-- ============================================================
-- 3. ЖК Компании
-- ============================================================
INSERT INTO utility_companies (company_name, contacts, geo, address, city_id, metro_station_id) VALUES
    (
        'СтройГрупп',
        '+7 (495) 123-45-67',
        ST_GeogFromText('POINT(37.6173 55.7558)'),
        'ул. Тверская, д. 10, офис 5',
        (SELECT id FROM cities WHERE city_name = 'Москва'),
        (SELECT id FROM metro_stations WHERE station_name = 'Тверская')
    ),
    (
        'ПремиумДом',
        '+7 (495) 987-65-43',
        ST_GeogFromText('POINT(37.5806 55.7495)'),
        'ул. Арбат, д. 20',
        (SELECT id FROM cities WHERE city_name = 'Москва'),
        (SELECT id FROM metro_stations WHERE station_name = 'Арбатская')
    ),
    (
        'НордСтрой',
        '+7 (812) 111-22-33',
        ST_GeogFromText('POINT(30.3141 59.9311)'),
        'Невский пр., д. 50',
        (SELECT id FROM cities WHERE city_name = 'Санкт-Петербург'),
        (SELECT id FROM metro_stations WHERE station_name = 'Невский проспект')
    ),
    (
        'КазаньИнвест',
        '+7 (843) 444-55-66',
        ST_GeogFromText('POINT(49.1221 55.7887)'),
        'ул. Баумана, д. 15',
        (SELECT id FROM cities WHERE city_name = 'Казань'),
        (SELECT id FROM metro_stations WHERE station_name = 'Казань Центральная')
    ),
    (
        'УралСтройКом',
        '+7 (343) 777-88-99',
        ST_GeogFromText('POINT(60.6122 56.8519)'),
        'ул. Ленина, д. 30',
        (SELECT id FROM cities WHERE city_name = 'Екатеринбург'),
        NULL
    );

-- ============================================================
-- 4. Пользователи
-- ============================================================
INSERT INTO users (email, hashed_password, provider, salt, company_id) VALUES
    (
        'ivan.petrov@mail.ru',
        '$2b$12$abcdefgh1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMN',
        NULL,
        'salt_ivan_001',
        (SELECT id FROM utility_companies WHERE company_name = 'СтройГрупп')
    ),
    (
        'anna.sokolova@yandex.ru',
        '$2b$12$zyxwvutsrqponmlkjihgfedcbaZYXWVUTSRQPONMLKJIHGFEDCBA0123',
        NULL,
        'salt_anna_002',
        (SELECT id FROM utility_companies WHERE company_name = 'ПремиумДом')
    ),
    (
        'sergey.kuzmin@gmail.com',
        NULL,
        'google',
        'salt_sergey_003',
        NULL
    ),
    (
        'olga.morozova@mail.ru',
        '$2b$12$1234567890abcdefgh1234567890abcdefghijklmnopqrstuvwxyzABC',
        NULL,
        'salt_olga_004',
        (SELECT id FROM utility_companies WHERE company_name = 'НордСтрой')
    ),
    (
        'dmitry.volkov@yandex.ru',
        NULL,
        'vk',
        'salt_dmitry_005',
        NULL
    ),
    (
        'maria.novikova@gmail.com',
        '$2b$12$qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM12',
        NULL,
        'salt_maria_006',
        (SELECT id FROM utility_companies WHERE company_name = 'КазаньИнвест')
    );

-- ============================================================
-- 5. Дома
-- ============================================================
INSERT INTO buildings (geo, address, district, company_id, city_id, metro_station_id) VALUES
    (
        ST_GeogFromText('POINT(37.6155 55.7520)'),
        'ул. Тверская, д. 25',
        'Центральный',
        (SELECT id FROM utility_companies WHERE company_name = 'СтройГрупп'),
        (SELECT id FROM cities WHERE city_name = 'Москва'),
        (SELECT id FROM metro_stations WHERE station_name = 'Тверская')
    ),
    (
        ST_GeogFromText('POINT(37.5870 55.7490)'),
        'ул. Арбат, д. 36',
        'Арбат',
        (SELECT id FROM utility_companies WHERE company_name = 'ПремиумДом'),
        (SELECT id FROM cities WHERE city_name = 'Москва'),
        (SELECT id FROM metro_stations WHERE station_name = 'Арбатская')
    ),
    (
        ST_GeogFromText('POINT(37.5990 55.7510)'),
        'Смоленская пл., д. 3',
        'Арбат',
        (SELECT id FROM utility_companies WHERE company_name = 'ПремиумДом'),
        (SELECT id FROM cities WHERE city_name = 'Москва'),
        (SELECT id FROM metro_stations WHERE station_name = 'Киевская')
    ),
    (
        ST_GeogFromText('POINT(30.3200 59.9340)'),
        'Невский пр., д. 88',
        'Центральный',
        (SELECT id FROM utility_companies WHERE company_name = 'НордСтрой'),
        (SELECT id FROM cities WHERE city_name = 'Санкт-Петербург'),
        (SELECT id FROM metro_stations WHERE station_name = 'Площадь Восстания')
    ),
    (
        ST_GeogFromText('POINT(30.2980 59.9280)'),
        'ул. Садовая, д. 14',
        'Адмиралтейский',
        (SELECT id FROM utility_companies WHERE company_name = 'НордСтрой'),
        (SELECT id FROM cities WHERE city_name = 'Санкт-Петербург'),
        (SELECT id FROM metro_stations WHERE station_name = 'Сенная площадь')
    ),
    (
        ST_GeogFromText('POINT(49.1250 55.7910)'),
        'ул. Баумана, д. 42',
        'Вахитовский',
        (SELECT id FROM utility_companies WHERE company_name = 'КазаньИнвест'),
        (SELECT id FROM cities WHERE city_name = 'Казань'),
        NULL
    );

-- ============================================================
-- 6. Типы помещений
-- ============================================================
INSERT INTO apartment_categories (name, description) VALUES
    ('Студия',          'Однокомнатная квартира без отдельной спальни'),
    ('1-комнатная',     'Квартира с одной жилой комнатой'),
    ('2-комнатная',     'Квартира с двумя жилыми комнатами'),
    ('3-комнатная',     'Квартира с тремя жилыми комнатами'),
    ('Апартаменты',     'Нежилое помещение с улучшенной отделкой'),
    ('Пентхаус',        'Квартира на последнем этаже с террасой');

-- ============================================================
-- 7. Помещения
-- ============================================================
INSERT INTO apartments (floor, number, building_id, category_id) VALUES
    -- ул. Тверская, д. 25
    (3,  12, (SELECT id FROM buildings WHERE address = 'ул. Тверская, д. 25'),   (SELECT id FROM apartment_categories WHERE name = 'Студия')),
    (5,  31, (SELECT id FROM buildings WHERE address = 'ул. Тверская, д. 25'),   (SELECT id FROM apartment_categories WHERE name = '1-комнатная')),
    (7,  54, (SELECT id FROM buildings WHERE address = 'ул. Тверская, д. 25'),   (SELECT id FROM apartment_categories WHERE name = '2-комнатная')),
    -- ул. Арбат, д. 36
    (2,   5, (SELECT id FROM buildings WHERE address = 'ул. Арбат, д. 36'),      (SELECT id FROM apartment_categories WHERE name = '1-комнатная')),
    (4,  18, (SELECT id FROM buildings WHERE address = 'ул. Арбат, д. 36'),      (SELECT id FROM apartment_categories WHERE name = '3-комнатная')),
    (10, 99, (SELECT id FROM buildings WHERE address = 'ул. Арбат, д. 36'),      (SELECT id FROM apartment_categories WHERE name = 'Пентхаус')),
    -- Смоленская пл., д. 3
    (1,   2, (SELECT id FROM buildings WHERE address = 'Смоленская пл., д. 3'),  (SELECT id FROM apartment_categories WHERE name = 'Студия')),
    (6,  45, (SELECT id FROM buildings WHERE address = 'Смоленская пл., д. 3'),  (SELECT id FROM apartment_categories WHERE name = '2-комнатная')),
    -- Невский пр., д. 88
    (2,   8, (SELECT id FROM buildings WHERE address = 'Невский пр., д. 88'),    (SELECT id FROM apartment_categories WHERE name = '1-комнатная')),
    (9,  77, (SELECT id FROM buildings WHERE address = 'Невский пр., д. 88'),    (SELECT id FROM apartment_categories WHERE name = 'Апартаменты')),
    -- ул. Садовая, д. 14
    (3,  22, (SELECT id FROM buildings WHERE address = 'ул. Садовая, д. 14'),    (SELECT id FROM apartment_categories WHERE name = '2-комнатная')),
    -- ул. Баумана, д. 42
    (4,  35, (SELECT id FROM buildings WHERE address = 'ул. Баумана, д. 42'),    (SELECT id FROM apartment_categories WHERE name = '3-комнатная'));

-- ============================================================
-- 8. Объявления
-- ============================================================
INSERT INTO posters (title, price, avatar_url, description, user_id, apartment_id) VALUES
    (
        'Студия у Тверской, свежий ремонт',
        65000.00,
        'https://cdn.example.com/posters/1/avatar.jpg',
        'Уютная студия после капитального ремонта. Новая кухня, встроенные шкафы. Рядом метро.',
        (SELECT id FROM users WHERE email = 'ivan.petrov@mail.ru'),
        (SELECT id FROM apartments WHERE number = 12 AND building_id = (SELECT id FROM buildings WHERE address = 'ул. Тверская, д. 25'))
    ),
    (
        'Просторная 2-комнатная на Тверской',
        120000.00,
        'https://cdn.example.com/posters/2/avatar.jpg',
        'Светлая квартира с панорамным видом. Паркинг в подарок.',
        (SELECT id FROM users WHERE email = 'ivan.petrov@mail.ru'),
        (SELECT id FROM apartments WHERE number = 54 AND building_id = (SELECT id FROM buildings WHERE address = 'ул. Тверская, д. 25'))
    ),
    (
        'Пентхаус на Арбате — эксклюзив',
        350000.00,
        'https://cdn.example.com/posters/3/avatar.jpg',
        'Уникальный пентхаус с открытой террасой 80 кв.м. Консьерж, закрытая территория.',
        (SELECT id FROM users WHERE email = 'anna.sokolova@yandex.ru'),
        (SELECT id FROM apartments WHERE number = 99 AND building_id = (SELECT id FROM buildings WHERE address = 'ул. Арбат, д. 36'))
    ),
    (
        '1-комнатная на Арбате',
        85000.00,
        'https://cdn.example.com/posters/4/avatar.jpg',
        'Отличное место для жизни. Арбат в шаговой доступности. Мебель остаётся.',
        (SELECT id FROM users WHERE email = 'anna.sokolova@yandex.ru'),
        (SELECT id FROM apartments WHERE number = 5 AND building_id = (SELECT id FROM buildings WHERE address = 'ул. Арбат, д. 36'))
    ),
    (
        'Студия у Смоленской площади',
        55000.00,
        'https://cdn.example.com/posters/5/avatar.jpg',
        'Компактная студия для одного или пары. Первый этаж, высокие потолки.',
        (SELECT id FROM users WHERE email = 'sergey.kuzmin@gmail.com'),
        (SELECT id FROM apartments WHERE number = 2 AND building_id = (SELECT id FROM buildings WHERE address = 'Смоленская пл., д. 3'))
    ),
    (
        '1-комнатная на Невском',
        75000.00,
        'https://cdn.example.com/posters/6/avatar.jpg',
        'Квартира в историческом центре Петербурга. Дизайнерский ремонт 2024 года.',
        (SELECT id FROM users WHERE email = 'olga.morozova@mail.ru'),
        (SELECT id FROM apartments WHERE number = 8 AND building_id = (SELECT id FROM buildings WHERE address = 'Невский пр., д. 88'))
    ),
    (
        'Апартаменты с видом на Неву',
        200000.00,
        'https://cdn.example.com/posters/7/avatar.jpg',
        'Элитные апартаменты, 9 этаж. Потрясающий вид на Неву. Подземный паркинг.',
        (SELECT id FROM users WHERE email = 'olga.morozova@mail.ru'),
        (SELECT id FROM apartments WHERE number = 77 AND building_id = (SELECT id FROM buildings WHERE address = 'Невский пр., д. 88'))
    ),
    (
        '2-комнатная на Садовой',
        90000.00,
        'https://cdn.example.com/posters/8/avatar.jpg',
        'Тихий двор, развитая инфраструктура. Школа и детсад в 5 минутах.',
        (SELECT id FROM users WHERE email = 'dmitry.volkov@yandex.ru'),
        (SELECT id FROM apartments WHERE number = 22 AND building_id = (SELECT id FROM buildings WHERE address = 'ул. Садовая, д. 14'))
    ),
    (
        '3-комнатная в центре Казани',
        95000.00,
        'https://cdn.example.com/posters/9/avatar.jpg',
        'Большая семейная квартира. Все комнаты изолированы. Лоджия 8 кв.м.',
        (SELECT id FROM users WHERE email = 'maria.novikova@gmail.com'),
        (SELECT id FROM apartments WHERE number = 35 AND building_id = (SELECT id FROM buildings WHERE address = 'ул. Баумана, д. 42'))
    );

-- ============================================================
-- 9. История изменения цены
-- ============================================================
INSERT INTO price_history (price, changed_at, poster_id) VALUES
    (70000.00, NOW() - INTERVAL '60 days', (SELECT id FROM posters WHERE title = 'Студия у Тверской, свежий ремонт')),
    (67000.00, NOW() - INTERVAL '30 days', (SELECT id FROM posters WHERE title = 'Студия у Тверской, свежий ремонт')),
    (65000.00, NOW() - INTERVAL '5 days',  (SELECT id FROM posters WHERE title = 'Студия у Тверской, свежий ремонт')),

    (130000.00, NOW() - INTERVAL '45 days', (SELECT id FROM posters WHERE title = 'Просторная 2-комнатная на Тверской')),
    (120000.00, NOW() - INTERVAL '10 days', (SELECT id FROM posters WHERE title = 'Просторная 2-комнатная на Тверской')),

    (400000.00, NOW() - INTERVAL '90 days', (SELECT id FROM posters WHERE title = 'Пентхаус на Арбате — эксклюзив')),
    (375000.00, NOW() - INTERVAL '60 days', (SELECT id FROM posters WHERE title = 'Пентхаус на Арбате — эксклюзив')),
    (350000.00, NOW() - INTERVAL '20 days', (SELECT id FROM posters WHERE title = 'Пентхаус на Арбате — эксклюзив')),

    (80000.00, NOW() - INTERVAL '30 days', (SELECT id FROM posters WHERE title = '1-комнатная на Арбате')),
    (85000.00, NOW() - INTERVAL '7 days',  (SELECT id FROM posters WHERE title = '1-комнатная на Арбате')),

    (210000.00, NOW() - INTERVAL '50 days', (SELECT id FROM posters WHERE title = 'Апартаменты с видом на Неву')),
    (200000.00, NOW() - INTERVAL '15 days', (SELECT id FROM posters WHERE title = 'Апартаменты с видом на Неву'));

-- ============================================================
-- 10. Лайки
-- ============================================================
INSERT INTO likes (user_id, poster_id) VALUES
    ((SELECT id FROM users WHERE email = 'sergey.kuzmin@gmail.com'),   (SELECT id FROM posters WHERE title = 'Студия у Тверской, свежий ремонт')),
    ((SELECT id FROM users WHERE email = 'dmitry.volkov@yandex.ru'),   (SELECT id FROM posters WHERE title = 'Студия у Тверской, свежий ремонт')),
    ((SELECT id FROM users WHERE email = 'maria.novikova@gmail.com'),  (SELECT id FROM posters WHERE title = 'Пентхаус на Арбате — эксклюзив')),
    ((SELECT id FROM users WHERE email = 'olga.morozova@mail.ru'),     (SELECT id FROM posters WHERE title = 'Пентхаус на Арбате — эксклюзив')),
    ((SELECT id FROM users WHERE email = 'ivan.petrov@mail.ru'),       (SELECT id FROM posters WHERE title = 'Апартаменты с видом на Неву')),
    ((SELECT id FROM users WHERE email = 'anna.sokolova@yandex.ru'),   (SELECT id FROM posters WHERE title = 'Апартаменты с видом на Неву')),
    ((SELECT id FROM users WHERE email = 'sergey.kuzmin@gmail.com'),   (SELECT id FROM posters WHERE title = '1-комнатная на Невском')),
    ((SELECT id FROM users WHERE email = 'dmitry.volkov@yandex.ru'),   (SELECT id FROM posters WHERE title = '3-комнатная в центре Казани'));

-- ============================================================
-- 11. Просмотры
-- ============================================================
INSERT INTO views (user_id, poster_id) VALUES
    ((SELECT id FROM users WHERE email = 'sergey.kuzmin@gmail.com'),   (SELECT id FROM posters WHERE title = 'Студия у Тверской, свежий ремонт')),
    ((SELECT id FROM users WHERE email = 'sergey.kuzmin@gmail.com'),   (SELECT id FROM posters WHERE title = 'Пентхаус на Арбате — эксклюзив')),
    ((SELECT id FROM users WHERE email = 'dmitry.volkov@yandex.ru'),   (SELECT id FROM posters WHERE title = 'Студия у Тверской, свежий ремонт')),
    ((SELECT id FROM users WHERE email = 'dmitry.volkov@yandex.ru'),   (SELECT id FROM posters WHERE title = 'Просторная 2-комнатная на Тверской')),
    ((SELECT id FROM users WHERE email = 'dmitry.volkov@yandex.ru'),   (SELECT id FROM posters WHERE title = '1-комнатная на Арбате')),
    ((SELECT id FROM users WHERE email = 'maria.novikova@gmail.com'),  (SELECT id FROM posters WHERE title = 'Пентхаус на Арбате — эксклюзив')),
    ((SELECT id FROM users WHERE email = 'maria.novikova@gmail.com'),  (SELECT id FROM posters WHERE title = 'Апартаменты с видом на Неву')),
    ((SELECT id FROM users WHERE email = 'olga.morozova@mail.ru'),     (SELECT id FROM posters WHERE title = 'Пентхаус на Арбате — эксклюзив')),
    ((SELECT id FROM users WHERE email = 'ivan.petrov@mail.ru'),       (SELECT id FROM posters WHERE title = 'Апартаменты с видом на Неву')),
    ((SELECT id FROM users WHERE email = 'anna.sokolova@yandex.ru'),   (SELECT id FROM posters WHERE title = 'Апартаменты с видом на Неву')),
    ((SELECT id FROM users WHERE email = 'anna.sokolova@yandex.ru'),   (SELECT id FROM posters WHERE title = '3-комнатная в центре Казани'));

-- ============================================================
-- 12. Фотографии объявлений
-- ============================================================
INSERT INTO poster_photos (img_url, sequence_order, poster_id) VALUES
    ('https://cdn.example.com/posters/1/photo1.jpg', 1, (SELECT id FROM posters WHERE title = 'Студия у Тверской, свежий ремонт')),
    ('https://cdn.example.com/posters/1/photo2.jpg', 2, (SELECT id FROM posters WHERE title = 'Студия у Тверской, свежий ремонт')),
    ('https://cdn.example.com/posters/1/photo3.jpg', 3, (SELECT id FROM posters WHERE title = 'Студия у Тверской, свежий ремонт')),

    ('https://cdn.example.com/posters/2/photo1.jpg', 1, (SELECT id FROM posters WHERE title = 'Просторная 2-комнатная на Тверской')),
    ('https://cdn.example.com/posters/2/photo2.jpg', 2, (SELECT id FROM posters WHERE title = 'Просторная 2-комнатная на Тверской')),

    ('https://cdn.example.com/posters/3/photo1.jpg', 1, (SELECT id FROM posters WHERE title = 'Пентхаус на Арбате — эксклюзив')),
    ('https://cdn.example.com/posters/3/photo2.jpg', 2, (SELECT id FROM posters WHERE title = 'Пентхаус на Арбате — эксклюзив')),
    ('https://cdn.example.com/posters/3/photo3.jpg', 3, (SELECT id FROM posters WHERE title = 'Пентхаус на Арбате — эксклюзив')),
    ('https://cdn.example.com/posters/3/photo4.jpg', 4, (SELECT id FROM posters WHERE title = 'Пентхаус на Арбате — эксклюзив')),

    ('https://cdn.example.com/posters/4/photo1.jpg', 1, (SELECT id FROM posters WHERE title = '1-комнатная на Арбате')),
    ('https://cdn.example.com/posters/4/photo2.jpg', 2, (SELECT id FROM posters WHERE title = '1-комнатная на Арбате')),

    ('https://cdn.example.com/posters/5/photo1.jpg', 1, (SELECT id FROM posters WHERE title = 'Студия у Смоленской площади')),

    ('https://cdn.example.com/posters/6/photo1.jpg', 1, (SELECT id FROM posters WHERE title = '1-комнатная на Невском')),
    ('https://cdn.example.com/posters/6/photo2.jpg', 2, (SELECT id FROM posters WHERE title = '1-комнатная на Невском')),

    ('https://cdn.example.com/posters/7/photo1.jpg', 1, (SELECT id FROM posters WHERE title = 'Апартаменты с видом на Неву')),
    ('https://cdn.example.com/posters/7/photo2.jpg', 2, (SELECT id FROM posters WHERE title = 'Апартаменты с видом на Неву')),
    ('https://cdn.example.com/posters/7/photo3.jpg', 3, (SELECT id FROM posters WHERE title = 'Апартаменты с видом на Неву')),

    ('https://cdn.example.com/posters/8/photo1.jpg', 1, (SELECT id FROM posters WHERE title = '2-комнатная на Садовой')),
    ('https://cdn.example.com/posters/8/photo2.jpg', 2, (SELECT id FROM posters WHERE title = '2-комнатная на Садовой')),

    ('https://cdn.example.com/posters/9/photo1.jpg', 1, (SELECT id FROM posters WHERE title = '3-комнатная в центре Казани')),
    ('https://cdn.example.com/posters/9/photo2.jpg', 2, (SELECT id FROM posters WHERE title = '3-комнатная в центре Казани')),
    ('https://cdn.example.com/posters/9/photo3.jpg', 3, (SELECT id FROM posters WHERE title = '3-комнатная в центре Казани'));
