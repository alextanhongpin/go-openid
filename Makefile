clean: # Cleans all the generate files from go:generate
	@find  gen_* -exec rm {} \;

json:
	@cd schema && go generate -x
	@go generate -x
