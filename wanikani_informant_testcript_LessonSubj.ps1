$token = "235d3b90-1eb4-4e06-962e-fc7ad57b09e8"
$h = @{
    'Authorization' = "Bearer $token"
}

$u = "https://api.wanikani.com/v2/assignments?srs_stages=1,2,3,4"
$content = Invoke-restmethod $u -Headers $h

$u = "https://api.wanikani.com/v2/summary"
$detail = Invoke-restmethod $u -Headers $h

# Locate specific subject ID
$content.data | % {
    #Write-host $_.id -f yellow
    $_.data | % {
        [int]$t = 582
         
        if ($_.subject_id -eq $t) {
            Write-Host "-> Found $t!" -f green
            Write-Host "> Stage: $($_.srs_stage)"
            pause
        }
    }
}