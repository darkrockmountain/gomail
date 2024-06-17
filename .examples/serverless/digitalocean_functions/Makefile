.PHONY: tidy wrap deploy clean

DOGOW_NPM_PACKAGE := digitalocean-go-wrapper

install-dogo-wrap:
	npm install -g $(DOGOW_NPM_PACKAGE)

wrap: install-dogo-wrap
	dogo-wrap -o gomail_serverless

deploy: wrap
	doctl serverless deploy gomail_serverless

clean:
	find . -type d -name ".deployed" -exec rm -rf {} +
	npm uninstall -g $(DOGOW_NPM_PACKAGE)
