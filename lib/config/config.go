package config

import "flag"

type Config struct {
	rootDirectory string
	verbosity int
	Rescan bool
}

var conf Config

func init() {
	conf.ReadCmd()
}

func (c *Config) ReadCmd() {
	flag.StringVar(&c.rootDirectory, "dir", "./", "directory to scan")
	flag.BoolVar(&c.Rescan, "forceRescan", false, "Force Drop Databases And Rescan")
	flag.Parse()
}

func (c *Config) RootDirecotry() string {
	return conf.rootDirectory
}

func GetConfig() *Config {
	return &conf
}

func Info() {
	return 

}