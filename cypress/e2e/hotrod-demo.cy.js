describe('hotrod demo spec', () => {
  it('request ride and check deliver',
    () => {
      // request a ride
      cy.requestRide('1', '123')

      // check if the driver has been delivered
      cy.get(':nth-child(6) > .css-bjcoli').contains(/Driver (.*) arriving in (.*)./);
    });
})