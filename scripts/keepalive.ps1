$IntervalSeconds = if ($env:KEEP_ALIVE_INTERVAL_SECONDS) { [int]$env:KEEP_ALIVE_INTERVAL_SECONDS } else { 600 }
$Path = if ($env:KEEP_ALIVE_PATH) { $env:KEEP_ALIVE_PATH } else { "/api/v1/device-management/health" }
$BaseUrl = if ($env:KEEP_ALIVE_URL) { $env:KEEP_ALIVE_URL } else { $env:RENDER_EXTERNAL_URL }

if ([string]::IsNullOrWhiteSpace($BaseUrl)) {
    Write-Host "keep-alive disabled: KEEP_ALIVE_URL or RENDER_EXTERNAL_URL is required"
    exit 0
}

$TargetUrl = $BaseUrl.TrimEnd("/") + $Path
Write-Host "keep-alive enabled: pinging $TargetUrl every $IntervalSeconds seconds"

while ($true) {
    Start-Sleep -Seconds $IntervalSeconds
    try {
        Invoke-WebRequest -Uri $TargetUrl -Method GET -UseBasicParsing | Out-Null
    }
    catch {
        Write-Host "keep-alive ping failed: $($_.Exception.Message)"
    }
}
