# Koment

Terminal comment viewer

## Install

```
go install github.com/limero/koment@latest
```

## Usage

```
koment <url>
koment --plain <url>   # print output and exit (no TUI)
koment demo
```

## Development

```
Run all tests offline:
go test ./...

Run all tests with actual calls to external services:
TEST_EXTERNAL=true go test ./...
```
