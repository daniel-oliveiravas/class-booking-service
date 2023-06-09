# Class booking service
Backend service to manage members, classes and bookings

# Running locally
To run the project locally, you need to have Go installed. There are some helpers which use docker-compose, so it's nice
to have it installed as well.

The command below will start a postgres DB (on localhost:5432) and the service (on localhost:8080)
```shell
make run
```

# Running tests
Unit tests:
```shell
make test-unit
```
All tests (unit + integration):
```shell
make test
```

# Package structure
The project's structure follow a dependency order. 

Everything starts at the app package which holds all the API handlers (http server). 

It is followed by the business package which holds all the use cases implementation, and within it,
one can find a integration package holding implementation the interfaces defined in the use cases.

And finally we have the foundation package which holds boilerplate code, probably common to several services. This package
could be replaced by a company's "service-kit". It where one can find code to connect to a postgres database and logger configuration.
In the future, it would hold common web service implementations, like common middlewares (JWT, CORS), cache management, and more.

# Future improvements
Here you find some thoughts of what could improve in the future
- Security: Add authentication and authorization to all endpoints
    - Add a JWT middleware to validate all requests
- Add a cache layer for read endpoints
  - I'd add a cache layer (like Redis) between the service and the database to reduce read endpoints latency
- Tracing / Metrics
  - Distributed tracing to track different processes of the application
  - Metrics to measure endpoints latency (how much time in DB, how much time in Go code, etc.)
- Improve the error handling and logging
  - I've experimented with a new http framework (Gin), so there is room for improvements on error handling in
    the handlers package and also better logging
- Implement idempotency to POST endpoints to support request retries from clients
  - Right now, if the service gets the same POST request twice, it will create two entries for the same member.
  - It would be nice to use a mechanism to avoid that, like a Request-Id header which could be stored temporarily in a in-memory cache
- Improve database transaction calls
  - I'd probably implement a function to wrap the code to begin and commit a DB transaction to avoid duplicated code
