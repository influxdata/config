// Package config implements a unified configuration package for InfluxData
// products.  It is intended to unify existing patterns of configuration, while
// facilitating a migration to newer ways of handling configuration. It
// implements the superset of the two previously-used TOML configuration
// parsers to maintain backward compatibility. It provides TOML documentation
// capabilities through the use of "doc" struct tags. Default configuration is
// handled through the use of an exemplary default configuration struct that
// can be provided to config.NewConfig
package config
