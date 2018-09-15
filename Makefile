.PHONY: restore update test cover cover_html clean

restore:
	dep ensure -v

update:
	dep ensure -update -v

test:
	go test -race ./...

cover:
	go test -cover -race ./...

cover_html:
	go test -race -coverprofile=profile.out
	go tool cover -html=profile.out -o cover.html

clean:
	-rm profile.out
	-rm cover.html