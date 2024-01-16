const request = require('request');
const assert = require('assert');

const HOST = process.env.ENDPOINT;
const SIGNADOT_API_KEY = process.env.SIGNADOT_API_KEY;

// Set the default headers for all requests.
if (SIGNADOT_API_KEY) {
    console.log("Setting Signadot API key");
    req = request.defaults({
        headers: {
            'signadot-api-key': SIGNADOT_API_KEY,
        }
    });
} else {
    req = request;
}

describe("Route Service", function() {
    it("should return a 200 status code and response body matches the expected JSON structure for /route", function(done) {
        // Make a request to the /route endpoint with query parameters
        req.get({
            url: `${HOST}/route`,
            qs: {
                pickup: '123',
                dropoff: '456'
            }
        }, function(error, response, body) {
            // Assert that the status code is 200
            assert.equal(response.statusCode, 200);

            // Parse the response body as JSON
            const responseBody = JSON.parse(body);

            // Assert that the response body has the expected properties and data types
            assert.ok(responseBody.hasOwnProperty('Pickup'));
            assert.ok(responseBody.hasOwnProperty('Dropoff'));
            assert.ok(responseBody.hasOwnProperty('ETA'));
            assert.strictEqual(typeof responseBody.Pickup, 'string');
            assert.strictEqual(typeof responseBody.Dropoff, 'string');
            assert.strictEqual(typeof responseBody.ETA, 'number');

            done();
        });
    });

    it("should return a 400 status code when missing required 'pickup' parameter' message for /route", function(done) {
        // Make a request to the /route endpoint
        req.get(`${HOST}/route`, function(error, response, body) {
            // Assert that the status code is 400
            assert.equal(response.statusCode, 400);

            // Assert that the response body contains the expected message
            assert.equal(body, "Missing required 'pickup' parameter\n");

            done();
        });
    });

    it("should return a 400 status code with 'Missing required 'dropoff' parameter' message for /route", function(done) {
        // Make a request to the /route endpoint with query parameter
        req.get(`${HOST}/route?pickup=123`, function(error, response, body) {
            // Assert that the status code is 400
            assert.equal(response.statusCode, 400);

            // Assert that the response body contains the expected message
            assert.equal(body, "Missing required 'dropoff' parameter\n");

            done();
        });
    });

    it("should return a non-negative ETA for /route", function(done) {
     // Make a request to the /route endpoint with query parameters
         req.get({
             url: `${HOST}/route`,
             qs: {
                 pickup: '123',
                 dropoff: '456'
             }
         }, function(error, response, body) {
             const data = JSON.parse(body);
             assert.ok(data.ETA >= 0);
             done();
         });
     });
});
