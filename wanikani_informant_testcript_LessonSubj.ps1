$token = "235d3b90-1eb4-4e06-962e-fc7ad57b09e8"
$h = @{
    'Authorization' = "Bearer $token"
}

$u = "https://api.wanikani.com/v2/summary"
$content = Invoke-restmethod $u -Headers $h

<#
$content.data.reviews
Data:
available_at                subject_ids               
------------                -----------               
2022-06-02T08:00:00.000000Z {}                        
2022-06-02T09:00:00.000000Z {539, 2728, 2674, 2651...}
2022-06-02T10:00:00.000000Z {2497, 471, 468}          
2022-06-02T11:00:00.000000Z {2689, 579, 580, 2690...} 
2022-06-02T12:00:00.000000Z {}                        
2022-06-02T13:00:00.000000Z {2594}                    
2022-06-02T14:00:00.000000Z {2653, 109, 856, 2495...} 
2022-06-02T15:00:00.000000Z {2631}                    
2022-06-02T16:00:00.000000Z {}                    
...

Proposal:
- Cache subjects, but acquire frequently
 > Target: GET https://api.wanikani.com/v2/subjects/2631
 > Also see: https://docs.api.wanikani.com/20170710/#conditional-requests
- Cross-reference subject with subject_ids to determine graduation chance

#>