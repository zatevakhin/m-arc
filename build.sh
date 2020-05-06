#!/usr/sbin/bash

go build -buildmode=plugin -o plugins/readmanga.so plugins/readmanga.go
go build -buildmode=plugin -o plugins/mintmanga.so plugins/mintmanga.go
go build -buildmode=plugin -o plugins/mangafox.so plugins/mangafox.go

go build -o app.elf main.go
