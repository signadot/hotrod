describe('hotrod spec', () => {
  it('request ride',
    () => {
      var frontendURL = 'http://frontend.' + Cypress.env('HOTROD_NAMESPACE') + ':8080';
      // inject routing key
      cy.intercept(frontendURL + '/*', (req) => {
        req.headers['baggage'] += ',sd-routing-key=' + Cypress.env('SIGNADOT_ROUTING_KEY');
      })

      cy.visit(frontendURL);
      cy.get(':nth-child(1) > .chakra-select__wrapper > .chakra-select').select('1');
      cy.get(':nth-child(3) > .chakra-select__wrapper > .chakra-select').select('123');
      cy.get('.chakra-button').click();

      // check if the driver has been delivered
      cy.get(':nth-child(6) > .css-bjcoli').contains(/Driver (.*) arriving in (.*)./);

      if (Cypress.env('HOTROD_E2E') != "1"){
        return
      }

      // check routing context
      var frontendSandboxName = Cypress.env('FRONTEND_SANDBOX_NAME');
      cy.get('.css-8g8ihq :nth-child(2) > :nth-child(2)').should('have.text', 'frontend');
      if (frontendSandboxName === "") {
        cy.get(':nth-child(2) > :nth-child(3)').should('have.text', '(baseline)');
      } else {
        cy.get(':nth-child(2) > :nth-child(3)').should('have.text', '(' + frontendSandboxName + ')');
      }
      var locationSandboxName = Cypress.env('LOCATION_SANDBOX_NAME');
      cy.get('.css-8g8ihq :nth-child(3) > :nth-child(2)').should('have.text', 'location');
      if (locationSandboxName === "") {
        cy.get(':nth-child(3) > :nth-child(3)').should('have.text', '(baseline)');
      } else {
        cy.get(':nth-child(3) > :nth-child(3)').should('have.text', '(' + locationSandboxName + ')');
      }
      var routeSandboxName = Cypress.env('ROUTE_SANDBOX_NAME');
      cy.get(':nth-child(5) > :nth-child(2)').should('have.text', 'route');
      if (routeSandboxName === "") {
        cy.get(':nth-child(5) > :nth-child(3)').should('have.text', '(baseline)');
      } else {
        cy.get(':nth-child(5) > :nth-child(3)').should('have.text', '(' + routeSandboxName + ')');
      }
      var driverSandboxName = Cypress.env('DRIVER_SANDBOX_NAME');
      cy.get(':nth-child(6) > :nth-child(2)').should('have.text', 'driver');
      if (driverSandboxName === "") {
        cy.get(':nth-child(6) > :nth-child(3)').should('have.text', '(baseline)');
      } else {
        cy.get(':nth-child(6) > :nth-child(3)').should('have.text', '(' + driverSandboxName + ')');
      }
    });
})