EXT=
ifeq ($(OS),Windows_NT)
EXT=.exe
endif

all: server$(EXT) client$(EXT)

server$(EXT): cmd/server/main.go api
	go build -o $@ cmd/server/main.go

client$(EXT): cmd/client/main.go api
	go build -o $@ cmd/client/main.go

api: tsp-output/schema/openapi.yaml
	ogen tsp-output/schema/openapi.yaml

tsp-output/schema/openapi.yaml: main.tsp
	tsp compile .
