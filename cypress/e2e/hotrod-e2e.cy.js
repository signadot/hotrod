describe('hotrod e2e spec', () => {
  it('request ride and check routing context',
    () => {
      // request a ride
      cy.requestRide('1', '123')

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