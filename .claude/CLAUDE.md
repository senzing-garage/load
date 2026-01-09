# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`load` is a Go CLI tool in the senzing-tools suite that pulls records from message queues (RabbitMQ, AWS SQS) and inserts them into the Senzing entity resolution database. It validates records against the Generic Entity Specification.

## Prerequisites

Senzing C library must be installed at:
- `/opt/senzing/er/lib` - shared objects
- `/opt/senzing/er/sdk/c` - SDK header files
- `/etc/opt/senzing` - configuration

See [How to Install Senzing for Go Development](https://github.com/senzing-garage/knowledge-base/blob/main/HOWTO/install-senzing-for-go-development.md) for setup instructions.

## Common Commands

```bash
# One-time dev setup (installs golangci-lint, gotestfmt, etc.)
make dependencies-for-development

# Install/update Go dependencies
make dependencies

# Lint (runs golangci-lint, govulncheck, cspell)
make lint

# Build binary (output in target/<os>-<arch>/)
make build

# Run tests (requires setup first)
make clean setup test

# Run single test
go test -v -run TestName ./path/to/package

# Coverage report (opens in browser)
make clean setup coverage

# Build Docker image
make docker-build

# Auto-fix lint issues
make fix
```

## Architecture

```
cmd/           CLI layer (cobra/viper)
  root.go      Main command definition with flags and execution
  context_*.go Platform-specific context variables

load/          Core loading logic
  load.go      BasicLoad struct implementing Load interface
               Orchestrates input reading and stats monitoring

input/         Input source adapters
  selector.go  URL-based routing to appropriate queue reader
  rabbitmq/    RabbitMQ consumer (amqp:// URLs)
  sqs/         AWS SQS consumer (sqs:// or https:// URLs)
```

### Key Dependencies

- `github.com/senzing-garage/go-cmdhelping` - CLI option handling
- `github.com/senzing-garage/go-queueing` - Queue consumers for RabbitMQ/SQS
- `github.com/senzing-garage/sz-sdk-go` - Senzing SDK interface
- `github.com/senzing-garage/go-sdk-abstract-factory` - Factory for Senzing engine instances

### Input URL Routing

The `input/selector.go` routes based on URL scheme:
- `amqp://` → RabbitMQ
- `sqs://` or `https://` → AWS SQS

## Linting Configuration

Uses golangci-lint v2 with extensive linter set. Config at `.github/linters/.golangci.yaml`. Notable exclusions:
- `exhaustruct` excludes `cobra.Command` and `load.BasicLoad`
- `ireturn` allows `logging.Logging` interface returns

## Testing

Tests require the Senzing database setup:
```bash
make setup  # Copies testdata/sqlite/G2C.db to /tmp/sqlite/
```

Test output uses gotestfmt for formatted display. Tests run sequentially (`-p 1`) due to shared resources.

## Environment Variables

- `LD_LIBRARY_PATH` - Path to Senzing libraries (default: `/opt/senzing/er/lib`)
- `SENZING_TOOLS_DATABASE_URL` - Database connection string
- `SENZING_LOG_LEVEL` - Logging level (TRACE, DEBUG, INFO, WARN, ERROR, FATAL, PANIC)
