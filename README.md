Hello world!

local deployment
```
export DOCKER_BUILDKIT=1
sam build --use-container --cached --parallel
sam local start-api --env-vars env.json
```

updating lambda (wsl/zsh)
```
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap main.go
zip function.zip bootstrap main.go
aws lambda update-function-code --function-name my:lambda:arn --zip-file fileb://function.zip
```