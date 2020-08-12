@echo off
go build
go build -ldflags "-H windowsgui" -o olms.exe
