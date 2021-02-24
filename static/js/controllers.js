"use strict";

(function (windows) {
  function uploadText(text) {
    return fetch("/", {
      method: "PUT",
      body: text,
    })
      .then((response) => {
        const contentType = response.headers.get("content-type");
        const isJson =
          contentType && contentType.indexOf("application/json") !== -1;
        // Success case is an HTTP 200 response and a JSON body.
        if (response.status === 200 && isJson) {
          return Promise.resolve(response.json());
        }
      })
      .then((data) => data.id);
  }
  if (!window.hasOwnProperty("controllers")) {
    window.controllers = {};
  }
  window.controllers.uploadText = uploadText;
})(window);
