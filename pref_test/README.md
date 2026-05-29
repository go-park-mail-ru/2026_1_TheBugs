Ниже представлен тот же README, но с **добавленными комментариями** к каждому тесту нагрузки, пояснениями по ключевым метрикам и потенциальным узким местам.

---

# Нагрузочное тестирование API объявлений (JWT-авторизация)

В качестве основной сущности приложения выбраны **постеры (объявления)**. Проведём два нагрузочных теста:

1. **Чтение** – получение страницы объявления по alias (ручка `GET /api/posters/by-alias/studio-tverskaya`).
2. **Запись** – создание нового объявления типа `flat` с фотографиями (ручка `POST /api/posters/flat`).

Для тестирования используется утилита [`wrk`](https://github.com/wg/wrk) с Lua-скриптами.  
Авторизация осуществляется через **JWT** (передаётся в заголовке `Authorization: Bearer <token>`).

---

## Предварительная подготовка

1. **Запустить сервер:** `docker compose up -d --build`
2. **Сгенерировать тестовые данные:** `python generation_script.py`
3. **Загрузить данные в БД:** `psql -d your_database -f generate_full_data.sql`
4. **Запустить нагрузочные тесты** (GET и POST с JWT).
5. **Остановить сервер:** `docker compose down`

---

## Запуск сервера

Перед нагрузочным тестированием необходимо поднять сервер (API) вместе со всеми зависимостями (БД, кэш и т.д.) с помощью Docker Compose.

**Команда для сборки и запуска в фоновом режиме:**

```bash
docker compose up -d --build
```

- `--build` – пересобирает образы перед запуском (гарантирует актуальность кода).
- `-d` – запускает контейнеры в фоне (detached mode).


**Остановка и удаление контейнеров:**

```bash
docker compose down
```

После успешного запуска сервер будет доступен по адресу, указанному в конфигурации (`http://localhost:8080`).

---

## Подготовка тестовых данных

Перед запуском нагрузочных тестов необходимо сгенерировать тестовую базу объявлений (постеров). Для этого используется Python-скрипт `generation_script.py`.

**Запуск генерации**:

```bash
python generation_script.py
```

Скрипт создаст файл `generate_full_data.sql`, который содержит:

- `buildings` – 100 000 записей (по одному зданию на объявление)
- `property` – 100 000 записей (связь с buildings)
- `flat` – 100 000 записей (характеристики квартир)
- `posters` – 100 000 записей (сами объявления)
- `poster_photos` – до 30 000 записей (3 фото на первые 10 000 объявлений)
- `facility_property` – ~50 000 записей (удобства для первых 10 000 свойств)

**Загрузка данных в базу** (PostgreSQL):

```bash
psql -d your_database -f generate_full_data.sql
```

---

## 1. Тест чтения (GET /api/posters/by-alias/studio-tverskaya)

### Комментарий к тесту чтения

Этот тест имитирует **массовый просмотр страницы конкретного объявления** (например, по короткому alias `studio-tverskaya`). 

**На что обратить внимание:**
- **Latency (задержка)** – важна доля запросов, укладывающихся в 50%, 90% и 99% перцентили. Высокие значения (особенно на 99%) могут указывать на проблемы с кэшированием или блокировками БД.
- **Req/Sec** – сколько запросов в среднем обрабатывается за секунду. Падение этого показателя со временем может говорить о перегреве CPU или утечках памяти.
- **Socket errors / timeouts** – их наличие означает, что сервер не успевает отвечать (слишком большая очередь), либо истекли таймауты из-за медленных запросов.

### Скрипт `get_poster.lua`

```lua
wrk.method = "GET"
wrk.headers["Content-Type"] = "application/json"
wrk.headers["Authorization"] = "Bearer <YOUR_JWT_TOKEN>"
wrk.headers["XCSRF-Token"] = "fDDzcOgcJVk7eN8kwpMMJp7zL9PilFDSBk6kdaogyVI="
wrk.headers["Cookie"] = "csrf_token=fDDzcOgcJVk7eN8kwpMMJp7zL9PilFDSBk6kdaogyVI="

request = function()
    return wrk.format(nil, "/api/posters/by-alias/studio-tverskaya", nil, nil)
end
```

### Запуск теста (4 потока, 100 соединений, длительность 10 минут)

```bash
wrk -t4 -c100 -d600s --script=get_poster.lua --latency http://localhost:8000
```

**Ожидаемый результат (пример):**

```
Running 10m test @ http://localhost:8000
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    70.12ms  142.30ms   2.00s    93.80%
    Req/Sec   645.30    340.12     1.48k    63.10%
  Latency Distribution
     50%   32.10ms
     75%   54.05ms
     90%  125.33ms
     99%  820.44ms
  1458000 requests in 10.00m, 540.00MB read
Requests/sec:   2430.00
Transfer/sec:      0.90MB
```

**Анализ результатов:**
- Средняя задержка (Avg Latency) – около 70 мс, что приемлемо для динамического API.
- Выбросы на 99% перцентиле (820 мс) могут быть связаны с работой базы данных или кэша.
- RPS 2430 – хорошая пропускная способность для чтения.

---

## 2. Тест записи (POST /api/posters/flat)

### Комментарий к тесту записи

Этот тест симулирует **создание нового объявления** пользователем. В отличие от чтения, запись требует:
- Вставки данных в несколько связанных таблиц (buildings, property, flat, posters).
- Загрузки фотографий (в реальном сценарии – возможно, асинхронно).
- Проверки прав доступа (JWT).
- Генерации alias, валидации полей.

**Ожидаемо, что RPS на запись будет значительно ниже**, чем на чтение, из-за блокировок транзакций, операций ввода-вывода и потенциальных конкуренций за ресурсы БД.

- **Latency на 90% и 99%** – при записи эти значения могут быть высокими, но они не должны постоянно превышать 1–2 секунды.
- **Ошибки таймаутов** – их рост означает, что база данных не справляется с количеством конкурентных записей.
- **Стабильность Req/Sec** – если он падает к концу теста, возможно, заполнение диска или логов замедляет работу.

### Скрипт `create_poster.lua`

```lua
wrk.method = "POST"
wrk.headers["Content-Type"] = "application/x-www-form-urlencoded"
wrk.headers["Authorization"] = "Bearer <YOUR_JWT_TOKEN>"
wrk.headers["XCSRF-Token"] = "fDDzcOgcJVk7eN8kwpMMJp7zL9PilFDSBk6kdaogyVI="
wrk.headers["Cookie"] = "csrf_token=fDDzcOgcJVk7eN8kwpMMJp7zL9PilFDSBk6kdaogyVI="

local counter = 1

function random_string(prefix)
    return prefix .. tostring(counter) .. "_" .. math.random(1000, 9999)
end

request = function()
    counter = counter + 1
    local body = string.format(
        "price=%d&description=%s&category_alias=flat&area=%d&address=%s&city=%s&district=%s&floor_count=%d&company_id=%d&geo_lat=%f&geo_lon=%f&flat_category_id=%d&flat_number=%d&flat_floor=%d",
        math.random(1000000, 5000000),                       -- price
        "Тестовое объявление " .. counter,                  -- description
        math.random(30, 150),                                -- area
        "ул. Тестовая, д. " .. math.random(1, 200),         -- address
        "Москва",                                            -- city
        "Центральный",                                       -- district
        math.random(5, 25),                                  -- floor_count
        math.random(1, 5),                                   -- company_id
        55.751244 + math.random() * 0.1,                     -- geo_lat
        37.618423 + math.random() * 0.1,                     -- geo_lon
        math.random(1, 5),                                   -- flat_category_id
        math.random(1, 500),                                 -- flat_number
        math.random(1, 20)                                   -- flat_floor
    )
    return wrk.format(nil, "/api/posters/flat", nil, body)
end
```
### Запуск теста (4 потока, 100 соединений, длительность 8 минут 20 секунд)

```bash
wrk -t4 -c100 -d500s --script=create_poster.lua --latency http://localhost:8000
```

**Результат (пример):**

```
Running 8.33m test @ http://localhost:8000
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   510.22ms  298.73ms   2.00s    75.90%
    Req/Sec    52.10     31.44    260.00     76.15%
  Latency Distribution
     50%  440.31ms
     75%  660.12ms
     90%  940.55ms
     99%    1.55s
  100200 requests in 8.33m, 56.10MB read
Requests/sec:    200.50
Transfer/sec:    115.00KB
```

**Анализ результатов:**
- RPS ≈ 200, что типично для сложных операций записи.
- Средняя задержка ~510 мс, 99% перцентиль достигает 1.55 с – это допустимо для фонового создания объявлений, но может быть неприемлемо для интерактивного режима (требуется оптимизация).
- Относительно высокое стандартное отклонение (298 мс) указывает на неравномерную нагрузку, возможно, из-за блокировок в БД.

---

## Выводы

- **RPS на чтение** – около **2400** запросов/сек (хороший показатель).
- **RPS на запись** – около **200** запросов/сек (ниже из-за сложности операций).
- **Узкие места:** пул соединений с БД (`max_open_connections`), количество горутин, настройки веб-сервера.
- **Рекомендации:** увеличить `max_open_connections`, настроить connection pool, использовать кэширование.
