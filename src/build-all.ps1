
go build topiik-server.go
if ($?){
    Write-Host "build succeed"
}
$hasExe= Get-Item .\target\topiik-server.exe | Measure-Object
if ($hasExe.count -gt 0){
    Remove-Item .\target\topiik-server.exe
}
Move-Item topiik-server.exe .\target\topiik-server.exe -Force