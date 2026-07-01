# 摄影作品测试数据生成器 (PowerShell 版本)
# 模拟摄影师上传图片，然后管理员审核，最后投票

$ErrorActionPreference = "Stop"
$BaseUrl = "http://localhost:8080"

# 1. 登录摄影师
$loginBody = '{"username":"photographer","password":"photo123"}'
$loginResp = Invoke-RestMethod -Uri "$BaseUrl/api/user/login" -Method POST -Body $loginBody -ContentType "application/json"
$photoToken = $loginResp.data.token
Write-Host "[OK] 摄影师登录"

# 2. 登录管理员
$adminBody = '{"username":"admin","password":"admin123"}'
$adminResp = Invoke-RestMethod -Uri "$BaseUrl/api/user/login" -Method POST -Body $adminBody -ContentType "application/json"
$adminToken = $adminResp.data.token
Write-Host "[OK] 管理员登录"

# 3. 登录投票用户
$voterBody = '{"username":"voter","password":"vote123"}'
$voterResp = Invoke-RestMethod -Uri "$BaseUrl/api/user/login" -Method POST -Body $voterBody -ContentType "application/json"
$voterToken = $voterResp.data.token
Write-Host "[OK] 投票用户登录"

# 4. 准备测试图片 - 用 PIL 风格的程序生成 (用 Go 工具)
# 这里直接用 Go 工具生成
Write-Host "[RUN] 生成测试图片并上传..."

# 检查 seed-tool.exe 是否存在
$seedTool = "D:\GitHub\wwmm\wwmm-system\backend\seed-tool.exe"
if (Test-Path $seedTool) {
    Set-Location "D:\GitHub\wwmm\wwmm-system\backend"
    & $seedTool 2>&1 | ForEach-Object { Write-Host "  $_" }
} else {
    Write-Host "[WARN] seed-tool.exe 不存在"
}

Start-Sleep -Seconds 2

# 5. 审核所有待审核的作品
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

# 6. 投票
Write-Host "[RUN] 模拟投票..."
$approved = Invoke-RestMethod -Uri "$BaseUrl/api/photo/list?size=20"
foreach ($p in @($approved.data.list)) {
    if ($null -eq $p.photoId) { continue }
    try {
        $r = Invoke-RestMethod -Uri "$BaseUrl/api/photo/$($p.photoId)/vote" -Method POST -Headers @{Authorization="Bearer $voterToken"}
        Write-Host "  [VOTE] photo #$($p.photoId) '$($p.title)' - tx: $($r.data.txHash.Substring(0,16))..."
    } catch {
        Write-Host "  [SKIP] photo #$($p.photoId): $($_.Exception.Message)"
    }
}

# 7. 用 admin 也投一次 (让 admin 不是作者的用户都能投)
$adminPending = Invoke-RestMethod -Uri "$BaseUrl/api/photo/list?size=20"
foreach ($p in @($adminPending.data.list)) {
    if ($null -eq $p.photoId) { continue }
    if ($p.photographerId -ne 100) {
        try {
            $r = Invoke-RestMethod -Uri "$BaseUrl/api/photo/$($p.photoId)/vote" -Method POST -Headers @{Authorization="Bearer $adminToken"}
            Write-Host "  [ADMIN VOTE] photo #$($p.photoId) '$($p.title)'"
        } catch {
            # 忽略
        }
    }
}

# 8. 显示最终状态
Write-Host "`n=== 最终状态 ==="
$state = Invoke-RestMethod -Uri "$BaseUrl/api/chain/state"
Write-Host "区块高度: $($state.data.latestIndex)"
Write-Host "总交易数: $($state.data.totalTxs)"
Write-Host "作品存证: $($state.data.txCertifyCount)"
Write-Host "投票记录: $($state.data.txVoteCount)"

$final = Invoke-RestMethod -Uri "$BaseUrl/api/photo/list?size=20"
Write-Host "`n=== 作品列表 ==="
foreach ($p in @($final.data.list)) {
    Write-Host "  #$($p.photoId) $($p.title) - 票数=$($p.voteCount) 链=$($p.isOnChain)"
}

Write-Host "`n[完成]"
