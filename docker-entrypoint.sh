#!/bin/sh
export CKPT_REDIS="${REDIS_PORT_6379_TCP_ADDR}:${REDIS_PORT_6379_TCP_PORT}"
./backend-services
