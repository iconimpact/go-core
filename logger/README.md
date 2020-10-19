# Logger

Colorful logging with additional time and level output.

Feel free to add new functions or improve the existing code.

## Install

```bash
go get github.com/iconmobile-dev/go-core/logger
```

## Usage and Examples

```go
// init logger
log := logging.Logger{MinLevel: "verbose"}

// set settings from a config
log.MinLevel = cfg.MinLevel
log.TimeFormat = cfg.TimeFormat
log.UseColor = cfg.UseColor
log.ReportCaller = cfg.ReportCaller

// disable colors for tests
if test {
    log.UseColor = false
}

// use JSON output for production
if cfg.Env == "prod" {
    log.UseJSON = true
}
```