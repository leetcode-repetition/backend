Hello world!
```
export DOCKER_BUILDKIT=1
sam build --use-container --cached --parallel
sam local start-api --env-vars env.json
```