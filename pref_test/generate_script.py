import random
import os

NUM_POSTERS = 100_000
NUM_USERS = 6
NUM_BUILDINGS = NUM_POSTERS  
NUM_CITIES = 3 
BATCH_SIZE = 5000
OUTPUT_FILE = "generate_full_data.sql"
CLEAN_START = True

def random_description():
    templates = [
        "Уютная студия после ремонта. Новая кухня, встроенные шкафы. Рядом метро.",
        "Светлая квартира с панорамным видом. Паркинг в подарок.",
        "Отличное место для жизни в центре. Мебель остаётся.",
        "Квартира в тихом районе. Дизайнерский ремонт 2024 года.",
        "Компактная студия для одного или пары. Первый этаж.",
        "Элитные апартаменты, 9 этаж. Потрясающий вид. Подземный паркинг.",
        "Тихий двор, развитая инфраструктура. Школа и детсад рядом.",
        "Большая семейная квартира. Все комнаты изолированы.",
    ]
    return random.choice(templates)

def generate_area():
    r = random.random()
    if r < 0.15: return random.randint(20, 35)    
    elif r < 0.40: return random.randint(30, 50)   
    elif r < 0.70: return random.randint(50, 80)   
    elif r < 0.90: return random.randint(70, 120)  
    else: return random.randint(100, 250)          

def generate_price(area):
    base = random.randint(80000, 250000)
    return round(area * base / 1000, 2)

def random_alias(idx):
    words = ['studio', '1room', '2room', '3room', 'penthouse', 'apartments']
    cities = ['tverskaya', 'arbatskaya', 'smolenskaya', 'nevskiy', 'kazan']
    return f"{random.choice(words)}-{random.choice(cities)}-{idx:05d}"

def random_address(city_id):
    streets = ['Тверская', 'Арбат', 'Смоленская', 'Невский пр.', 'Садовая', 'Баумана', 'Ленина', 'Кирова']
    house = random.randint(1, 200)
    return f"ул. {random.choice(streets)}, д. {house}"

def random_metro(city_id):
    moscow_metro = ['Арбатская', 'Тверская', 'Белорусская', 'Киевская', 'Лубянка']
    spb_metro = ['Невский проспект', 'Площадь Восстания', 'Сенная площадь']
    kazan_metro = ['Казань Центральная', 'Площадь Ленина']
    if city_id == 1: return random.choice(moscow_metro)
    elif city_id == 2: return random.choice(spb_metro)
    else: return random.choice(kazan_metro)

def random_district(city_id):
    if city_id == 1: return random.choice(['Центральный', 'Арбат', 'Хамовники', 'Замоскворечье'])
    elif city_id == 2: return random.choice(['Центральный', 'Адмиралтейский', 'Васильевский'])
    else: return random.choice(['Вахитовский', 'Тукая', 'Ново-Савиновский'])

