$ErrorActionPreference = "Continue"
$BaseUrl = "http://localhost:8080"
$adminBody = '{"username":"admin","password":"admin123"}'
$adminResp = Invoke-RestMethod -Uri "$BaseUrl/api/user/login" -Method POST -Body $adminBody -ContentType "application/json"
$adminToken = $adminResp.data.token

$voterBody = '{"username":"voter","password":"vote123"}'
$voterResp = Invoke-RestMethod -Uri "$BaseUrl/api/user/login" -Method POST -Body $voterBody -ContentType "application/json"
$voterToken = $voterResp.data.token

# 审核
Write-Host "[RUN] 审核所有待审核作品..."
$pending = Invoke-RestMethod -Uri "$BaseUrl/api/photo/pending" -Headers @{Authorization="Bearer $adminToken"}
foreach ($p in @($pending.data.list)) {
    if ($null -eq $p.photoId) { continue }
    $body = '{"approve":true,"comment":"内容符合规范，审核通过"}'
    try {
        $r = Invoke-RestMethod -Uri "$BaseUrl/api/photo/$($p.photoId)/audit" -Method POST -Body $body -ContentType "application/json" -Headers @{Authorization="Bearer $adminToken"}
        Write-Host "  [OK] photo #$($p.photoId) '$($p.title)' 审核通过"
    } catch {
        Write-Host "  [ERR] photo #$($p.photoId): $_"
    }
}

# 投票 (voter 用户)
Write-Host "[RUN] voter 投票..."
for ($i = 1; $i -le 8; $i++) {
    try {
        $r = Invoke-RestMethod -Uri "$BaseUrl/api/photo/$i/vote" -Method POST -Headers @{Authorization="Bearer $voterToken"}
        Write-Host "  [VOTE] photo #$($i) - tx: $($r.data.txHash.Substring(0,20))..."
    } catch {
        Write-Host "  [SKIP] photo #$($i): $($_.Exception.Message)"
    }
}

# 投票 (admin 用户)
Write-Host "[RUN] admin 投票..."
for ($i = 1; $i -le 8; $i++) {
    try {
        $r = Invoke-RestMethod -Uri "$BaseUrl/api/photo/$i/vote" -Method POST -Headers @{Authorization="Bearer $adminToken"}
        Write-Host "  [VOTE] photo #$($i)"
    } catch {
        # ignore
    }
}

# 显示状态
Write-Host "`n=== 最终状态 ==="
$state = Invoke-RestMethod -Uri "$BaseUrl/api/chain/state"
$state.data | Format-List

$final = Invoke-RestMethod -Uri "$BaseUrl/api/photo/list?size=20"
Write-Host "`n=== 作品列表 ==="
foreach ($p in @($final.data.list)) {
    Write-Host "  #$($p.photoId) $($p.title) - 票数=$($p.voteCount) 链=$($p.isOnChain)"
}
