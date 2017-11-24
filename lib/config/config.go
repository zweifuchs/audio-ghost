package config

import "flag"

type Config struct {
	rootDirectory string
	verbosity int
}

func (c *Config) ReadCmd() {
	flag.StringVar(&c.rootDirectory, "dir", "./", "directory to scan")
	flag.Parse()
}

func (c *Config) RootDirecotry() string {
	return c.rootDirectory
}