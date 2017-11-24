package config

import "flag"

type Config struct {
	rootDirectory string
	verbosity int
}

var conf Config

func init() {
	conf.ReadCmd()
}

func (c *Config) ReadCmd() {
	flag.StringVar(&c.rootDirectory, "dir", "./", "directory to scan")
	flag.Parse()
}

func (c *Config) RootDirecotry() string {
	return conf.rootDirectory
}

func GetConfig() *Config {
	return &conf
}