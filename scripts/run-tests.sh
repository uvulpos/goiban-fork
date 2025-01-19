go mod download 
go install golang.org/x/tools/cmd/goimports
go install gotest.tools/gotestsum@latest

gotestsum --format testname ./...