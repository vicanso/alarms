.PHONY: default dev

# for dev
dev:
	fresh
build:
	packr2
	go build -o alarms