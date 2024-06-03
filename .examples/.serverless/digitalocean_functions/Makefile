.PHONY: deploy clean

deploy:
	doctl serverless deploy ../digitalocean_functions

clean:
	rm -f ./packages/email-sender/send-post/.deployed
