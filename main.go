package main

import (
	"flag"
	"log"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sunshineplan/olms-go/olms"
	"github.com/vharitonsky/iniflags"
)

func main() {
	flag.StringVar(&olms.UNIX, "UNIX", "", "UNIX-domain Socket")
	flag.StringVar(&olms.Host, "Host", "127.0.0.1", "Server Host")
	flag.StringVar(&olms.Port, "port", "12345", "Server Port")
	flag.StringVar(&olms.LogPath, "log", filepath.Join(filepath.Dir(olms.Self), "access.log"), "Log Path")
	iniflags.SetConfigFile(filepath.Join(filepath.Dir(olms.Self), "config.ini"))
	iniflags.SetAllowMissingConfigFile(true)
	iniflags.Parse()

	switch flag.NArg() {
	case 0:
		olms.Run()
	case 1:
		switch flag.Arg(0) {
		case "run":
			olms.Run()
		case "backup":
			olms.Backup()
		case "init":
			olms.Restore("")
		default:
			log.Fatalf("Unknown argument: %s", flag.Arg(0))
		}
	case 2:
		switch flag.Arg(0) {
		case "restore":
			olms.Restore(flag.Arg(1))
		default:
			log.Fatalf("Unknown arguments: %s", strings.Join(flag.Args(), " "))
		}
	default:
		log.Fatalf("Unknown arguments: %s", strings.Join(flag.Args(), " "))
	}
}
