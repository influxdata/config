package config_test

import (
	"fmt"
	"log"
	"os"

	"github.com/influxdata/config"
)

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
