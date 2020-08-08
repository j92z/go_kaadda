go build -x -v -trimpath -ldflags "-s -w" -o kaadda main.go
go build -x -v -trimpath -ldflags "-s -w" -o kaadda-init initx/main.go
upx -9 kaadda
upx -9 kaadda-init