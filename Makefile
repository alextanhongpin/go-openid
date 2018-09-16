clean: # Cleans all the generate files from go:generate
	@find  gen_* -exec rm {} \;

json:
	@go generate -x 
