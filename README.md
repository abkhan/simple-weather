# Simple Weather

A simple weather web server.

The code is done only to make a simple web server and call to openWeather api work.

It does not fully follow all good rest api and programming best practice guidelines, using go-clean-architecture.

## Running

make dev-air


## Getting weather info

 curl -X POST -d '{"lat": "3.41", "lon": "4.15"}' localhost:8080/weather

 ## Things done intentionally

 only run the weather server, and not the default server that comes with go clean arch.


 ## Notes:

 This is work in progress, but the idea is that are framework like go-clean-arch helps in organizing the code in a way that makes a decoupled and priority inversion code easier.

 It also helps in testing and testability.


