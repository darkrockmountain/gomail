.PHONY: tidy vendor deploy clean
tidy:
	go mod tidy

vendor: tidy
	go mod vendor

build: vendor
	@(GOOS=linux GOARCH=amd64 go build -o main gcf_email_sender.go)

deploy: build
	gcloud functions deploy go-http-function --gen2 --runtime=go121 --region=$(region) --source=. --entry-point=SendEmail --trigger-http --allow-unauthenticated
	@(rm -rf .vendor)
	
clean:
	@(rm -rf .vendor)
	@(rm -rf go.sum)
	@(rm -rf .main)
