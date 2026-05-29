wrk.method = "GET"
wrk.headers["Content-Type"] = "application/json"
wrk.headers["Authorization"] = "Bearer <YOUR_JWT_TOKEN>"
wrk.headers["XCSRF-Token"] = "fDDzcOgcJVk7eN8kwpMMJp7zL9PilFDSBk6kdaogyVI="
wrk.headers["Cookie"] = "csrf_token=fDDzcOgcJVk7eN8kwpMMJp7zL9PilFDSBk6kdaogyVI="

request = function()
    return wrk.format(nil, "/api/posters/by-alias/studio-tverskaya", nil, nil)
end