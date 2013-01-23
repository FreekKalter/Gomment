all: gomment.go
	go build gomment.go
install: gomment 
	cp gomment /usr/local/bin/
