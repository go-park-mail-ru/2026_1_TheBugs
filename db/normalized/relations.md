# Описание схемы базы данных

## Сущности

### `cities`
Хранит список городов.
- `id` — уникальный идентификатор
- `city_name` — название города (уникальное, до 40 символов)

### `metro_stations`
Хранит список станций метро.
- `id` — уникальный идентификатор
- `station_name` — название станции (уникальное, до 40 символов)

### `utility_companies`
Хранит информацию об управляющих компаниях (ЖК).
- `id` — уникальный идентификатор
- `company_name` — название компании
- `contacts` — контактная информация
- `geo` — географические координаты (тип PostGIS GEOGRAPHY POINT)
- `address` — адрес
- `created_at` — дата создания записи
- `updated_at` — дата последнего обновления записи
- `city_id` — идентификатор города (FK → `cities`)
- `metro_station_id` — идентификатор ближайшей станции метро (FK → `metro_stations`)

### `users`
Хранит информацию о пользователях платформы.
- `id` — уникальный идентификатор
- `email` — адрес электронной почты
- `hashed_password` — хеш пароля (может быть NULL при входе через OAuth)
- `provider` — провайдер OAuth (может быть NULL при входе по паролю)
- `created_at` — дата регистрации
- `updated_at` — дата последнего обновления
- `salt` — соль для хеширования пароля
- `company_id` — идентификатор управляющей компании (FK → `utility_companies`)

### `buildings`
Хранит информацию о домах/зданиях.
- `id` — уникальный идентификатор
- `geo` — географические координаты здания (тип PostGIS GEOGRAPHY POINT)
- `address` — адрес здания
- `district` — район
- `company_id` — управляющая компания (FK → `utility_companies`)
- `city_id` — город (FK → `cities`)
- `metro_station_id` — ближайшая станция метро (FK → `metro_stations`)

### `apartment_categories`
Хранит типы/категории помещений.
- `id` — уникальный идентификатор
- `name` — название категории (например, «квартира», «коммерческое», «паркинг»)
- `description` — описание категории

### `apartments`
Хранит информацию о помещениях (квартирах).
- `id` — уникальный идентификатор
- `floor` — этаж
- `number` — номер помещения
- `building_id` — здание, в котором находится помещение (FK → `buildings`)
- `category_id` — категория помещения (FK → `apartment_categories`)

### `posters`
Хранит объявления о продаже/аренде помещений.
- `id` — уникальный идентификатор
- `title` — заголовок объявления
- `price` — цена
- `avatar_url` — ссылка на главное изображение объявления
- `description` — описание
- `created_at` — дата создания объявления
- `updated_at` — дата последнего обновления
- `user_id` — пользователь, разместивший объявление (FK → `users`)
- `apartment_id` — помещение, о котором объявление (FK → `apartments`)

### `price_history`
Хранит историю изменения цены по объявлению.
- `id` — уникальный идентификатор
- `price` — цена на момент изменения
- `changed_at` — дата и время изменения цены
- `poster_id` — идентификатор объявления (FK → `posters`)

### `likes`
Хранит лайки пользователей на объявления.
- `id` — уникальный идентификатор
- `created_at` — дата создания лайка
- `updated_at` — дата последнего обновления
- `user_id` — пользователь, поставивший лайк (FK → `users`)
- `poster_id` — объявление, которому поставлен лайк (FK → `posters`)

### `views`
Хранит просмотры объявлений пользователями.
- `id` — уникальный идентификатор
- `created_at` — дата первого просмотра
- `updated_at` — дата последнего просмотра
- `user_id` — пользователь, просмотревший объявление (FK → `users`)
- `poster_id` — просмотренное объявление (FK → `posters`)

### `poster_photos`
Хранит фотографии, прикреплённые к объявлению.
- `id` — уникальный идентификатор
- `img_url` — ссылка на изображение
- `sequence_order` — порядковый номер фото в галерее
- `poster_id` — идентификатор объявления (FK → `posters`)

---

## Сущности и их функциональные зависимости

### cities
`{id} → {city_name}`  
**1NF**: ✓ **2NF**: ✓ **3NF**: ✓ **BCNF**: ✓

### metro_stations
`{id} → {station_name}`  
**1NF**: ✓ **2NF**: ✓ **3NF**: ✓ **BCNF**: ✓

### utility_companies
`{id} → {company_name, contacts, geo, address, created_at, updated_at, city_id, metro_station_id}`  
**1NF**: ✓ **2NF**: ✓ **3NF**: ✓ **BCNF**: ✓

