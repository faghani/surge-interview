# surge interview project
Surge, price coefficient service

* fast and lightweight price coefficient calculator based on demand rate of tehran districts
* scalable service ( for the 21st century )
* light communication with other services and notify them (Near?) realtime ( Messaging based, thanks to Nats.io )
* find location district ( OSM data )
* metrics ( Soon :) )

## REST
REST v1 Document ( [Postman](https://documenter.getpostman.com/view/909541/SzS1SoR5) )

## Installation
Docker compose
```bash
docker-compose up -d
```

## Built With
* [Nats](https://nats.io) - Used to message-based communication with pricing service
* [Redis](https://redis.io) - Used to collect demands
