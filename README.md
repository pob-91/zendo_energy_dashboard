# Zendo Test Energy Dashboard

# Architecture

![archtecture](architecture.png "Architecture")

# TODO:

- Fix issue with timeout in data processor, need to switch to long poll not continuous or have no timeout
- Write better docs
- Add tests

### Requirements

- docker & docker compose
- curl
- an electricity maps api key (they have a free account tier)

### Getting Started

In order to start this amazing energy dashboard first setup required environment variables.

In each folder where you find a `.env.example` file, create a `.env` file and fill in the missing values.


```sh
make start
```

this will start the stack, setup the database and seed the database with an initial data set.

Navigate to `http://localhost:3000` and see your dashboard!

### Assumptions

- Weather is taken from York
- Auth is very basic / non existent for internal networked services. In prod workloads this would be more advanced (certificats, security groups, networking rules e.t.c.)

### Improvements

- Add logging middleware
- Add authentication and firewalls
- Rely less on the _changes feed and introduce a robust queue like RabbitMQ
- Database is fine to get started as is very lightweight however, would probably use something like MongoDB to get started and then if sclae and queries become an issue supplement with something like Hypertable.
