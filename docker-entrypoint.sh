#!/bin/sh
export CKPT_REDIS="redis:6379"
export CKPT_AMQP_URL="amqp://guest:guest@rabbitmq:5672"
sleep 5
./backend-services
