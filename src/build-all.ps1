### topiik-server
go build topiik-server.go
if ($?){
    Write-Host "build succeed"
}
$hasExe= Get-Item .\target\topiik-server.exe | Measure-Object
if ($hasExe.count -gt 0){
    Remove-Item .\target\topiik-server.exe
}
Move-Item topiik-server.exe .\target\topiik-server.exe -Force

### topiik-cli
Set-Location .\cli
go build topiik-cli.go
if ($?){
    Write-Host "build cli succeed"
}
Set-Location ..
$hasExe= Get-Item .\target\topiik-cli.exe | Measure-Object
if ($hasExe.count -gt 0){
    Remove-Item .\target\topiik-cli.exe
}
Move-Item .\cli\topiik-cli.exe .\target\topiik-cli.exe -Force