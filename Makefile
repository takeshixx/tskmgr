build:
ifeq ($(OS),Windows_NT)
	go build -o ./bin/tsk.exe ./cmd/tsk
else
	go build -o ./bin/tsk ./cmd/tsk
endif

clean:
	rm -rf ./bin

test:
	go test -timeout 20s -cover ./...

run:
	go run cmd/tsk/main.go