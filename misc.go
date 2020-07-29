package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sunshineplan/utils/mail"
)

var mailSetting mail.Setting

func backup() {
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

func restore(file string) {
	log.Println("Start!")
	if file == "" {
		file = joinPath(dir(self), "scripts/schema.sql")
	} else {
		if _, err := os.Stat(file); err != nil {
			log.Fatalf("File not found: %v", err)
		}
	}
	dropAll := joinPath(dir(self), "scripts/drop_all.sql")
	execScript(dropAll)
	execScript(file)
	log.Println("Done!")
}
