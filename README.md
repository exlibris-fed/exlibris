# exlibris
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-v2.0%20adopted-ff69b4.svg)](code_of_conduct.md)

`exlibris` is a social network dedicated to tracking and discussing what you're reading.

This will be filled out more once we have a more fleshed out project.

## Developing
exlibris is written in [Vue.js](https://vuejs.org/) and [go](https://golang.org/), with ActivityPub support coming from [go-fed](https://github.com/go-fed/activity). It uses [postgreSQL](https://postgresql.org/) as a database, though you don't need it installed. [Docker](https://www.docker.com/) is required to build locally.

### Environment files
You need to provide environment variables for the database, app (Vue) and api (go). Copy any file ending in `.env.dist` to `.env` (ie `app.env.dist` to `app.env`) and fill in as necessary.

### Hot reload
To start a development server:

```bash
make run-local
```

This will start both the back- and front-ends with hot reloading, so making a change in a file will automatically recompile.

## History

exlibris was created during the 2020 employee hackathon at [ACV Auctions](https://acvauctions.com) and is being actively developed by its creators. We'd love to have your help too!
