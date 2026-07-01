$ErrorActionPreference = "Continue"
$BaseUrl = "http://localhost:8080"
$adminBody = '{"username":"admin","password":"admin123"}'
$adminResp = Invoke-RestMethod -Uri "$BaseUrl/api/user/login" -Method POST -Body $adminBody -ContentType "application/json"
$adminToken = $adminResp.data.token

# 直接更新 status
$ids = @(1, 3, 4, 5, 6, 7)
foreach ($id in $ids) {
    $body = '{"approve":true,"comment":"内容符合规范，审核通过"}'
    try {
        $r = Invoke-RestMethod -Uri "$BaseUrl/api/photo/$id/audit" -Method POST -Body $body -ContentType "application/json" -Headers @{Authorization="Bearer $adminToken"}
        Write-Host "  [OK] photo #$id - $($r.message)"
    } catch {
        Write-Host "  [ERR] photo #$id : $_"
    }
}
