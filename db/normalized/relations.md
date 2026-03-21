# Описание схемы базы данных

## Сущности

### profiles
**Описание**: Профили пользователей с контактной информацией.
- `id` — уникальный идентификатор (PK, GENERATED ALWAYS AS IDENTITY)
- `phone` — номер телефона (с CHECK регуляркой для РФ)
- `first_name` — имя (до 40 символов)
- `last_name` — фамилия (до 40 символов)
- `created_at` — дата создания (TIMESTAMPTZ DEFAULT NOW())
- `updated_at` — дата обновления (TIMESTAMPTZ DEFAULT NOW())

### users
**Описание**: Пользователи платформы с авторизацией.
- `id` — уникальный идентификатор (PK, GENERATED ALWAYS AS IDENTITY)
- `email` — email (UNIQUE, CHECK формат)
- `hashed_password` — хеш пароля (может быть NULL)
- `provider` — OAuth провайдер (может быть NULL)
- `provider_id` — ID в OAuth (UNIQUE)
- `salt` — соль для пароля
- `profile_id` — профиль (FK → profiles)
- `created_at`, `updated_at`

### refresh_tokens
**Описание**: Рефреш-токены пользователей.
- `id` — уникальный идентификатор (PK, GENERATED ALWAYS AS IDENTITY)
- `token_id` — UUID токена (UNIQUE)
- `user_id` — владелец (FK → users)
- `expires_at` — срок действия
- `created_at`

### cities
**Описание**: Города.
- `id` — уникальный идентификатор (PK)
- `city_name` — название (UNIQUE, до 40 символов)

### metro_stations
**Описание**: Станции метро с геоданными.
- `id` — уникальный идентификатор (PK)
- `station_name` — название (до 40 символов)
- `geo` — координаты (GEOGRAPHY(POINT, 4326))

### utility_companies
**Описание**: Управляющие компании (ЖК).
- `id` — уникальный идентификатор (PK)
- `company_name` — название
- `phone` — телефон (CHECK РФ формат)
- `geo` — координаты (GEOGRAPHY)
- `address` — адрес (CHECK кириллица, длина)
- `avatar_url` — аватар
- `alias` — уникальный алиас (UNIQUE)

### buildings
**Описание**: Дома/здания.
- `id` — уникальный идентификатор (PK)
- `address` — адрес (CHECK формат)
- `geo` — координаты (GEOGRAPHY)
- `city_id` — город (FK → cities)
- `metro_station_id` — метро (FK → metro_stations, NULLABLE)
- `district` — район (до 30 символов)
- `floor_count` — этажность (SMALLINT < 100)
- `company_id` — ЖК (FK → utility_companies, NULLABLE)

### property_categories
**Описание**: Типы недвижимости.
- `id` — уникальный идентификатор (PK)
- `name` — название (до 30 символов)

### property
**Описание**: Объекты недвижимости.
- `id` — уникальный идентификатор (PK)
- `category_id` — тип (FK → property_categories)
- `building_id` — дом (FK → buildings)
- `area` — площадь (NUMERIC(10,2) > 0)

### flat_categories
**Описание**: Категории квартир.
- `id` — уникальный идентификатор (PK)
- `name` — название

### flat
**Описание**: Детали квартир (расширение property).
- `property_id` — объект (PK, FK → property)
- `floor` — этаж (SMALLINT > 0)
- `number` — номер (INT > 0)
- `category_id` — категория (FK → flat_categories)

### posters
**Описание**: Объявления о продаже/аренде.
- `id` — уникальный идентификатор (PK)
- `price` — цена (NUMERIC(10,2) > 0)
- `avatar_url` — главное фото
- `description` — описание (до 500 символов)
- `alias` — уникальный slug (UNIQUE)
- `user_id` — автор (FK → users)
- `property_id` — объект (FK → property)
- `created_at`, `updated_at`

### poster_photos
**Описание**: Галерея объявлений.
- `id` — уникальный идентификатор (PK)
- `img_url` — ссылка на фото
- `sequence_order` — порядок (SMALLINT 1-15)
- `poster_id` — объявление (FK → posters)

### utility_companies_photos
**Описание**: Фото ЖК.
- `id` — уникальный идентификатор (PK)
- `img_url` — ссылка
- `sequence_order` — порядок (SMALLINT 1-15)
- `utility_company_id` — ЖК (FK)

### likes, views
**Описание**: Взаимодействия с объявлениями.
- `id` — PK
- `user_id`, `poster_id` — FK
- `created_at`, `updated_at`

## Нормализация

### profiles `{id} → {phone, first_name, last_name, created_at, updated_at}`
**1NF**: ✓  
**2NF**: ✓  
**3NF**: ✓  
**BCNF**: ✓

### users `{id} → {email, hashed_password, provider, provider_id, salt, profile_id, created_at, updated_at}`
**1NF**: ✓  
**2NF**: ✓  
**3NF**: ✓  
**BCNF**: ✓

### refresh_tokens `{id} → {token_id, user_id, expires_at, created_at}`
**1NF**: ✓  
**2NF**: ✓  
**3NF**: ✓  
**BCNF**: ✓

### cities `{id} → {city_name}`
**1NF**: ✓  
**2NF**: ✓  
**3NF**: ✓  
**BCNF**: ✓

### metro_stations `{id} → {station_name, geo}`
**1NF**: ✓  
**2NF**: ✓  
**3NF**: ✓  
**BCNF**: ✓

