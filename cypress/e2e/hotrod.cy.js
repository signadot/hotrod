describe('hotrod spec', () => {
  it('request ride',
    () => {
      var frontendURL = 'http://frontend.' + Cypress.env('HOTROD_NAMESPACE') + ':8080';
      // inject routing key
      cy.intercept(frontendURL + '/*', (req) => {
        req.headers['baggage'] += ',sd-routing-key=' + Cypress.env('SIGNADOT_ROUTING_KEY');
      })

      cy.visit(frontendURL);
      cy.get('#requestRide').click();
      cy.get(':nth-child(1) > .col').click();
      cy.get(':nth-child(7) > .col > .text-success').click();
      cy.get(':nth-child(7) > .col > .text-success').contains(/Driver (.*) arriving in (.*)./);

      // check routing context
      var frontendSandboxName = Cypress.env('FRONTEND_SANDBOX_NAME');
      if (frontendSandboxName === "") {
        cy.get(':nth-child(3) > .col > .frontend').should('have.text', 'frontend (baseline)');
      } else {
        cy.get(':nth-child(3) > .col > .frontend').should('have.text', 'frontend (sandbox=' + frontendSandboxName + ')');
      }
      var locationSandboxName = Cypress.env('LOCATION_SANDBOX_NAME');
      if (locationSandboxName === "") {
        cy.get(':nth-child(4) > .col > .location').should('have.text', 'location (baseline)');
      } else {
        cy.get(':nth-child(4) > .col > .location').should('have.text', 'location (sandbox=' + locationSandboxName + ')');
      }
      var routeSandboxName = Cypress.env('ROUTE_SANDBOX_NAME');
      if (routeSandboxName === "") {
        cy.get(':nth-child(6) > .col > .route').should('have.text', 'route (baseline)');
      } else {
        cy.get(':nth-child(6) > .col > .route').should('have.text', 'route (sandbox=' + routeSandboxName + ')');
      }
      var driverSandboxName = Cypress.env('DRIVER_SANDBOX_NAME');
      if (driverSandboxName === "") {
        cy.get(':nth-child(7) > .col > .driver').should('have.text', 'driver (baseline)');
      } else {
        cy.get(':nth-child(7) > .col > .driver').should('have.text', 'driver (sandbox=' + driverSandboxName + ')');
      }
    });
})