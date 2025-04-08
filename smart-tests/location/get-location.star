res = http.get(
    url="http://location.hotrod-istio.svc:8081/locations", # can also be cluster internal URLs like http://name.namespace.svc
    capture=True, # enables SmartDiff
    name="getLocations"
)

ck = smart_test.check("get-location-status")
if res.status_code != 200:
    ck.error("bad status code: {}", res.status_code)

print(res.status_code)
print(res.body())