### users
`{id} → {email, hashed_password, provider, created_at, updated_at, salt, company_id}`  
**1NF**: ✓ **2NF**: ✓ **3NF**: ✓ **BCNF**: ✓

### buildings
`{id} → {geo, address, district, company_id, city_id, metro_station_id}`  
**1NF**: ✓ **2NF**: ✓ **3NF**: ✓ **BCNF**: ✓

### apartment_categories
`{id} → {name, description}`  
**1NF**: ✓ **2NF**: ✓ **3NF**: ✓ **BCNF**: ✓

### apartments
`{id} → {floor, number, building_id, category_id}`  
**1NF**: ✓ **2NF**: ✓ **3NF**: ✓ **BCNF**: ✓

### posters
`{id} → {title, price, avatar_url, description, created_at, updated_at, user_id, apartment_id}`  
**1NF**: ✓ **2NF**: ✓ **3NF**: ✓ **BCNF**: ✓

### price_history
`{id} → {price, changed_at, poster_id}`  
**1NF**: ✓ **2NF**: ✓ **3NF**: ✓ **BCNF**: ✓

### likes
`{id} → {created_at, updated_at, user_id, poster_id}`  
**1NF**: ✓ **2NF**: ✓ **3NF**: ✓ **BCNF**: ✓

### views
`{id} → {created_at, updated_at, user_id, poster_id}`  
**1NF**: ✓ **2NF**: ✓ **3NF**: ✓ **BCNF**: ✓

### poster_photos
`{id} → {img_url, sequence_order, poster_id}`  
**1NF**: ✓ **2NF**: ✓ **3NF**: ✓ **BCNF**: ✓

---

## Общее заключение

Все отношения в схеме соответствуют требованиям:
- **1NF**: Все атрибуты атомарны, нет повторяющихся групп
- **2NF**: Нет частичных зависимостей (все PK состоят из одного атрибута)
- **3NF**: Нет транзитивных зависимостей
- **BCNF**: Все детерминанты являются потенциальными ключами

Схема полностью нормализована и не содержит аномалий вставки, обновления и удаления.

---

## ER-диаграмма базы данных

```mermaid
erDiagram

    cities {
        bigint id PK
        text city_name
    }

    metro_stations {
        bigint id PK
        text station_name
    }

    utility_companies {
        bigint id PK
        text company_name
        text contacts
        geography geo
        text address
        timestamptz created_at
        timestamptz updated_at
        bigint city_id FK
        bigint metro_station_id FK
    }

    users {
        bigint id PK
        text email
        text hashed_password
        text provider
        timestamptz created_at
        timestamptz updated_at
        text salt
        bigint company_id FK
    }

    buildings {
        bigint id PK
        geography geo
        text address
        text district
        bigint company_id FK
        bigint city_id FK
        bigint metro_station_id FK
    }

    apartment_categories {
        bigint id PK
        text name
        text description
    }

    apartments {
        bigint id PK
        smallint floor
        smallint number
        bigint building_id FK
        bigint category_id FK
    }

    posters {
        bigint id PK
        text title
        numeric price
        text avatar_url
        text description
        timestamptz created_at
        timestamptz updated_at
        bigint user_id FK
        bigint apartment_id FK
    }

    price_history {
        bigint id PK
        numeric price
        timestamptz changed_at
        bigint poster_id FK
    }

    likes {
        bigint id PK
        timestamptz created_at
        timestamptz updated_at
        bigint user_id FK
        bigint poster_id FK
    }

    views {
        bigint id PK
        timestamptz created_at
        timestamptz updated_at
        bigint user_id FK
        bigint poster_id FK
    }

    poster_photos {
        bigint id PK
        text img_url
        smallint sequence_order
        bigint poster_id FK
    }

    cities ||--o{ utility_companies : ""
    cities ||--o{ buildings : ""
    metro_stations ||--o{ utility_companies : ""
    metro_stations ||--o{ buildings : ""
    utility_companies ||--o{ users : ""
    utility_companies ||--o{ buildings : ""
    buildings ||--o{ apartments : ""
    apartment_categories ||--o{ apartments : ""
    apartments ||--o{ posters : ""
    users ||--o{ posters : ""
    posters ||--o{ price_history : ""
    posters ||--o{ likes : ""
    posters ||--o{ views : ""
    posters ||--o{ poster_photos : ""
    users ||--o{ likes : ""
    users ||--o{ views : ""
```
