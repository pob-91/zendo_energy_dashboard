# Zendo Test Energy Dashboard


# TODO:

- Fix issue with timeout in data processor, need to switch to long poll not continuous or have no timeout
- Write better docs

### Requirements

- TODO

### Getting Started

In order to start this amazing energy dashboard run:

```sh
make start
```

this will start the stack, setup the database and seed the database with an initial data set.

Navigate to `http://localhost:3000` and see your dashboard!

# Assumptions

- Weather is taken from York
- Auth is very basic / non existent for internal networked services. In prod workloads this would be more advanced (certificats, security groups, networking rules e.t.c.)
