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
INSERT INTO property_categories (name, alias) VALUES
    ('Квартира', 'flat'),
    ('Дом', 'house'),
    ('Апартаменты', 'apartments');

-- ============================================================
-- Застройщики (developers)
-- ============================================================
INSERT INTO developers (developer_name, avatar_url) VALUES
    ('Донстрой', 'https://upload.wikimedia.org/wikipedia/commons/1/1e/%D0%9B%D0%BE%D0%B3%D0%BE%D1%82%D0%B8%D0%BF_%D0%94%D0%BE%D0%BD%D1%81%D1%82%D1%80%D0%BE%D0%B9.png'),
    ('ПремиумДом Девелопмент', 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTHoLsgU0o4GhmieKk6j0Bq2gILNxDrQLALHg&s'),
    ('НордСтрой Девелопмент', 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcStmh7s8ef89_ITaTytzK4mP9O8ReVnrh-3UQ&s'),
    ('КазаньИнвест Девелопмент', 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRw0YaHp30Xwau65hiGgdHeglHZVI9tFZDzoQ&s'),
    ('УралСтройКом Девелопмент', 'https://domselect.ru/storage/main/sk-psk.jpg');

-- ============================================================
-- 6. ЖК Компании (utility_companies)
-- ============================================================
INSERT INTO utility_companies 
(company_name, phone, geo, address, avatar_url, alias, description, developer_id) 
VALUES
    ('Символ', '+7 495 123 48 77',
     ST_GeogFromText('SRID=4326;POINT(37.6173 55.7558)'),
     'г. Москва, ул. Тверская, д. 10, офис 5',
     'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTr_mhDYpgBiUdG8GfGpnse45MCmYiSZIAu9w&s',
     'stroigroup',
     'СИМВОЛ — новая городская среда в окружении исторического центра Москвы. Высокий статус СИМВОЛА подтверждают десятки профессиональных наград, в числе которых лучший квартал Москвы и России, лучший городской дизайн и самая комфортная среда.
Ультрасовременная архитектура от звездных бюро Москвы и Лондона и 10 га парковой территории сделали СИМВОЛ одним из самых престижных и узнаваемых жилых пространств столицы. Парк «Зеленая река» сегодня — не просто одно из любимых мест отдыха и прогулок москвичей, но и настоящая визитная карточка СИМВОЛА.
Здесь все готово для жизни. Продуманные общественные пространства в домах, более 100 действующих объектов инфраструктуры и все возможности центра с его галереями, музеями, театрами, торговыми и бизнес-центрами превратили СИМВОЛ в особый мир, где реализуются любые сценарии вашей жизни.',
     (SELECT id FROM developers WHERE developer_name = 'Донстрой')),

    ('ПремиумДом', '+7 495 987 65 43',
     ST_GeogFromText('SRID=4326;POINT(37.5806 55.7495)'),
     'г. Москва, ул. Арбат, д. 20',
     'https://profi-storage.storage.yandexcloud.net/iblock/6c3/ymsd8l4okdnq64hrej0l6mnjkunh3ym4/logo-_11_.svg',
     'premiumdom',
     'Элитный жилой комплекс с уникальной архитектурной концепцией, разработанной ведущими бюро, где каждая деталь продумана до мелочей. Проект сочетает современные технологии строительства, премиальные материалы и изысканный дизайн общественных пространств. Закрытая охраняемая территория обеспечивает высокий уровень безопасности и приватности для жителей, а система видеонаблюдения и контроль доступа работают круглосуточно. Внутренний двор благоустроен по принципу «двор без машин»: ландшафтный дизайн, зоны отдыха, детские и спортивные площадки создают комфортную среду для жизни. В шаговой доступности находятся рестораны, бутики, образовательные учреждения и культурные объекты центра Москвы. Подземный паркинг, консьерж-сервис и развитая инфраструктура делают комплекс идеальным выбором для тех, кто ценит статус, комфорт и высокий уровень сервиса.',
     (SELECT id FROM developers WHERE developer_name = 'ПремиумДом Девелопмент')),

    ('НордСтрой', '+7 812 111 22 33',
     ST_GeogFromText('SRID=4326;POINT(30.3141 59.9311)'),
     'г. Санкт-Петербург, Невский пр., д. 50',
     'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRH_W7hzypC9OPgWn77SgqQ2OOZHYzGY1ZYXw&s',
     'nordstroy',
     'Современный жилой комплекс, расположенный в историческом центре Санкт-Петербурга, гармонично вписанный в архитектурный облик города. Из окон открываются живописные виды на набережные и акваторию Невы, создавая атмосферу уюта и вдохновения. Проект предусматривает удобные планировки квартир с высокими потолками и большими окнами, обеспечивающими максимальное естественное освещение. Внутренняя инфраструктура включает коммерческие помещения, зоны отдыха и благоустроенные дворы с озеленением. Отличная транспортная доступность позволяет быстро добраться до ключевых районов города, а близость к станциям метро делает передвижение особенно удобным. В пешей доступности находятся культурные достопримечательности, театры, музеи, а также лучшие рестораны и кафе города. Комплекс идеально подходит для тех, кто хочет жить в центре событий, не отказываясь от комфорта и современного уровня жизни.',
     (SELECT id FROM developers WHERE developer_name = 'НордСтрой Девелопмент')),

    ('КазаньИнвест', '+7 843 444 55 66',
     ST_GeogFromText('SRID=4326;POINT(49.1221 55.7887)'),
     'г. Казань, ул. Баумана, д. 15',
     'https://mir-s3-cdn-cf.behance.net/projects/404/171d7853318135.Y3JvcCwxMDIyLDgwMCwxODcsMA.jpg',
     'kazaninvest',
     'Жилой комплекс комфорт-класса в самом сердце Казани, сочетающий современную архитектуру и продуманную городскую среду. Проект ориентирован на удобство повседневной жизни: эргономичные планировки квартир, качественные строительные материалы и современные инженерные решения обеспечивают высокий уровень комфорта. Особое внимание уделено благоустройству территории — зеленые дворы, прогулочные зоны, детские и спортивные площадки формируют безопасное и уютное пространство для жителей всех возрастов. Комплекс расположен в районе с развитой инфраструктурой: рядом находятся торговые центры, школы, детские сады, медицинские учреждения и остановки общественного транспорта. Исторический центр города и основные достопримечательности находятся в шаговой доступности, что делает локацию особенно привлекательной. Это идеальный вариант для тех, кто ищет баланс между динамичной городской жизнью и комфортом.',
     (SELECT id FROM developers WHERE developer_name = 'КазаньИнвест Девелопмент')),

    ('УралСтройКом', '+7 343 777 88 99',
     ST_GeogFromText('SRID=4326;POINT(60.6122 56.8519)'),
     'г. Екатеринбург, ул. Ленина, д. 30',
     'https://sh.agency/upload/iblock/76b/76b329d4d06d8a87939c571a4601aa60.jpg',
     'uralstroy',
     'Современный жилой комплекс в Екатеринбурге, созданный с учетом актуальных требований к качеству жизни в мегаполисе. Архитектурная концепция проекта сочетает лаконичный стиль и функциональность, а разнообразие планировок позволяет подобрать оптимальное решение для любого образа жизни. Просторные квартиры с панорамными окнами наполняются естественным светом, создавая ощущение открытого пространства. Комплекс располагает собственной инфраструктурой: подземный паркинг, коммерческие помещения, зоны отдыха и благоустроенные дворы с озеленением. Удобное расположение обеспечивает быстрый доступ к деловому центру города, а также к основным транспортным магистралям. В непосредственной близости находятся бизнес-центры, учебные заведения, магазины и рестораны. Проект ориентирован на активных городских жителей, ценящих комфорт, практичность и современную городскую среду.',
     (SELECT id FROM developers WHERE developer_name = 'УралСтройКом Девелопмент'));



INSERT INTO utility_companies_photos (img_url, sequence_order, utility_company_id) VALUES
    ('https://simvol-kvartal.ru/upload/dev2fun.imagecompress/webp/iblock/00c/y13z8xnvwh3kfx8uzjew62t43ux519b6.webp', 1, (SELECT id FROM utility_companies WHERE alias = 'stroigroup')),
    ('https://simvol-kvartal.ru/upload/dev2fun.imagecompress/webp/iblock/465/ef7v1x7whetdlyx82k02i1k3jovx3zli.webp', 2, (SELECT id FROM utility_companies WHERE alias = 'stroigroup')),
    ('https://simvol-kvartal.ru/upload/dev2fun.imagecompress/webp/iblock/235/ns84d8jqjfximts8vcupphhbbsglx1zs.webp', 3, (SELECT id FROM utility_companies WHERE alias = 'stroigroup')),
    ('https://cdn.a101.ru/proxy/insecure/w:2560/q:80/plain/https://cdn.a101.ru/mmedia/p/pag/i/25cd7f9aff.jpg@webp', 1, (SELECT id FROM utility_companies WHERE alias = 'premiumdom')),
    ('https://cdn.a101.ru/proxy/insecure/w:2560/q:80/plain/https://cdn.a101.ru/mmedia/p/pag/i/3d6f8f1613.jpg@webp', 2, (SELECT id FROM utility_companies WHERE alias = 'premiumdom')),
    ('https://cdn.a101.ru/proxy/insecure/w:2560/q:80/plain/https://cdn.a101.ru/mmedia/p/pag/i/8f84c7b936.jpg@webp', 3, (SELECT id FROM utility_companies WHERE alias = 'premiumdom')),
    ('https://cdn.a101.ru/proxy/insecure/w:2560/q:80/plain/https://cdn.a101.ru/mmedia/p/pag/i/584a3cf6d9.jpg@webp', 1, (SELECT id FROM utility_companies WHERE alias = 'nordstroy')),
    ('https://ss.metronews.ru/userfiles/materials/181/1819761/858x540_d3c9c6c9.jpg', 2, (SELECT id FROM utility_companies WHERE alias = 'nordstroy')),
    ('https://cdn.a101.ru/proxy/insecure/w:2560/q:80/plain/https://cdn.a101.ru/mmedia/p/pag/i/b3ae409b35.jpg@webp', 3, (SELECT id FROM utility_companies WHERE alias = 'nordstroy')),
    ('https://donstroy.moscow/upload/iblock/958/iwi8wynvfpd9pfq48mcn9gma2bkfnlws.jpg', 1, (SELECT id FROM utility_companies WHERE alias = 'kazaninvest')),
    ('https://donstroy.moscow/upload/iblock/5a5/a5edkxw0y093s6ulkr7ono86tei7m1da.jpg', 2, (SELECT id FROM utility_companies WHERE alias = 'kazaninvest')),
    ('https://donstroy.moscow/upload/iblock/049/j9tydlo0tmyvfmj21jq573msglaqbvj8.jpg', 3, (SELECT id FROM utility_companies WHERE alias = 'kazaninvest')),
    ('https://cdn.samolet.ru/imgproxy/insecure/q:90/rs:fill:1260:695/g:ce/bl:0/c:0/plain/https://media.samolet.ru/r/pp/pptgc/image/1_B4fcK8e.jpg@webp', 1, (SELECT id FROM utility_companies WHERE alias = 'uralstroy')),
    ('https://cdn.samolet.ru/imgproxy/insecure/q:90/rs:fill:1260:695/g:ce/bl:0/c:0/plain/https://media.samolet.ru/r/pp/pptgc/image/3_UQDjomO.jpg@webp', 2, (SELECT id FROM utility_companies WHERE alias = 'uralstroy')),
    ('https://cdn.samolet.ru/imgproxy/insecure/q:90/rs:fill:1260:695/g:ce/bl:0/c:0/plain/https://media.samolet.ru/r/pp/pptgc/image/%D0%9A%D0%B2%D0%B0%D1%80%D1%82%D0%B0%D0%BB_%D0%BD%D0%B0_%D0%92%D0%BE%D0%B4%D0%B5_5_%D0%BA%D0%BE%D1%80%D0%BF%D1%83%D1%81_%D0%B8%D0%BC%D0%B8%D0%B4%D0%B6_2_2025-03-06.jpg@webp', 3, (SELECT id FROM utility_companies WHERE alias = 'uralstroy'));


-- ============================================================
-- 7. Дома (buildings)
-- ============================================================
INSERT INTO buildings (address, geo, city_id, metro_station_id, district, floor_count, company_id) VALUES
    -- Москва
    ('ул. Тверская, д. 25', ST_GeogFromText('SRID=4326;POINT(37.6155 55.7520)'), 1, 2, 'Центральный', 12, (SELECT id FROM utility_companies WHERE alias = 'stroigroup')),
    ('ул. Тверская, д. 30', ST_GeogFromText('SRID=4326;POINT(37.6155 55.7520)'), 1, 2, 'Центральный', 12, (SELECT id FROM utility_companies WHERE alias = 'stroigroup')),
    ('ул. Арбат, д. 36', ST_GeogFromText('SRID=4326;POINT(37.5870 55.7490)'), 1, 1, 'Арбат', 15, (SELECT id FROM utility_companies WHERE alias = 'premiumdom')),
    ('Смоленская пл., д. 3', ST_GeogFromText('SRID=4326;POINT(37.5990 55.7510)'), 1, 4, 'Арбат', 8, (SELECT id FROM utility_companies WHERE alias = 'premiumdom')),
    -- Санкт-Петербург
    ('Невский пр., д. 88', ST_GeogFromText('SRID=4326;POINT(30.3200 59.9340)'), 2, 7, 'Центральный', 18, (SELECT id FROM utility_companies WHERE alias = 'nordstroy')),
    ('ул. Садовая, д. 14', ST_GeogFromText('SRID=4326;POINT(30.2980 59.9280)'), 2, 8, 'Адмиралтейский', 10, (SELECT id FROM utility_companies WHERE alias = 'nordstroy')),
    ('Невский пр., д. 90', ST_GeogFromText('SRID=4326;POINT(30.3200 59.9340)'), 2, 7, 'Центральный', 18, (SELECT id FROM utility_companies WHERE alias = 'nordstroy')),
    ('ул. Галошина, д. 15', ST_GeogFromText('SRID=4326;POINT(30.2980 59.9280)'), 2, 8, 'Адмиралтейский', 10, (SELECT id FROM utility_companies WHERE alias = 'nordstroy')),
    -- Казань
    ('ул. Белый, д. 42', ST_GeogFromText('SRID=4326;POINT(49.1260 55.7910)'), 3, 9, 'Вахитовский', 14, (SELECT id FROM utility_companies WHERE alias = 'kazaninvest')),
    ('ул. Тукая, д. 42', ST_GeogFromText('SRID=4326;POINT(49.1270 55.7910)'), 3, 9, 'Тукая', 14, NULL),
    ('ул. Баумана, д. 42', ST_GeogFromText('SRID=4326;POINT(49.1250 55.7910)'), 3, 9, 'Вахитовский', 14, (SELECT id FROM utility_companies WHERE alias = 'kazaninvest'));


-- ============================================================
-- 8. Категории квартир (flat_categories)
-- ============================================================
INSERT INTO flat_categories (name, room_count) VALUES 
    ('Студия', 0),
    ('1-комн.', 1),
    ('2-комн.', 2),
    ('3-комн.', 3),
    ('4-комн.', 4),
    ('5-комн.', 5),
    ('6+ комн.', 6);


-- ============================================================
-- 9. Объекты недвижимости (property)
-- ============================================================
INSERT INTO property (category_id, building_id, area) VALUES
    (1, 1, 20),  -- Тверская 25, студия
    (1, 2, 43),  -- Тверская 25, 2-к
    (1, 3, 93),  -- Арбат 36, пентхаус
    (1, 4, 34),  -- Арбат 36, 1-к
    (1, 5, 12),  -- Смоленская, студия
    (1, 6, 110), -- Смоленская, 2-к
    (1, 7, 223), -- Невский 88, 1-к
    (1, 8, 123), -- Невский 88, апартаменты
    (1, 9, 222), -- Садовая, 2-к
    (1, 10, 222); 

-- ============================================================
-- 10. Квартиры (flat)
-- ============================================================
INSERT INTO flat (property_id, floor, number, category_id) VALUES
    (1, 3, 12, (SELECT id FROM flat_categories WHERE room_count = 0)),
    (2, 7, 54, (SELECT id FROM flat_categories WHERE room_count = 1)),
    (3, 10, 99, (SELECT id FROM flat_categories WHERE room_count = 2)),
    (4, 2, 5, (SELECT id FROM flat_categories WHERE room_count = 2)),
    (5, 1, 2, (SELECT id FROM flat_categories WHERE room_count = 0)),
    (6, 6, 45, (SELECT id FROM flat_categories WHERE room_count = 3)),
    (7, 2, 8, (SELECT id FROM flat_categories WHERE room_count = 4)),
    (8, 9, 77, (SELECT id FROM flat_categories WHERE room_count = 5)),
    (9, 3, 22, (SELECT id FROM flat_categories WHERE room_count = 1)),
    (10, 3, 22, (SELECT id FROM flat_categories WHERE room_count = 1));

-- ============================================================
-- 11. Объявления (posters)
-- ============================================================
INSERT INTO posters (price, avatar_url, description, user_id, property_id, alias, created_at) VALUES
    (65000.00, 'https://design-cube.ru/wp-content/uploads/2022/06/View16-1.jpg', 'Уютная студия после капитального ремонта. Новая кухня, встроенные шкафы. Рядом метро.', 1, 1, 'studio-tverskaya', NOW() - INTERVAL '65 days'),
    (120000.00, 'https://design-cube.ru/wp-content/uploads/2025/07/4-2.webp', 'Светлая квартира с панорамным видом. Паркинг в подарок.', 1, 2, '2room-tverskaya', NOW() - INTERVAL '50 days'),
    (350000.00, 'https://design-cube.ru/wp-content/uploads/2026/01/23.webp', 'Уникальный пентхаус с открытой террасой 80 кв.м. Консьерж, закрытая территория.', 2, 3, 'penthouse-arbatskaya', NOW() - INTERVAL '95 days'),
    (85000.00, 'https://design-cube.ru/wp-content/uploads/2023/10/5-Mnogo-podushek-na-rozovom-divane-pridajut-emu-stil.jpg', 'Отличное место для жизни. Арбат в шаговой доступности. Мебель остаётся.', 2, 4, '1room-arbatskaya', NOW() - INTERVAL '35 days'),
    (55000.00, 'https://design-cube.ru/wp-content/uploads/2022/11/300_kuhnya_gostinaya_prihozhaya-15.jpg', 'Компактная студия для одного или пары. Первый этаж, высокие потолки.', 3, 5, 'studio-smolenskaya', NOW() - INTERVAL '25 days'),
    (75000.00, 'https://design-cube.ru/wp-content/uploads/2024/01/324_kuhnya-gostinaya-1.jpg', 'Квартира в историческом центре Петербурга. Дизайнерский ремонт 2024 года.', 4, 7, '1room-nevskiy', NOW() - INTERVAL '40 days'),
    (200000.00, 'https://design-cube.ru/wp-content/uploads/2022/04/276-Kuhnya-gostinnaya-2.jpg', 'Элитные апартаменты, 9 этаж. Потрясающий вид на Неву. Подземный паркинг.', 4, 8, 'apartments-neva-view', NOW() - INTERVAL '60 days'),
    (90000.00, 'https://design-cube.ru/wp-content/uploads/2022/04/278_obshhaya-zona_-8.jpg', 'Тихий двор, развитая инфраструктура. Школа и детсад в 5 минутах.', 5, 9, '2room-sadovaya', NOW() - INTERVAL '30 days'),
    (95000.00, 'https://design-cube.ru/wp-content/uploads/2025/07/5-4.webp', 'Большая семейная квартира. Все комнаты изолированы. Лоджия 8 кв.м.', 6, 10, '3room-kazan', NOW() - INTERVAL '35 days');

-- ============================================================
-- 12. Фотографии объявлений
-- ============================================================
-- Фотографии для объявлений
INSERT INTO poster_photos (img_url, sequence_order, poster_id) VALUES
    ('https://design-cube.ru/wp-content/uploads/2022/06/View16-1.jpg', 1, (SELECT id FROM posters WHERE alias = 'studio-tverskaya')),
    ('https://design-cube.ru/wp-content/uploads/2022/06/View10.jpg', 2, (SELECT id FROM posters WHERE alias = 'studio-tverskaya')),
    ('https://design-cube.ru/wp-content/uploads/2022/06/View01.jpg', 3, (SELECT id FROM posters WHERE alias = 'studio-tverskaya')),
    ('https://design-cube.ru/wp-content/uploads/2022/06/View04-1.jpg', 4, (SELECT id FROM posters WHERE alias = 'studio-tverskaya')),
    ('https://design-cube.ru/wp-content/uploads/2022/06/View02-4.jpg', 5, (SELECT id FROM posters WHERE alias = 'studio-tverskaya')),

    ('https://design-cube.ru/wp-content/uploads/2025/07/5.webp', 1, (SELECT id FROM posters WHERE alias = '2room-tverskaya')),
    ('https://design-cube.ru/wp-content/uploads/2025/07/6.webp', 2, (SELECT id FROM posters WHERE alias = '2room-tverskaya')),
    ('https://design-cube.ru/wp-content/uploads/2025/07/3.webp', 3, (SELECT id FROM posters WHERE alias = '2room-tverskaya')),
    ('https://design-cube.ru/wp-content/uploads/2025/07/4.webp', 4, (SELECT id FROM posters WHERE alias = '2room-tverskaya')),
    ('https://design-cube.ru/wp-content/uploads/2025/07/1.webp', 5, (SELECT id FROM posters WHERE alias = '2room-tverskaya')),
    ('https://design-cube.ru/wp-content/uploads/2025/07/3-1.webp', 6, (SELECT id FROM posters WHERE alias = '2room-tverskaya')),
    ('https://design-cube.ru/wp-content/uploads/2025/07/2-1.webp', 7, (SELECT id FROM posters WHERE alias = '2room-tverskaya')),
    ('https://design-cube.ru/wp-content/uploads/2025/07/4-2.webp', 8, (SELECT id FROM posters WHERE alias = '2room-tverskaya')),
    ('https://design-cube.ru/wp-content/uploads/2025/07/5-2.webp', 9, (SELECT id FROM posters WHERE alias = '2room-tverskaya')),
    ('https://design-cube.ru/wp-content/uploads/2025/07/3-3.webp', 10, (SELECT id FROM posters WHERE alias = '2room-tverskaya')),
    ('https://design-cube.ru/wp-content/uploads/2025/07/5-3.webp', 11, (SELECT id FROM posters WHERE alias = '2room-tverskaya')),

    ('https://design-cube.ru/wp-content/uploads/2026/01/8.webp', 1, (SELECT id FROM posters WHERE alias = 'penthouse-arbatskaya')),
    ('https://design-cube.ru/wp-content/uploads/2026/01/7.webp', 2, (SELECT id FROM posters WHERE alias = 'penthouse-arbatskaya')),
    ('https://design-cube.ru/wp-content/uploads/2026/01/5.webp', 3, (SELECT id FROM posters WHERE alias = 'penthouse-arbatskaya')),
    ('https://design-cube.ru/wp-content/uploads/2026/01/6.webp', 4, (SELECT id FROM posters WHERE alias = 'penthouse-arbatskaya')),
    ('https://design-cube.ru/wp-content/uploads/2026/01/2-7.webp', 5, (SELECT id FROM posters WHERE alias = 'penthouse-arbatskaya')),
    ('https://design-cube.ru/wp-content/uploads/2026/01/10-4.webp', 6, (SELECT id FROM posters WHERE alias = 'penthouse-arbatskaya')),
    ('https://design-cube.ru/wp-content/uploads/2026/01/4.2.webp', 7, (SELECT id FROM posters WHERE alias = 'penthouse-arbatskaya')),
    ('https://design-cube.ru/wp-content/uploads/2026/01/3-3.webp', 8, (SELECT id FROM posters WHERE alias = 'penthouse-arbatskaya')),
    ('https://design-cube.ru/wp-content/uploads/2026/01/1-5.webp', 9, (SELECT id FROM posters WHERE alias = 'penthouse-arbatskaya')),

    ('https://design-cube.ru/wp-content/uploads/2023/10/5-Mnogo-podushek-na-rozovom-divane-pridajut-emu-stil.jpg', 1, (SELECT id FROM posters WHERE alias = '1room-arbatskaya')),
    ('https://design-cube.ru/wp-content/uploads/2023/10/7-Fartuk-kuhni-vylozhen-belosnezhnoj-keramicheskoj-plitkoj.jpg', 2, (SELECT id FROM posters WHERE alias = '1room-arbatskaya')),
    ('https://design-cube.ru/wp-content/uploads/2023/10/13-Stilnuju-zonu-TV-ukrasil-vorsistyj-kover.jpg', 3, (SELECT id FROM posters WHERE alias = '1room-arbatskaya')),
    ('https://design-cube.ru/wp-content/uploads/2023/10/8-Stilnaya-vstroennaya-mebel-cveta-temnoj-polyni-funkcionalna.jpg', 4, (SELECT id FROM posters WHERE alias = '1room-arbatskaya')),
    ('https://design-cube.ru/wp-content/uploads/2023/10/2-Dvuhyarusna-krovat-pudrovogo-cveta-osnashhena-lestnicej-i-svetilnikami.jpg', 5, (SELECT id FROM posters WHERE alias = '1room-arbatskaya')),
    ('https://design-cube.ru/wp-content/uploads/2023/10/7-Beliznu-tualetnogo-stolika-podcherknula-zelenaya-podvesnaya-ljustra.jpg', 6, (SELECT id FROM posters WHERE alias = '1room-arbatskaya')),

    ('https://design-cube.ru/wp-content/uploads/2022/11/300_kuhnya_gostinaya_prihozhaya-15.jpg', 1, (SELECT id FROM posters WHERE alias = 'studio-smolenskaya')),
    ('https://design-cube.ru/wp-content/uploads/2022/11/300_kuhnya_gostinaya_prihozhaya-2.jpg', 2, (SELECT id FROM posters WHERE alias = 'studio-smolenskaya')),
    ('https://design-cube.ru/wp-content/uploads/2022/11/300_s.u-9.jpg', 3, (SELECT id FROM posters WHERE alias = 'studio-smolenskaya')),

    ('https://design-cube.ru/wp-content/uploads/2024/01/324_kuhnya-gostinaya-1.jpg', 1, (SELECT id FROM posters WHERE alias = '1room-nevskiy')),
    ('https://design-cube.ru/wp-content/uploads/2024/01/324_kuhnya-gostinaya-5.jpg', 2, (SELECT id FROM posters WHERE alias = '1room-nevskiy')),
    ('https://design-cube.ru/wp-content/uploads/2024/01/324_kuhnya-gostinaya-10.jpg', 3, (SELECT id FROM posters WHERE alias = '1room-nevskiy')),
    ('https://design-cube.ru/wp-content/uploads/2024/01/324_spalnya-2.jpg', 4, (SELECT id FROM posters WHERE alias = '1room-nevskiy')),
    ('https://design-cube.ru/wp-content/uploads/2024/01/324_spalnya-4.jpg', 5, (SELECT id FROM posters WHERE alias = '1room-nevskiy')),
    ('https://design-cube.ru/wp-content/uploads/2024/01/324_postirochnaya-9.jpg', 6, (SELECT id FROM posters WHERE alias = '1room-nevskiy')),

    ('https://design-cube.ru/wp-content/uploads/2022/04/276-Kuhnya-gostinnaya-2.jpg', 1, (SELECT id FROM posters WHERE alias = 'apartments-neva-view')),
    ('https://design-cube.ru/wp-content/uploads/2022/04/276-Kuhnya-gostinnaya-1.jpg', 2, (SELECT id FROM posters WHERE alias = 'apartments-neva-view')),
    ('https://design-cube.ru/wp-content/uploads/2022/04/276-Prihozhaya-1.jpg', 3, (SELECT id FROM posters WHERE alias = 'apartments-neva-view')),
    ('https://design-cube.ru/wp-content/uploads/2022/04/276-Sanuzel-1.jpg', 4, (SELECT id FROM posters WHERE alias = 'apartments-neva-view')),

    ('https://design-cube.ru/wp-content/uploads/2022/04/278_obshhaya-zona_-4.jpg', 1, (SELECT id FROM posters WHERE alias = '2room-sadovaya')),
    ('https://design-cube.ru/wp-content/uploads/2022/04/278_obshhaya-zona_-8.jpg', 2, (SELECT id FROM posters WHERE alias = '2room-sadovaya')),
    ('https://design-cube.ru/wp-content/uploads/2022/04/278_obshhaya-zona_-14.jpg', 3, (SELECT id FROM posters WHERE alias = '2room-sadovaya')),
    ('https://design-cube.ru/wp-content/uploads/2022/04/278_obshhaya-zona_-19.jpg', 4, (SELECT id FROM posters WHERE alias = '2room-sadovaya')),
    ('https://design-cube.ru/wp-content/uploads/2022/04/278_spalnya_-3.jpg', 5, (SELECT id FROM posters WHERE alias = '2room-sadovaya')),
    ('https://design-cube.ru/wp-content/uploads/2022/04/278_spalnya_-4.jpg', 6, (SELECT id FROM posters WHERE alias = '2room-sadovaya')),
    ('https://design-cube.ru/wp-content/uploads/2022/04/278_vanna_-1.jpg', 7, (SELECT id FROM posters WHERE alias = '2room-sadovaya')),
    ('https://design-cube.ru/wp-content/uploads/2022/04/278_vanna_-3.jpg', 8, (SELECT id FROM posters WHERE alias = '2room-sadovaya')),

    ('https://design-cube.ru/wp-content/uploads/2025/07/7-2.webp', 1, (SELECT id FROM posters WHERE alias = '3room-kazan')),
    ('https://design-cube.ru/wp-content/uploads/2025/07/5-4.webp', 2, (SELECT id FROM posters WHERE alias = '3room-kazan')),
    ('https://design-cube.ru/wp-content/uploads/2025/07/4-4.webp', 3, (SELECT id FROM posters WHERE alias = '3room-kazan')),
    ('https://design-cube.ru/wp-content/uploads/2025/07/3-4.webp', 4, (SELECT id FROM posters WHERE alias = '3room-kazan')),
    ('https://design-cube.ru/wp-content/uploads/2025/07/5-5.webp', 5, (SELECT id FROM posters WHERE alias = '3room-kazan')),
    ('https://design-cube.ru/wp-content/uploads/2025/07/4_Post.webp', 6, (SELECT id FROM posters WHERE alias = '3room-kazan')),
    ('https://design-cube.ru/wp-content/uploads/2025/07/40100_Post.webp', 7, (SELECT id FROM posters WHERE alias = '3room-kazan')),
    ('https://design-cube.ru/wp-content/uploads/2025/07/10000_Post.webp', 8, (SELECT id FROM posters WHERE alias = '3room-kazan')),
    ('https://design-cube.ru/wp-content/uploads/2025/07/30000_Post.webp', 9, (SELECT id FROM posters WHERE alias = '3room-kazan'));


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


-- ============================================================
-- 14. Удобства (facilities)
-- ============================================================
INSERT INTO facilities (name, alias) VALUES
    ('Wi-Fi', 'wifi'),
    ('Кондиционер', 'conditioner'),
    ('Стиральная машина', 'washing-machine'),
    ('Сушилка', 'dryer'),
    ('Гладильная доска', 'ironing-board'),
    ('Утюг', 'iron'),
    ('Телевизор', 'tv'),
    ('Холодильник', 'fridge'),
    ('Микроволновка', 'microwave'),
    ('Электроплита', 'stove'),
    ('Посудомойка', 'dishwasher'),
    ('Лифт', 'elevator'),
    ('Парковка', 'parking'),
    ('Консьерж', 'concierge'),
    ('Детская площадка', 'playground');


-- ============================================================
-- 15. Связь удобства → property (для каждого постера)
-- ============================================================
-- studio-tverskaya (property_id=1) - базовые удобства
INSERT INTO facility_property (property_id, facility_id) VALUES
    (1, (SELECT id FROM facilities WHERE alias='wifi')),
    (1, (SELECT id FROM facilities WHERE alias='conditioner')),
    (1, (SELECT id FROM facilities WHERE alias='washing-machine')),
    (1, (SELECT id FROM facilities WHERE alias='fridge')),
    (1, (SELECT id FROM facilities WHERE alias='microwave'));

-- 2room-tverskaya (property_id=2) - расширенный набор
INSERT INTO facility_property (property_id, facility_id) VALUES
    (2, (SELECT id FROM facilities WHERE alias='wifi')),
    (2, (SELECT id FROM facilities WHERE alias='conditioner')),
    (2, (SELECT id FROM facilities WHERE alias='washing-machine')),
    (2, (SELECT id FROM facilities WHERE alias='dryer')),
    (2, (SELECT id FROM facilities WHERE alias='dishwasher')),
    (2, (SELECT id FROM facilities WHERE alias='elevator')),
    (2, (SELECT id FROM facilities WHERE alias='parking'));

-- penthouse-arbatskaya (property_id=3) - премиум
INSERT INTO facility_property (property_id, facility_id) VALUES
    (3, (SELECT id FROM facilities WHERE alias='wifi')),
    (3, (SELECT id FROM facilities WHERE alias='conditioner')),
    (3, (SELECT id FROM facilities WHERE alias='dishwasher')),
    (3, (SELECT id FROM facilities WHERE alias='concierge')),
    (3, (SELECT id FROM facilities WHERE alias='parking')),
    (3, (SELECT id FROM facilities WHERE alias='elevator'));

-- 1room-arbatskaya (property_id=4) - стандарт
INSERT INTO facility_property (property_id, facility_id) VALUES
    (4, (SELECT id FROM facilities WHERE alias='wifi')),
    (4, (SELECT id FROM facilities WHERE alias='washing-machine')),
    (4, (SELECT id FROM facilities WHERE alias='fridge')),
    (4, (SELECT id FROM facilities WHERE alias='tv')),
    (4, (SELECT id FROM facilities WHERE alias='iron'));

-- studio-smolenskaya (property_id=5) - минимальный набор
INSERT INTO facility_property (property_id, facility_id) VALUES
    (5, (SELECT id FROM facilities WHERE alias='wifi')),
    (5, (SELECT id FROM facilities WHERE alias='fridge')),
    (5, (SELECT id FROM facilities WHERE alias='microwave'));

-- 1room-nevskiy (property_id=7) - питерский стандарт
INSERT INTO facility_property (property_id, facility_id) VALUES
    (7, (SELECT id FROM facilities WHERE alias='wifi')),
    (7, (SELECT id FROM facilities WHERE alias='conditioner')),
    (7, (SELECT id FROM facilities WHERE alias='washing-machine')),
    (7, (SELECT id FROM facilities WHERE alias='elevator'));

-- apartments-neva-view (property_id=8) - люкс
INSERT INTO facility_property (property_id, facility_id) VALUES
    (8, (SELECT id FROM facilities WHERE alias='wifi')),
    (8, (SELECT id FROM facilities WHERE alias='dishwasher')),
    (8, (SELECT id FROM facilities WHERE alias='concierge')),
    (8, (SELECT id FROM facilities WHERE alias='playground')),
    (8, (SELECT id FROM facilities WHERE alias='parking'));

-- 2room-sadovaya (property_id=9) - семейный
INSERT INTO facility_property (property_id, facility_id) VALUES
    (9, (SELECT id FROM facilities WHERE alias='wifi')),
    (9, (SELECT id FROM facilities WHERE alias='washing-machine')),
    (9, (SELECT id FROM facilities WHERE alias='dryer')),
    (9, (SELECT id FROM facilities WHERE alias='playground')),
    (9, (SELECT id FROM facilities WHERE alias='elevator'));

-- 3room-kazan (property_id=10) - просторная
INSERT INTO facility_property (property_id, facility_id) VALUES
    (10, (SELECT id FROM facilities WHERE alias='wifi')),
    (10, (SELECT id FROM facilities WHERE alias='conditioner')),
    (10, (SELECT id FROM facilities WHERE alias='dishwasher')),
    (10, (SELECT id FROM facilities WHERE alias='stove')),
    (10, (SELECT id FROM facilities WHERE alias='parking'));

 ANALYSE;
