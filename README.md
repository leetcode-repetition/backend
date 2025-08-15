Hello world!

local deployment
```
export DOCKER_BUILDKIT=1
sam build --use-container --cached --parallel
sam local start-api --env-vars env.json
```

updating lambda (wsl/zsh)
```
./deploy.sh
```

updating shared package (wsl/zsh)
```
./update_shared.sh
```