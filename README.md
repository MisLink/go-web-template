# go-web-template

## How to start

```
go install golang.org/x/tools/cmd/gonew@latest

gonew github.com/MisLink/go-web-template <your-module-path>

grep -R 'MODULE_NAME' -I -l --exclude README.md . | xargs sed -i 's/MODULE_NAME/<your-module-name>/g'
```
