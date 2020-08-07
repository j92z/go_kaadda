set CGO_ENABLED=0
set GOARCH=amd64
set GOOS=linux
go build -o kaadda main.go
go build -o kaadda-init initx/main.go