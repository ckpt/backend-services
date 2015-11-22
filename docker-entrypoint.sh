#!/bin/sh
export CKPT_REDIS="${REDIS_PORT_6379_TCP_ADDR}:${REDIS_PORT_6379_TCP_PORT}"
export CKPT_AMQP_URL="amqp://guest:guest@${RABBITMQ_PORT_5672_TCP_ADDR}:${RABBITMQ_PORT_5672_TCP_PORT}"
sleep 5
./backend-services
