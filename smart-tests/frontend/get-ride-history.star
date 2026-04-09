ck = smart_test.check("ride-history-status")

payload = {
    "sessionID": 998877,
    "requestID": 445566,
    "pickupLocationID": 1,
    "dropoffLocationID": 731,
}

dispatch = http.post(
    url="http://frontend.hotrod-devmesh.svc:8080/dispatch",
    json_body=payload,
    capture=True,
    name="dispatchForRideHistory",
)
if dispatch.status_code != 200:
    ck.error("dispatch request failed with status {}", dispatch.status_code)

history = http.get(
    url="http://frontend.hotrod-devmesh.svc:8080/ride-history?sessionID=998877",
    capture=True,
    name="getRideHistory",
)

if history.status_code != 200:
    ck.error("ride history request failed with status {}", history.status_code)
else:
    body = history.json()
    if type(body) != "dict":
        ck.error("unexpected response type: {}", type(body))
    else:
        if not body.has("totalCount"):
            ck.error("missing totalCount field in response")
        if not body.has("entries"):
            ck.error("missing entries field in response")

        entries = body["entries"]
        if type(entries) != "list":
            ck.error("entries is not a list, got {}", type(entries))
        elif len(entries) == 0:
            ck.error("entries list is empty")
        else:
            expected_request_id = 445566
            matching = None
            for e in entries:
                if type(e) != "dict":
                    continue
                if e.get("requestID") == expected_request_id:
                    matching = e
                    break

            if matching == None:
                ck.error("no ride history entry found for requestID {}", expected_request_id)
            else:
                if matching.get("pickupLocation") == "":
                    ck.error("pickupLocation is empty")
                if matching.get("dropoffLocation") == "":
                    ck.error("dropoffLocation is empty")
                if matching.get("requestedAt") == "":
                    ck.error("requestedAt is empty")
                if not matching.has("driverPlate"):
                    ck.error("driverPlate field missing")
