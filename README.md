# Go API TDD Starter

This is a demo project for building a well tested API in Go using TDD. This example showcases the following:
* Test coverage of http routers and handlers.
* Test coverage of DB Repositories.
* Test coverage of DB connection code.
* Running DB migrations (both forwards and backwards) as part of tests.
* Data fixtures setup and teardown on each integration test.


## Getting Started

1. Run `make dev-up` to run a local docker instance of postgres.
2. Run `make test` to run the tests and generate coverage.
3. Run `make coverage` to view the coverage report.
4. To tear down the stack, run `make dev-down`.
