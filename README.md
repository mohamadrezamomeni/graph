# Graph

# Project Build & Run

## Configuration

in this project, we don’t have any configurations that need to be injected

## Build

Run the following command to build pbserver and pbclient.

```bash
make build
```

The compiled binaries will be placed in the `bin/` directory.

---

## Run

After building, you can run each binary:

```bash
./bin/pbserver

./bin/pbclient
```

You would run them separately after building.

---

## Testing

For testing, you must run following command:

```bash
go test ./...
```

---

## Packages Used

- **[sql-migrate](https://github.com/rubenv/sql-migrate)** → A database migration tool for Go that helps manage schema changes safely and in a version-controlled manner.
- **[logrus](https://github.com/sirupsen/logrus) v1.9.3** → A structured logger for Go, providing leveled logging and easy integration with log management systems.
- **[ozzo-validation](https://github.com/go-ozzo/ozzo-validation) v1.37.0** → A ozzo-validation is a popular Go validation library that makes it easier to validate struct fields, function parameters, and user input in a clean, declarative way..
- **[echo](https://github.com/labstack/echon) vv4.13.4** → A high-performance, minimalist web framework for Go, featuring routing, middleware support, and easy JSON handling.
