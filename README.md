Config
======

This is the unified configuration management package for InfluxData. The
intention of this package is to unify existing patterns of interacting with
configuration across the elements of the TICK+E stack. As such, it implements
the superset of APIs from both `github.com/naoina/toml` and
`github.com/BurntSushi/toml`, while also providing a small API for the common
case of loading and storing configuration from a particular file. It also
provides wrapper types for formatting Durations and Sizes in TOML, which were
previously held in a sub-package within InfluxDB. Also provided is the ability
to document configuration fields using a "doc" struct tag. When the config
struct is marshalled as TOML, any doc struct tags found will be inserted as
TOML comments in an appropriate place to document the corresponding field.

Usage
=====

It's possible to use this package like `BurntSushi/toml` or `naoina/toml` for
backwards compatibility. However, the recommended usage going forward is as
follows:

```go
package main

import (
  "fmt"
  "log"

  "github.com/influxdata/config"
)

type Config struct {
  Name        string `toml:"name" doc:"Your name"`
  LuckyNumber int    `toml:"lucky_number" doc:"Your lucky number"`
}

func main() {
  ex := Config{"Tim", 42}
  cfg, err := NewConfig("config.toml", ex)
  if err != nil {
    log.Fatalf("Error loading configuration: %s", err.Error())
  }

  // Write default configuration to stdout
  fmt.Println("Example Config:")
  fmt.Print(cfg.EncodeDefault())

  // Load config
  var conf Config
  err := cfg.Decode(conf)
  if err != nil {
    log.Fatalf("Error parsing configuration: %s", err.Error())
  }

  fmt.Printf("Your name: %s\n", conf.Name)
  fmt.Printf("Your lucky number: %s\n", conf.LuckyNumber)
}
```
