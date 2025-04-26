EXT=
ifeq ($(OS),Windows_NT)
EXT=.exe
endif

all: server$(EXT) client$(EXT)

server$(EXT): api/*.go  cmd/server/main.go
	go build -o $@ cmd/server/main.go

client$(EXT):  api/*.go cmd/client/main.go
	go build -o $@ cmd/client/main.go

api/*.go: tsp-output/schema/openapi.yaml
	ogen tsp-output/schema/openapi.yaml

tsp-output/schema/openapi.yaml: main.tsp
	tsp compile .
