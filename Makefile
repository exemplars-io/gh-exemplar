
build:
	go build -o _output/bin/exemplars main.go

# Build , Remove and Install the extension
install:
	echo "building the exemplar extension"
	go build main.go

	gh extension remove exemplar
	
	gh extension install .
