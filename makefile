build:
	go test
	go build

run: blackbox
	./blackbox

blackbox:
	go build

clean: 
	rm -rf blackbox
