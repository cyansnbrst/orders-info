#!/bin/bash

TOPIC="orders"
BROKER="kafka:9092"
FILE_PATH="message.json"

docker exec -i kafka kafka-console-producer.sh --broker-list $BROKER --topic $TOPIC < $FILE_PATH
