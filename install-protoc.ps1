# Protocol Buffers Compiler (protoc) インストールスクリプト
# Windows PowerShell用

$protocVersion = "25.1"  # 最新版のバージョンを指定
$protocZip = "protoc-$protocVersion-win64.zip"
$protocUrl = "https://github.com/protocolbuffers/protobuf/releases/download/v$protocVersion/$protocZip"
$installPath = "C:\protoc"

Write-Host "Installing Protocol Buffers Compiler v$protocVersion..." -ForegroundColor Green

# ダウンロードディレクトリ作成
$tempDir = "$env:TEMP\protoc_install"
New-Item -ItemType Directory -Force -Path $tempDir | Out-Null

# ダウンロード
Write-Host "Downloading from $protocUrl..."
$zipPath = "$tempDir\$protocZip"
Invoke-WebRequest -Uri $protocUrl -OutFile $zipPath

# 解凍
Write-Host "Extracting to $installPath..."
if (Test-Path $installPath) {
    Write-Host "Removing existing installation..."
    Remove-Item -Recurse -Force $installPath
}
Expand-Archive -Path $zipPath -DestinationPath $installPath -Force

# PATH環境変数の確認と追加
$binPath = "$installPath\bin"
$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($currentPath -notlike "*$binPath*") {
    Write-Host "Adding $binPath to PATH..."
    [Environment]::SetEnvironmentVariable("Path", "$currentPath;$binPath", "User")
    $env:Path = "$env:Path;$binPath"
    Write-Host "PATH updated. Please restart your terminal for changes to take effect." -ForegroundColor Yellow
} else {
    Write-Host "PATH already contains $binPath"
}

# クリーンアップ
Remove-Item -Recurse -Force $tempDir

# インストール確認
Write-Host "`nVerifying installation..." -ForegroundColor Green
& "$binPath\protoc.exe" --version

Write-Host "`nInstallation complete!" -ForegroundColor Green
Write-Host "You may need to restart your terminal or IDE for PATH changes to take effect." -ForegroundColor Yellow