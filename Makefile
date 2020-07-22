.PHONY: update master release setup update_master update_release build clean version

setup:
	git config --global --add url."git@gitlab.com:".insteadOf "https://gitlab.com/"

version:
	go run main.go generate
	mv version_vars.go cmd/version_vars.go

clean:
	rm -rf vendor/
	go mod vendor

update:
	-GOFLAGS="" go get -u all

build:
	go build ./...
	go mod tidy

update_release:
	GOFLAGS="" go get -u gitlab.com/elixxir/primitives@"XX-2433/MoveRingBuffer"
	GOFLAGS="" go get -u gitlab.com/elixxir/crypto@release
	GOFLAGS="" go get -u gitlab.com/elixxir/comms@"XX-2433/MoveRingBuffer"
	GOFLAGS="" go get -u gitlab.com/elixxir/client@"XX-2433/MoveRingBuffer"

update_master:
	GOFLAGS="" go get -u gitlab.com/elixxir/primitives@master
	GOFLAGS="" go get -u gitlab.com/elixxir/crypto@master
	GOFLAGS="" go get -u gitlab.com/elixxir/comms@master
	GOFLAGS="" go get -u gitlab.com/elixxir/client@master

master: clean update_master build version

release: clean update_release build version
