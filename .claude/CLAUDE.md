# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

go-messaging is a Go library for creating structured log messages in JSON and slog-friendly formats. It provides a `Messenger` interface for generating consistent, structured messages with configurable fields.

## Common Commands

```bash
# Build
make build

# Run all tests
make test

# Run a single test
go test -v -run TestName ./...

# Run tests for a specific package
go test -v ./messenger/...

# Lint (runs golangci-lint, govulncheck, and cspell)
make lint

# Run only golangci-lint
make golangci-lint

# Coverage report (opens in browser)
make coverage

# Install development dependencies (one-time)
make dependencies-for-development

# Update Go dependencies
make dependencies

# Generate type definitions for all languages
make generate

# Clean build artifacts
make clean
```

## Architecture

### Core Packages

- **messenger/** - Main package implementing the `Messenger` interface
  - `New()` creates a messenger with configurable options
  - `NewJSON()` returns JSON-formatted messages
  - `NewSlog()` / `NewSlogLevel()` returns slog-compatible output
  - `NewError()` wraps messages as Go errors

- **parser/** - Parses JSON messages back into `SenzingMessage` structs

- **go/typedef/** - Auto-generated Go types from `message-RFC8927.json` schema (do not edit manually)

### Message Level Convention

Message numbers determine log level automatically:
- 0-999: TRACE
- 1000-1999: DEBUG
- 2000-2999: INFO
- 3000-3999: WARN
- 4000-4999: ERROR
- 5000-5999: FATAL
- 6000+: PANIC

### Code Generation

The `message-RFC8927.json` schema defines the message structure using JSON Type Definition. Running `make generate` creates type definitions in multiple languages (Go, Python, Java, TypeScript, etc.) using `jtd-codegen`.

## Testing Patterns

- Tests use `github.com/stretchr/testify` (assert/require)
- Table-driven tests with `testCasesForMessage` pattern
- Tests run in parallel (`test.Parallel()`)
- Example functions provide godoc documentation

## Environment Variables

- `SENZING_MESSAGE_FIELDS` - Comma-separated list of fields to include in output (e.g., "id, level, text") or "all" for all fields
