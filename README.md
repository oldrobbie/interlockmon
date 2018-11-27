# Interlock Mon

Monitoring UCP Interlock.

## How to Deploy

## Build and run from source

`docker-compose up --build --remove-orphans`

## Deploy from Docker Hub

`docker stack deploy -c .\prod.docker-compose.yml --prune interlockmon`

## Usage

Connect to the Splunk dashboard e.g. like this http://localhost:8000/en-US/app/search/interlock_mon.