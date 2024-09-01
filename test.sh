#!/bin/bash

URL="http://localhost:8080/reservations"

PAYLOAD='{
    "room_id": "411",
    "start_time": "2024-09-01T10:00:00Z",
    "end_time": "2024-09-01T11:00:00Z"
}'

send_request() {
    curl -X POST -H "Content-Type: application/json" -d "$PAYLOAD" $URL
}

for i in {1..10}; do
    send_request &
done

wait

echo "All req done"


