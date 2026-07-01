$goBin = 'D:\Programs\go\bin'
$cur = [Environment]::GetEnvironmentVariable('Path','User')
Write-Host "Current User PATH:"
Write-Host $cur
Write-Host "---"
$parts = $cur -split ';' | ForEach-Object { $_.Trim() } | Where-Object { $_ -ne '' }
if ($parts -contains $goBin) {
    Write-Host "$goBin is already in PATH, nothing to do."
    exit 0
}
$new = ($parts + $goBin) -join ';'
[Environment]::SetEnvironmentVariable('Path', $new, 'User')
Write-Host "Done. New User PATH:"
Write-Host ([Environment]::GetEnvironmentVariable('Path','User'))
Write-Host "---"
Write-Host "You MUST log out / sign out (or restart) for the change to take effect for new processes."
