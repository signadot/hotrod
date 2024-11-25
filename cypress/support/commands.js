// Request a HotRod ride
Cypress.Commands.add('requestRide', (from, to) => {
    var frontendURL = 'http://frontend.hotrod-istio.svc:8080/';
    // inject routing key
    cy.intercept(frontendURL + '/*', (req) => {
      req.headers['baggage'] += ',sd-routing-key=' + Cypress.env('SIGNADOT_ROUTING_KEY');
    })
    
    cy.visit(frontendURL);
    cy.get(':nth-child(1) > .chakra-select__wrapper > .chakra-select').select(from);
    cy.get(':nth-child(3) > .chakra-select__wrapper > .chakra-select').select(to);
    cy.get('.chakra-button').click();
})