### utility_companies `{id} → {company_name, phone, geo, address, avatar_url, alias}`
**1NF**: ✓  
**2NF**: ✓  
**3NF**: ✓  
**BCNF**: ✓

### buildings `{id} → {address, geo, city_id, metro_station_id, district, floor_count, company_id}`
**1NF**: ✓  
**2NF**: ✓  
**3NF**: ✓  
**BCNF**: ✓

### property_categories `{id} → {name}`
**1NF**: ✓  
**2NF**: ✓  
**3NF**: ✓  
**BCNF**: ✓

### property `{id} → {category_id, building_id, area}`
**1NF**: ✓  
**2NF**: ✓  
**3NF**: ✓  
**BCNF**: ✓

### flat_categories `{id} → {name}`
**1NF**: ✓  
**2NF**: ✓  
**3NF**: ✓  
**BCNF**: ✓

### flat `{property_id} → {floor, number, category_id}`
**1NF**: ✓  
**2NF**: ✓  
**3NF**: ✓  
**BCNF**: ✓

### posters `{id} → {price, avatar_url, description, user_id, property_id, alias, created_at, updated_at}`
**1NF**: ✓  
**2NF**: ✓  
**3NF**: ✓  
**BCNF**: ✓

### poster_photos `{id} → {img_url, sequence_order, poster_id}`
**1NF**: ✓  
**2NF**: ✓  
**3NF**: ✓  
**BCNF**: ✓

### utility_companies_photos `{id} → {img_url, sequence_order, utility_company_id}`
**1NF**: ✓  
**2NF**: ✓  
**3NF**: ✓  
**BCNF**: ✓

### likes `{id} → {created_at, updated_at, user_id, poster_id}`
**1NF**: ✓  
**2NF**: ✓  
**3NF**: ✓  
**BCNF**: ✓

### views `{id} → {created_at, updated_at, user_id, poster_id}`
**1NF**: ✓  
**2NF**: ✓  
**3NF**: ✓  
**BCNF**: ✓

## ER-диаграмма
```mermaid
erDiagram
    %% Profiles
    profiles {
        bigint id PK
        text phone
        text first_name
        text last_name
        timestamptz created_at
        timestamptz updated_at
    }

    %% Users
    users {
        bigint id PK
        text email
        text hashed_password
        text provider
        text provider_id
        text salt
        bigint profile_id FK
        timestamptz created_at
        timestamptz updated_at
    }

    %% Refresh tokens
    refresh_tokens {
        bigint id PK
        uuid token_id
        bigint user_id FK
        timestamptz expires_at
        timestamptz created_at
    }

    %% Cities
    cities {
        bigint id PK
        text city_name
    }

    %% Metro stations
    metro_stations {
        bigint id PK
        text station_name
        geography geo
    }

    %% Utility companies
    utility_companies {
        bigint id PK
        text company_name
        text phone
        geography geo
        text address
        text avatar_url
        text alias
    }

    %% Utility companies photos
    utility_companies_photos {
        bigint id PK
        text img_url
        smallint sequence_order
        bigint utility_company_id FK
    }

    %% Buildings
    buildings {
        bigint id PK
        text address
        geography geo
        bigint city_id FK
        bigint metro_station_id FK
        text district
        smallint floor_count
        bigint company_id FK
    }

    %% Property categories
    property_categories {
        bigint id PK
        text name
        text description
    }

    %% Property
    property {
        bigint id PK
        bigint category_id FK
        bigint building_id FK
        numeric area
    }

    %% Flat categories
    flat_categories {
        bigint id PK
        text name
    }

    %% Flat
    flat {
        bigint property_id PK FK
        smallint floor
        int number
        bigint category_id FK
    }

    %% Posters
    posters {
        bigint id PK
        numeric price
        text avatar_url
        text description
        timestamptz created_at
        timestamptz updated_at
        bigint user_id FK
        bigint property_id FK
        text alias
    }

    %% Poster photos
    poster_photos {
        bigint id PK
        text img_url
        smallint sequence_order
        bigint poster_id FK
    }

    %% Likes
    likes {
        bigint id PK
        timestamptz created_at
        timestamptz updated_at
        bigint user_id FK
        bigint poster_id FK
    }

    %% Views
    views {
        bigint id PK
        timestamptz created_at
        timestamptz updated_at
        bigint user_id FK
        bigint poster_id FK
    }

    %% Relations
    profiles ||--o{ users : ""
    users ||--o{ refresh_tokens : ""
    users ||--o{ posters : ""
    users ||--o{ likes : ""
    users ||--o{ views : ""

    cities ||--o{ buildings : ""
    metro_stations ||--o{ buildings : ""
    utility_companies ||--o{ buildings : ""

    utility_companies ||--o{ utility_companies_photos : ""

    buildings ||--o{ property : ""
    property_categories ||--o{ property : ""
    property ||--o{ flat : ""
    
    property ||--o{ posters : ""
    posters ||--o{ poster_photos : ""
    posters ||--o{ likes : ""
    posters ||--o{ views : ""
```

## Индексы
```
- ФК: user_id, poster_id, property_id, building_id
- Поиск: price, area, alias, created_at (DESC)
- Гео: GiST на GEOGRAPHY(POINT, 4326)
- Композитные: idx_buildings_city_metro_company
```

