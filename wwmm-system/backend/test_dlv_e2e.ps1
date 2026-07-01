$ErrorActionPreference = 'Stop'
Set-Location 'D:\GitHub\wwmm\wwmm-system\backend'
$out = 'C:\Users\ych1289\AppData\Local\Temp\dlv_out.log'
$err = 'C:\Users\ych1289\AppData\Local\Temp\dlv_err.log'
if (Test-Path $out) { Remove-Item $out }
if (Test-Path $err) { Remove-Item $err }
$proc = Start-Process -FilePath 'D:\Programs\go-workspace\bin\dlv.exe' `
    -ArgumentList @('debug','--headless=true','--listen=:2345','--api-version=2','--accept-multiclient','--continue','.') `
    -RedirectStandardOutput $out -RedirectStandardError $err -WindowStyle Hidden -PassThru
Write-Host "Started dlv PID: $($proc.Id)"
Start-Sleep -Seconds 4
Write-Host "--- procs ---"
Get-Process -Id $proc.Id -ErrorAction SilentlyContinue | Select-Object Id, ProcessName | Format-Table -AutoSize
Get-Process -Name 'wwmm-server' -ErrorAction SilentlyContinue | Select-Object Id, ProcessName | Format-Table -AutoSize
Write-Host "--- ports ---"
netstat -ano | Select-String ':2345 |:8080 ' | ForEach-Object { $_.Line }
Write-Host "--- HTTP probe /api/user/login ---"
try {
    $r = Invoke-WebRequest -Uri 'http://localhost:8080/api/user/login' -Method Post -ContentType 'application/json' -Body '{"username":"admin","password":"admin123"}' -UseBasicParsing -TimeoutSec 5
    Write-Host "HTTP $($r.StatusCode)"
    Write-Host $r.Content.Substring(0, [Math]::Min(200, $r.Content.Length))
} catch {
    Write-Host "ERROR: $($_.Exception.Message)"
}
Write-Host "--- dlv stdout ---"
if (Test-Path $out) { Get-Content $out } else { Write-Host "(no stdout file)" }
Write-Host "--- dlv stderr ---"
if (Test-Path $err) { Get-Content $err } else { Write-Host "(no stderr file)" }
Write-Host "--- cleanup ---"
Get-Process -Name 'dlv' -ErrorAction SilentlyContinue | Stop-Process -Force
Get-Process -Name 'wwmm-server' -ErrorAction SilentlyContinue | Stop-Process -Force
Write-Host "done"
