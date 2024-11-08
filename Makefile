BINARY_NAME=Fynemdown
APP_NAME=Fynemdown
VERSION=0.0.1

build:
	rm -rf ${BINARY_NAME}
	rm -rf fyne-md
	fyne package -appVersion ${VERSION} -name ${APP_NAME} -release

run:
	go run .

clean:
	@echo "Cleaning..."
	@go clean
	@rm -rf ${BINARY_NAME}
	@echo "Cleaned!"

test:
	go test -v .