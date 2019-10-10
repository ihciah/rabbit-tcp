rabbit:
	go build -o bin/rabbit cmd/rabbit.go

release:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -o bin/rabbit-linux-amd64 cmd/rabbit.go
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -a -ldflags '-extldflags "-static"' -o bin/rabbit-linux-386 cmd/rabbit.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -a -ldflags '-extldflags "-static"' -o bin/rabbit-linux-arm cmd/rabbit.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -ldflags '-extldflags "-static"' -o bin/rabbit-linux-arm64 cmd/rabbit.go
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -a -ldflags '-extldflags "-static"' -o bin/rabbit-windows-386.exe cmd/rabbit.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -o bin/rabbit-windows-amd64.exe cmd/rabbit.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=386 go build -a -ldflags '-extldflags "-static"' -o bin/rabbit-darwin-386 cmd/rabbit.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -o bin/rabbit-darwin-amd64 cmd/rabbit.go
