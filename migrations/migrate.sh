#!/bin/sh

goose postgres "user=postgres password=postgres host=postgres dbname=fin-agg-db sslmode=disable" up -dir ./migrations
