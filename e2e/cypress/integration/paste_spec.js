it("pastes an entry", () => {
  cy.visit("/");

  cy.get("#upload-textarea").type("test upload data");
  cy.get("#upload").click();
  cy.get("#result").click();
});
