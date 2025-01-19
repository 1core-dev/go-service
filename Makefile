# Select the preferred available shell, prioritizing ash > zsh > bash
SHELL_PATHS := /bin/ash /bin/zsh /bin/bash 
SHELL := $(firstword $(wildcard $(SHELL_PATHS)))

run: 
	go run app/services/sales-api/main.go | go run app/tooling/logfmt/main.go