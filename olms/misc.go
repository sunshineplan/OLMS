package olms

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sunshineplan/utils/mail"
)

// Dialer contain backup account information
var Dialer mail.Dialer

// To is backup receiver
var To string

var (
	joinPath = filepath.Join
	dir      = filepath.Dir
)

// Backup database
func Backup() {
	log.Println("Start!")
	file := dump()
	defer os.Remove(file)
	if err := Dialer.Send(&mail.Message{
		To:          strings.Split(To, ","),
		Subject:     fmt.Sprintf("OLMS Backup-%s", time.Now().Format("20060102")),
		Attachments: []*mail.Attachment{{Path: file, Filename: "database"}},
	}); err != nil {
		log.Fatalln("Failed to send mail:", err)
	}
	log.Println("Done!")
}

// Restore database
func Restore(file string) {
	log.Println("Start!")
	if file == "" {
		file = joinPath(dir(Self), "scripts/schema.sql")
	} else {
		if _, err := os.Stat(file); err != nil {
			log.Fatalln("File not found:", err)
		}
	}
	dropAll := joinPath(dir(Self), "scripts/drop_all.sql")
	execScript(dropAll)
	execScript(file)
	log.Println("Done!")
}
