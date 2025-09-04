package cli

import (
	"flag"
	"fmt"
	"path/filepath"
)

type Config struct {
	RootPath string
	Port     int
	Version  bool
}

func ParseArgs(versionStr string) *Config {
	rootPath := flag.String("root", ".", "Root directory to serve files from")
	port := flag.Int("port", 8000, "Port to run the server on")
	showVersion := flag.Bool("version", false, "Show version and exit")
	showHelp := flag.Bool("help", false, "Show help and exit")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options]\n", "goserve")
		fmt.Println("Options:")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *showHelp {
		flag.Usage()
		return nil
	}

	if *showVersion {
		fmt.Println("GoServe version", versionStr)
		return nil
	}

	absRoot, err := filepath.Abs(*rootPath)
	if err != nil {
		panic("Error resolving path: " + err.Error())
	}

	return &Config{
		RootPath: absRoot,
		Port:     *port,
		Version:  *showVersion,
	}
}
