build:
	GOOS=linux GOARCH=amd64 go build -o main lambda_email_sender.go

package: build
	zip function.zip main

deploy: package
	aws lambda update-function-code --function-name SendEmail --zip-file fileb://function.zip

clean:
	rm -f main function.zip
