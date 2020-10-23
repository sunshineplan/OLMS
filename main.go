package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"olms-go/olms"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sunshineplan/utils"
	"github.com/sunshineplan/utils/winsvc"
	"github.com/vharitonsky/iniflags"
)

func init() {
	winsvc.SetServiceName("OLMS")
	winsvc.SetDescription("Overtime and Leave Management System")
	winsvc.SetExecution(olms.Run)
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	fmt.Println(`
  run
		run OLMS web service mode (Default)
  backup
		run OLMS database`)
}

func main() {
	flag.Usage = usage
	flag.StringVar(&olms.UNIX, "unix", "", "UNIX-domain Socket")
	flag.StringVar(&olms.Host, "host", "0.0.0.0", "Server Host")
	flag.StringVar(&olms.Port, "port", "12345", "Server Port")
	flag.StringVar(&olms.SiteKey, "sitekey", "", "reCAPTCHA Site Key")
	flag.StringVar(&olms.SecretKey, "secretkey", "", "reCAPTCHA Secret Key")
	flag.StringVar(&olms.MailSetting.From, "from", "", "Backup sender")
	flag.StringVar(&olms.To, "to", "", "Backup receiver")
	flag.StringVar(&olms.MailSetting.Password, "password", "", "Backup sender password")
	flag.StringVar(&olms.MailSetting.SMTPServer, "server", "", "Backup sender server")
	flag.IntVar(&olms.MailSetting.SMTPServerPort, "bport", 587, "Backup sender server port")
	flag.StringVar(&olms.LogPath, "log", "", "Log Path")
	//flag.StringVar(&olms.LogPath, "log", filepath.Join(filepath.Dir(olms.Self), "access.log"), "Log Path")
	iniflags.SetConfigFile(filepath.Join(filepath.Dir(olms.Self), "config.ini"))
	iniflags.SetAllowMissingConfigFile(true)
	iniflags.Parse()

	if winsvc.IsWindowsService() {
		winsvc.RunService(false)
		return
	}

	var err error
	switch flag.NArg() {
	case 0:
		olms.Run()
	case 1:
		switch flag.Arg(0) {
		case "run", "debug":
			olms.Run()
		case "install":
			err = winsvc.InstallService()
		case "remove":
			err = winsvc.RemoveService()
		case "start":
			err = winsvc.StartService()
		case "stop":
			err = winsvc.StopService()
		case "backup":
			olms.Backup()
		case "init":
			if utils.Confirm("Do you want to initialize database?", 3) {
				olms.Restore("")
			}
		default:
			log.Fatalf("Unknown argument: %s", flag.Arg(0))
		}
	case 2:
		switch flag.Arg(0) {
		case "restore":
			if utils.Confirm("Do you want to restore database?", 3) {
				olms.Restore(flag.Arg(1))
			}
		default:
			log.Fatalf("Unknown arguments: %s", strings.Join(flag.Args(), " "))
		}
	default:
		log.Fatalf("Unknown arguments: %s", strings.Join(flag.Args(), " "))
	}
	if err != nil {
		log.Fatalf("failed to %s OLMS: %v", flag.Arg(0), err)
	}
}
