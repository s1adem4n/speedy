cd frontend || exit

bun installl
bun run build

cd ..
go build -o build/server main.go
