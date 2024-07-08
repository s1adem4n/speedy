cd frontend || exit

bun installl
bun run build

cd ..
go build -o bin/server main.go
