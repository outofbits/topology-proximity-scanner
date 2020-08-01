build-all: build-linux-amd64 build-linux-arm build-linux-arm64

build-linux-amd64:
	mkdir -p build/linux
	env GOOS=linux GOARCH=amd64 go build -o build/linux/topology-proximity-scanner
	cd build/linux && tar -czvf ../topology-proximity-scanner-linux-amd64.tar.gz topology-proximity-scanner
	rm -rf build/linux

build-linux-arm:
	mkdir -p build/linux
	env GOOS=linux GOARCH=arm go build -o build/linux/topology-proximity-scanner
	cd build/linux && tar -czvf ../topology-proximity-scanner-linux-arm.tar.gz topology-proximity-scanner
	rm -rf build/linux

build-linux-arm64:
	mkdir -p build/linux
	env GOOS=linux GOARCH=arm64 go build -o build/linux/topology-proximity-scanner
	cd build/linux && tar -czvf ../topology-proximity-scanner-linux-arm64.tar.gz topology-proximity-scanner
	rm -rf build/linux