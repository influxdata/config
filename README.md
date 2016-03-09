Config
======
This is the unified configuration management package for InfluxData. The
intention of this package is to unify existing patterns of interacting with
configuration across the elements of the TICK+E stack. As such, it implements
the superset of APIs from both `github.com/influxdata/toml` and
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
func ExampleConfig() {
	type DieConfig struct {
		Sides int    `toml:"sides" doc:"Number of sides on your die"`
		Name  string `toml:"name" doc:"The name of the die"`
	}

	type Config struct {
		Name        string    `toml:"name" doc:"Your name"`
		LuckyNumber int       `toml:"lucky_number" doc:"Your lucky number"`
		Die         DieConfig `toml:"die" doc:"Your die config"`
	}

	cfg := Config{
		Name:        "Tim",
		LuckyNumber: 42,
		Die: DieConfig{
			Sides: 20,
			Name:  "d20",
		},
	}
	// Write default configuration to stdout
	fmt.Println("Default Config:")
	config.NewEncoder(os.Stdout).Encode(cfg)

	// Decode a configfile
	err := config.Decode(`
name = "Jim"
[die]
    name = "d8"
	sides = 8
`, &cfg)
	if err != nil {
		log.Fatalf("Error loading configuration: %s", err.Error())
	}

	// Write default configuration to stdout
	fmt.Println("Decoded Config:")
	config.NewEncoder(os.Stdout).Encode(cfg)

	//Output:
	// Default Config:
	// name = "Tim" # Your name
	// lucky_number = 42 # Your lucky number
	// # Your die config
	// [die]
	//     sides = 20 # Number of sides on your die
	//     name = "d20" # The name of the die
	// Decoded Config:
	// name = "Jim" # Your name
	// lucky_number = 42 # Your lucky number
	// # Your die config
	// [die]
	//     sides = 8 # Number of sides on your die
	//     name = "d8" # The name of the die
	//

}
```

TODO
====

There are still somethings that should be moved to this package to ensure consitancy across packages.

* Environment variable handling. Currently each project allows for config to be specified via some form of env vars.
    We should make move that logic here so it is consistent across projects.

