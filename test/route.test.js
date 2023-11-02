const request = require('request');
const assert = require('assert');

const host = "localhost:8083";

describe("Route Service", function() {
    it("should return a 200 status code and response body matches the expected JSON structure for /route", function(done) {
        // Make a request to the /route endpoint with query parameters
        request.get({
            url: 'http://localhost:8083/route',
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
        request.get('http://localhost:8083/route', function(error, response, body) {
            // Assert that the status code is 400
            assert.equal(response.statusCode, 400);

            // Assert that the response body contains the expected message
            assert.equal(body, "Missing required 'pickup' parameter\n");

            done();
        });
    });

    it("should return a 400 status code with 'Missing required 'dropoff' parameter' message for /route", function(done) {
        // Make a request to the /route endpoint with query parameter
        request.get('http://localhost:8083/route?pickup=123', function(error, response, body) {
            // Assert that the status code is 400
            assert.equal(response.statusCode, 400);

            // Assert that the response body contains the expected message
            assert.equal(body, "Missing required 'dropoff' parameter\n");

            done();
        });
    });
});