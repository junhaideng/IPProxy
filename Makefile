output = build
filename = proxy_server

build-linux:
	echo "build on linux"
	go env -w GOOS=linux
	go build -o ${output}/${filename}"_linux"

build:
	go env -w CGO_ENABLED=0

	echo "build on windows"
	go env -w GOOS=windows
	go build  -o ${output}/${filename}".exe"

	echo "build on Mac"
	go env -w GOOS=darwin
	go build -o ${output}/${filename}"_mac"

	echo "build on linux"
	go env -w GOOS=linux
	go build -o ${output}/${filename}"_linux"

	echo "build finished"
