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