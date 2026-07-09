# 完整重建面板单文件 exe：.\build.ps1
# 封装了本机的 Go / MinGW(gcc) 路径与 cgo 开关，改完代码跑这个就行
$ErrorActionPreference = 'Stop'
$env:Path = $env:Path + ';D:\Software\go\bin;D:\Software\mingw64\bin'
$env:CGO_ENABLED = '1'          # sqlite 驱动是 cgo，必须开
$go = 'D:\Software\go\bin\go.exe'

Write-Host '[1/2] 构建前端 (npm run build -> internal/web/dist)…' -ForegroundColor Cyan
Set-Location "$PSScriptRoot\frontend"
npm run build
if ($LASTEXITCODE -ne 0) { throw '前端构建失败' }

Write-Host '[2/2] 构建后端并嵌入前端 (go build -o xui.exe)…' -ForegroundColor Cyan
Set-Location $PSScriptRoot
& $go build -o xui.exe .
if ($LASTEXITCODE -ne 0) { throw 'go build 失败' }
Write-Host ("完成：{0:N1} MB  ->  .\run-dev.ps1 启动" -f ((Get-Item "$PSScriptRoot\xui.exe").Length/1MB)) -ForegroundColor Green
