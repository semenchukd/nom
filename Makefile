.PHONY: all
install:
	go build main.go && sudo cp main /usr/local/bin/nom

