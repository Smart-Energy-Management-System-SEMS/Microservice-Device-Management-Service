#!/bin/sh
set -eu

BOOTSTRAP_SERVER="${KAFKA_BOOTSTRAP_SERVER:-kafka:9092}"
TOPICS_CSV="${KAFKA_TOPICS:-}"

if [ -z "$TOPICS_CSV" ]; then
  echo "No Kafka topics configured in KAFKA_TOPICS; skipping topic creation."
  exit 0
fi

echo "Ensuring Kafka topics exist on ${BOOTSTRAP_SERVER}"

echo "$TOPICS_CSV" | tr ',' '\n' | while IFS= read -r raw_topic; do
  topic="$(echo "$raw_topic" | sed 's/^[[:space:]]*//; s/[[:space:]]*$//')"

  if [ -z "$topic" ]; then
    continue
  fi

  echo "Creating topic if missing: $topic"
  kafka-topics.sh \
    --bootstrap-server "$BOOTSTRAP_SERVER" \
    --create \
    --if-not-exists \
    --topic "$topic" \
    --partitions 1 \
    --replication-factor 1
done
