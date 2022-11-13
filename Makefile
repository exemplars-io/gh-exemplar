# clean befoe the build
clean :

	-rm gh-exemplar

build :
	@echo "build the exemplar extension in current directory"
	go build -o gh-exemplar main.go

# Build , Remove and Install the extension
install:	

	-gh extension remove exemplar
	
	gh extension install .

all: clean build install