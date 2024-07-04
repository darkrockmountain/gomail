

build:
	GOOS=linux GOARCH=amd64 go build -o SendEmail/azure_email_sender azure_email_sender.go

deploy: build
	func azure functionapp publish SendEmail

clean:
	rm -f SendEmail/azure_email_sender
