clean: # Cleans all the generate files from go:generate
	@find  gen_* -exec rm {} \;

json:
	@cd pkg/schema && go generate -x
	@go generate -x

client:
	@cd cmd/client && go run *.go

server:
	@cd cmd/server && find . ! -name '*test.go' -exec go run {} \;
