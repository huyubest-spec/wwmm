$dir = 'D:\Programs\go-workspace\bin'
$cur = [Environment]::GetEnvironmentVariable('Path','User')
$parts = $cur -split ';' | ForEach-Object { $_.Trim() } | Where-Object { $_ -ne '' }
if ($parts -contains $dir) {
    Write-Host "$dir is already in PATH, nothing to do."
    exit 0
}
$new = ($parts + $dir) -join ';'
[Environment]::SetEnvironmentVariable('Path', $new, 'User')
Write-Host "Added: $dir"
Write-Host "New User PATH ends with: ...$(($parts | Select-Object -Last 2) -join ';');$dir"
