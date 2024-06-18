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
      var sandboxName = Cypress.env('SIGNADOT_SANDBOX_NAME');

      cy.get('.css-8g8ihq :nth-child(2) > :nth-child(2)').should('have.text', 'frontend');
      if (Cypress.env('SANDBOXED_FRONTEND') != "1") {
        cy.get(':nth-child(2) > :nth-child(3)').should('have.text', '(baseline)');
      } else {
        cy.get(':nth-child(2) > :nth-child(3)').should('have.text', '(' + sandboxName + ')');
      }
      cy.get('.css-8g8ihq :nth-child(3) > :nth-child(2)').should('have.text', 'location');
      if (Cypress.env('SANDBOXED_LOCATION') != "1") {
        cy.get(':nth-child(3) > :nth-child(3)').should('have.text', '(baseline)');
      } else {
        cy.get(':nth-child(3) > :nth-child(3)').should('have.text', '(' + sandboxName + ')');
      }
      cy.get(':nth-child(5) > :nth-child(2)').should('have.text', 'route');
      if (Cypress.env('SANDBOXED_ROUTE') != "1") {
        cy.get(':nth-child(5) > :nth-child(3)').should('have.text', '(baseline)');
      } else {
        cy.get(':nth-child(5) > :nth-child(3)').should('have.text', '(' + sandboxName + ')');
      }
      cy.get(':nth-child(6) > :nth-child(2)').should('have.text', 'driver');
      if (Cypress.env('SANDBOXED_DRIVER') != "1") {
        cy.get(':nth-child(6) > :nth-child(3)').should('have.text', '(baseline)');
      } else {
        cy.get(':nth-child(6) > :nth-child(3)').should('have.text', '(' + sandboxName + ')');
      }
    });
})