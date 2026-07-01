Write-Host "=== Current environment ==="
Write-Host ""
Write-Host "User PATH (HKCU\Environment\Path):"
$user = [Environment]::GetEnvironmentVariable('Path','User')
if ([string]::IsNullOrEmpty($user)) { Write-Host "  (empty)" } else { $user -split ';' | ForEach-Object { Write-Host "  $_" } }
Write-Host ""
Write-Host "Machine PATH (HKLM\System\CurrentControlSet\Control\Session Manager\Environment\Path):"
$machine = [Environment]::GetEnvironmentVariable('Path','Machine')
$machine -split ';' | ForEach-Object { Write-Host "  $_" }
Write-Host ""
Write-Host "=== Go on disk ==="
$candidates = @(
    "D:\Programs\go\bin\go.exe",
    "C:\Program Files\Go\bin\go.exe",
    "C:\Go\bin\go.exe"
)
foreach ($c in $candidates) {
    if (Test-Path $c) {
        $v = & $c version 2>&1 | Select-Object -First 1
        Write-Host "FOUND: $c  ->  $v"
    } else {
        Write-Host "miss : $c"
    }
}
Write-Host ""
Write-Host "=== Where 'go' resolves from this PowerShell ==="
$which = Get-Command go -ErrorAction SilentlyContinue
if ($which) {
    Write-Host "  $($which.Source)"
} else {
    Write-Host "  (not found in PATH for this PowerShell session)"
}
