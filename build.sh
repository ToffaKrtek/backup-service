go build -o builds/CGO_ENABLED/backup-service-daemon cmd/daemon/main.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o builds/CGO_DISABLED/backup-service-cli cmd/cli/main.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o builds/CGO_DISABLED/backup-service-daemon cmd/daemon/main.go
go build -o builds/CGO_ENABLED/backup-service-cli cmd/cli/main.go
