"use strict";

(function (windows) {
  function uploadText(text, baseUrl = "") {
    return fetch(baseUrl + "/", {
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
        // Treat any other response as an error.
        return response.text().then((text) => {
          if (text) {
            return Promise.reject(new Error(text));
          } else {
            return Promise.reject(new Error(response.statusText));
          }
        });
      })
      .then((data) => data.id);
  }
  if (!window.hasOwnProperty("logpaste")) {
    window.logpaste = {};
  }
  window.logpaste.uploadText = uploadText;
})(window);