def main():
    with open(OUTPUT_FILE, 'w', encoding='utf-8') as f:
        if CLEAN_START:
            f.write("-- ============================================================\n")
            f.write("-- Clean existing data (order respects foreign keys)\n")
            f.write("-- ============================================================\n")
            f.write("TRUNCATE TABLE poster_photos, facility_property, posters, flat, property, buildings RESTART IDENTITY CASCADE;\n\n")
        
        # ------------------------------------------------------
        # Buildings (NUM_BUILDINGS = NUM_POSTERS = 100 000)
        f.write("-- ============================================================\n")
        f.write(f"-- Buildings ({NUM_BUILDINGS:,} records)\n")
        f.write("-- ============================================================\n\n")
        
        total_build_batches = NUM_BUILDINGS // BATCH_SIZE
        for batch in range(total_build_batches):
            f.write(f"-- Buildings batch {batch+1}/{total_build_batches}\n")
            f.write("INSERT INTO buildings (address, geo, city_id, metro_station_id, district, floor_count, company_id) VALUES\n")
            batch_lines = []
            for i in range(BATCH_SIZE):
                idx = batch * BATCH_SIZE + i + 1  # building_id
                city_id = random.randint(1, NUM_CITIES)
                address = random_address(city_id)
                district = random_district(city_id)
                floor_count = random.randint(5, 30)
                company_id = random.randint(1, 5) if random.random() > 0.2 else None
                # координаты
                if city_id == 1:
                    lon = round(random.uniform(37.5, 37.7), 6)
                    lat = round(random.uniform(55.7, 55.85), 6)
                elif city_id == 2:
                    lon = round(random.uniform(30.2, 30.4), 6)
                    lat = round(random.uniform(59.9, 60.0), 6)
                else:
                    lon = round(random.uniform(49.0, 49.2), 6)
                    lat = round(random.uniform(55.7, 55.9), 6)
                geo = f"ST_SetSRID(ST_MakePoint({lon}, {lat}), 4326)"
                metro_id = random.randint(1, 10)
                company_str = 'NULL' if company_id is None else str(company_id)
                
                line = f"    ('{address}', {geo}, {city_id}, {metro_id}, '{district}', {floor_count}, {company_str})"
                if i < BATCH_SIZE - 1:
                    line += ","
                batch_lines.append(line)
            f.write("\n".join(batch_lines))
            f.write("\n;\n\n")
        f.write(f"-- Total buildings: {NUM_BUILDINGS:,}\n\n")
        
        # ------------------------------------------------------
        # Property (одна запись на каждое здание, building_id уникальный)
        f.write("-- ============================================================\n")
        f.write(f"-- Property ({NUM_POSTERS:,} records) -- each building has one property\n")
        f.write("-- ============================================================\n\n")
        property_areas = []
        total_prop_batches = NUM_POSTERS // BATCH_SIZE
        for batch in range(total_prop_batches):
            f.write(f"-- Property batch {batch+1}/{total_prop_batches}\n")
            f.write("INSERT INTO property (category_id, building_id, area) VALUES\n")
            batch_lines = []
            for i in range(BATCH_SIZE):
                idx = batch * BATCH_SIZE + i + 1 
                category_id = 1
                building_id = idx
                area = generate_area()
                property_areas.append(area)
                line = f"    ({category_id}, {building_id}, {area})"
                if i < BATCH_SIZE - 1:
                    line += ","
                batch_lines.append(line)
            f.write("\n".join(batch_lines))
            f.write("\n;\n\n")
        f.write(f"-- Total property: {NUM_POSTERS:,}\n\n")
        
        # ------------------------------------------------------
        # Flat (одна квартира на каждое property)
        f.write("-- ============================================================\n")
        f.write(f"-- Flat ({NUM_POSTERS:,} records)\n")
        f.write("-- ============================================================\n\n")
        
        for batch in range(total_prop_batches):
            f.write(f"-- Flat batch {batch+1}/{total_prop_batches}\n")
            f.write("INSERT INTO flat (property_id, floor, number, category_id) VALUES\n")
            batch_lines = []
            for i in range(BATCH_SIZE):
                idx = batch * BATCH_SIZE + i + 1  # property_id
                property_id = idx
                floor = random.randint(1, 30)
                flat_number = random.randint(1, 500)
                area = property_areas[idx-1]
                if area < 35: category_id = 1
                elif area < 50: category_id = 2
                elif area < 80: category_id = 3
                elif area < 120: category_id = 4
                else: category_id = 5
                line = f"    ({property_id}, {floor}, {flat_number}, {category_id})"
                if i < BATCH_SIZE - 1:
                    line += ","
                batch_lines.append(line)
            f.write("\n".join(batch_lines))
            f.write("\n;\n\n")
        f.write(f"-- Total flat: {NUM_POSTERS:,}\n\n")
        
        f.write("-- ============================================================\n")
        f.write(f"-- Posters ({NUM_POSTERS:,} records)\n")
        f.write("-- ============================================================\n\n")
        
        for batch in range(total_prop_batches):
            f.write(f"-- Posters batch {batch+1}/{total_prop_batches}\n")
            f.write("INSERT INTO posters (price, avatar_url, description, user_id, property_id, alias, created_at) VALUES\n")
            records = []
            for i in range(BATCH_SIZE):
                idx = batch * BATCH_SIZE + i + 1
                area = property_areas[idx-1]
                price = generate_price(area)
                user_id = random.randint(1, NUM_USERS)
                property_id = idx
                alias = random_alias(idx)
                days_ago = random.randint(1, 365)
                avatar = f"https://design-cube.ru/wp-content/uploads/202{random.randint(2,6)}/{random.randint(1,12):02d}/{random.randint(1,300)}.jpg"
                desc = random_description().replace("'", "''")
                line = f"    ({price:.2f}, '{avatar}', '{desc}', {user_id}, {property_id}, '{alias}', NOW() - INTERVAL '{days_ago} days')"
                if i < BATCH_SIZE - 1:
                    line += ","
                records.append(line)
            f.write("\n".join(records))
            f.write("\n;\n\n")
        f.write(f"-- Total posters: {NUM_POSTERS:,}\n\n")
        
        f.write("-- ============================================================\n")
        f.write("-- Poster Photos (3 photos per poster for first 10,000 posters)\n")
        f.write("-- ============================================================\n\n")
        
        photo_count = 0
        PHOTO_BATCH = 1000
        for start in range(1, 10001, PHOTO_BATCH):
            end = min(start + PHOTO_BATCH - 1, 10000)
            f.write(f"-- Photos for posters {start} to {end}\n")
            f.write("BEGIN;\n")
            for poster_id in range(start, end+1):
                for seq in range(1, 4):
                    url = f"https://design-cube.ru/wp-content/uploads/202{random.randint(2,6)}/{random.randint(1,12):02d}/{random.randint(1,300)}.jpg"
                    f.write(f"INSERT INTO poster_photos (img_url, sequence_order, poster_id) VALUES ('{url}', {seq}, {poster_id});\n")
                    photo_count += 1
            f.write("COMMIT;\n\n")
        f.write(f"-- Total poster_photos: {photo_count}\n\n")
        f.write("-- ============================================================\n")
        f.write("-- Facility Property (sample for first 10,000 properties)\n")
        f.write("-- ============================================================\n\n")
        
        facilities = list(range(1, 16))
        fac_count = 0
        FAC_BATCH = 1000
        for start in range(1, 10001, FAC_BATCH):
            end = min(start + FAC_BATCH - 1, 10000)
            f.write(f"-- Facilities for properties {start} to {end}\n")
            f.write("BEGIN;\n")
            for property_id in range(start, end+1):
                num_fac = random.randint(3, 8)
                selected = random.sample(facilities, num_fac)
                for fac_id in selected:
                    f.write(f"INSERT INTO facility_property (property_id, facility_id) VALUES ({property_id}, {fac_id});\n")
                    fac_count += 1
            f.write("COMMIT;\n\n")
        f.write(f"-- Total facility_property: {fac_count}\n")
        

        f.write("-- ============================================================\n")
        f.write(f"-- SUMMARY\n")
        f.write(f"-- Buildings: {NUM_BUILDINGS:,}\n")
        f.write(f"-- Property: {NUM_POSTERS:,}\n")
        f.write(f"-- Flat: {NUM_POSTERS:,}\n")
        f.write(f"-- Posters: {NUM_POSTERS:,}\n")
        f.write(f"-- Poster Photos: {photo_count}\n")
        f.write(f"-- Facility Property: {fac_count}\n")
        f.write("-- ============================================================\n")
    
    print(f"Generated: {OUTPUT_FILE}")
    print(f"Buildings: {NUM_BUILDINGS:,}")
    print(f"Property: {NUM_POSTERS:,}")
    print(f"Flat: {NUM_POSTERS:,}")
    print(f"Posters: {NUM_POSTERS:,}")
    print(f"Poster Photos: {photo_count:,}")
    print(f"Facility Property: {fac_count:,}")
    
    file_size = os.path.getsize(OUTPUT_FILE)
    print(f"File size: {file_size / (1024*1024):.1f} MB")

if __name__ == "__main__":
    main()