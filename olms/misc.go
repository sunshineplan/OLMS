package olms

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/sunshineplan/utils/mail"
)

var mailSetting mail.Setting

var (
	joinPath = filepath.Join
	dir      = filepath.Dir
)

// Backup database
func Backup() {
	log.Println("Start!")
	file := dump()
	defer os.Remove(file)
	if err := mail.SendMail(
		&mailSetting,
		fmt.Sprintf("My Bookmarks Backup-%s", time.Now().Format("20060102")),
		"",
		&mail.Attachment{FilePath: file, Filename: "database"},
	); err != nil {
		log.Fatalf("Failed to send mail: %v", err)
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
			log.Fatalf("File not found: %v", err)
		}
	}
	dropAll := joinPath(dir(Self), "scripts/drop_all.sql")
	execScript(dropAll)
	execScript(file)
	log.Println("Done!")
}
