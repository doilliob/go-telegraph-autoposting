all:
	mkdir -p build
	go build -o build/autoposting
	cp configuration.yaml build/

clean:
	rm -rf build/

run:
	cd build && ./autoposting
