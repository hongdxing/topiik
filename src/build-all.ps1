
go build .\server\topiik-server.go .\server\executor.go .\server\vote.go
Remove-Item .\target\topiik-server.exe
Move-Item topiik-server.exe .\target\topiik-server.exe