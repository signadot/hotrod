

export const apiGet = async <R>(path: string) => {
    const d = await fetch(path)
    const json = await d.json()
    return json as R
}

export const apiPost = async <R>(path: string, data: Record<string, string> | object, headers: Record<string, string>) => {
    const d = await fetch(path, {
        headers: {
            ...headers,
            "Content-Type": "application/json",
        },
        method: "POST",
        body: JSON.stringify(data),
    })
    const json = await d.json()
    return json as R
}

//
// const headers = {
//     'baggage': 'session=' + sessionID + ', request=' + lastRequestID
// };
// console.log('sending headers', headers);
//
// var pickupLocationID = parseInt($('#pickupLocation').val());
// var dropoffLocationID = parseInt($('#dropoffLocation').val());
//


// logSuccess(new Date(), 'Requesting a ride', {
//     request: {
//         id: lastRequestID,
//         pickupLocationID: pickupLocationID,
//         dropoffLocationID: dropoffLocationID
//     },
// })
//
// $.ajax(pathPrefix + '/dispatch?nonse=' + Math.random(), {
//     headers: headers,
//     method: 'POST',
//     contentType : 'application/json',
//     data: JSON.stringify({
//         sessionID: sessionID,
//         requestID: lastRequestID,
//         pickupLocationID: pickupLocationID,
//         dropoffLocationID: dropoffLocationID
//     }),
//     success: function (data, textStatus) {
//     },
//     error: function (xhr, status, error) {
//         if (xhr.responseText) {
//             error += " (" + xhr.responseText + ")"
//         }
//         logError(new Date(), 'Error requesting a ride to frontend API', {
//             request: {
//                 id: lastRequestID,
//                 pickupLocationID: pickupLocationID,
//                 dropoffLocationID: dropoffLocationID
//             },
//             error: error,
//         })
//     }
// });