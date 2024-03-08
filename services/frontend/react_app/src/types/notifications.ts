export type NotificationResponse = {
    cursor: number,
    notifications: Notification[],
}

export type Notification = {
    id: string,
    timestamp: string,
    context: {
        request: {
            id: number,
            sessionID: number,
            pickupLocationID: number,
            dropoffLocationID: number,
        },
        routingKey: string,
        baselineWorkload: string,
        sandboxName: string,
    },
    body: string,
}

// {
//     "cursor": 3,
//     "notifications": [
//     {
//         "id": "req-1-frontend-dispatching-driver",
//         "timestamp": "2024-03-07T16:13:28.141199-05:00",
//         "context": {
//             "request": {
//                 "id": 1,
//                 "sessionID": 9402,
//                 "pickupLocationID": 1,
//                 "dropoffLocationID": 1
//             },
//             "routingKey": "",
//             "baselineWorkload": "",
//             "sandboxName": ""
//         },
//         "body": "Processing dispatch driver request"
//     },
//     {
//         "id": "req-1-location-resolve",
//         "timestamp": "2024-03-07T21:13:29.201590298Z",
//         "context": {
//             "request": {
//                 "id": 1,
//                 "sessionID": 9402,
//                 "pickupLocationID": 1,
//                 "dropoffLocationID": 1
//             },
//             "routingKey": "",
//             "baselineWorkload": "location",
//             "sandboxName": ""
//         },
//         "body": "Resolving locations"
//     },
//     {
//         "id": "req-1-finding-driver",
//         "timestamp": "2024-03-07T21:13:29.834238282Z",
//         "context": {
//             "request": {
//                 "id": 1,
//                 "sessionID": 9402,
//                 "pickupLocationID": 1,
//                 "dropoffLocationID": 1
//             },
//             "routingKey": "",
//             "baselineWorkload": "driver",
//             "sandboxName": ""
//         },
//         "body": "Finding an available driver"
//     },
//     {
//         "id": "req-1-route-resolve",
//         "timestamp": "2024-03-07T21:13:30.051313207Z",
//         "context": {
//             "request": {
//                 "id": 1,
//                 "sessionID": 9402,
//                 "pickupLocationID": 1,
//                 "dropoffLocationID": 1
//             },
//             "routingKey": "",
//             "baselineWorkload": "route",
//             "sandboxName": ""
//         },
//         "body": "Resolving routes"
//     }
// ]
// }