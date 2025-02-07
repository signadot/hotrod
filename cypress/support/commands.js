// Request a HotRod ride
Cypress.Commands.add('requestRide', (from, to) => {
    var frontendURL = 'http://frontend.' + Cypress.env('HOTROD_NAMESPACE') + ':8080/';
    // inject routing key
    cy.intercept('**', (req) => {
      req.headers['baggage'] += ',sd-routing-key=' + Cypress.env('SIGNADOT_ROUTING_KEY');
    })
    
    cy.visit(frontendURL);
    cy.get(':nth-child(1) > .chakra-select__wrapper > .chakra-select').select(from);
    cy.get(':nth-child(3) > .chakra-select__wrapper > .chakra-select').select(to);
    cy.contains('button', 'Request Ride').click();
    cy.contains('button', 'Show Logs').click();
})