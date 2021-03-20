it("pastes an entry", () => {
  cy.visit("/");

  cy.get("#upload-textarea").type("test upload data");
  cy.get("#upload").click();
  cy.get("#result a").should("not.be.empty");
  cy.get("#result a")
    .invoke("attr", "href")
    .then((href) => {
      cy.log(href);
      cy.request(href).its("body").should("equal", "test upload data");
    });
});
