# 启动本地面板：在 PowerShell 里运行  .\run-dev.ps1
# 面板地址 http://127.0.0.1:2053   首次登录 admin / admin   (Ctrl+C 停止)
Set-Location -Path $PSScriptRoot
Write-Host "启动面板… 浏览器打开 http://127.0.0.1:2053  (登录 admin/admin)" -ForegroundColor Green
& "$PSScriptRoot\xui.exe"
