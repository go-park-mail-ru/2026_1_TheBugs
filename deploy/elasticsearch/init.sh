#!/bin/bash

until curl -s http://localhost:9200/_cluster/health > /dev/null; do
  echo "Waiting for Elasticsearch..."
  sleep 2
done

curl -v -X DELETE "localhost:9200/posters" -s || true

curl -v -X PUT "localhost:9200/posters" \
  -H 'Content-Type: application/json' \
  -d '{
  "settings": {
    "analysis": {
      "filter": {
        "russian_stop": {"type": "stop", "stopwords": "_russian_"},
        "russian_stemmer": {"type": "stemmer", "language": "russian"},
        "russian_keywords": {"type": "keyword_marker", "keywords": ["квартира"]}
      },
      "analyzer": {
        "russian_analyzer": {
          "type": "custom",
          "tokenizer": "standard",
          "filter": ["lowercase", "russian_stop", "russian_keywords", "russian_stemmer"]
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "geo": { "type": "geo_point" },
      "description": {"type": "text", "analyzer": "russian_analyzer"},
      "city": {"type": "text", "analyzer": "russian_analyzer"},
      "station_name": {"type": "text", "analyzer": "russian_analyzer"},
      "district": {"type": "text", "analyzer": "russian_analyzer"},
      "address": {"type": "text", "analyzer": "russian_analyzer"},
      "company_name": {"type": "text", "analyzer": "russian_analyzer"},
      "facilities": {
        "type": "object",
        "properties": {
          "name": {"type": "text"},
          "alias": {"type": "text"}
        }
      }
    }
  }
}'

echo "✅ Индекс posters создан!"