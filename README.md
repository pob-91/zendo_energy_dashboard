# Zendo Test Energy Dashboard


# TODO:

- Setup database auth and put in env files
- Write documentation and start up instructions and requirements
- Add database seeding to setup.sh script

### Requirements

- TODO

### Getting Started

In order to start this amazing energy dashboard run:

```sh
make up
make setup
```

this will start the stack, setup the database and seed the database with an initial data set.

# Assumptions

- Weather is taken from York
- Auth is very basic / non existent for internal networked services. In prod workloads this would be more advanced (certificats, security groups, networking rules e.t.c.)
