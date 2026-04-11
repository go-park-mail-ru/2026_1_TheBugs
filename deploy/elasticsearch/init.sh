#!/bin/bash
set -e

until curl -s http://localhost:9200/_cluster/health > /dev/null; do
  echo "Waiting for Elasticsearch..."
  sleep 2
done

curl -X DELETE "localhost:9200/posters" -s > /dev/null || true

curl -X PUT "localhost:9200/posters" \
  -H 'Content-Type: application/json' << 'EOF'
{
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
      "description": {"type": "text", "analyzer": "russian_analyzer"},
      "facilities": {
        "type": "nested",
        "properties": {
          "name": {"type": "text", "analyzer": "russian_analyzer"}
        }
      }
    }
  }
}
EOF

echo "✅ Индекс posters создан!"