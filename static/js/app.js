"use strict";

document.getElementById("curl-cmd").innerText = `
      echo "some data I want to upload" | \\
        curl \\
          --silent \\
          --show-error \\
          --form 'logpaste=<-' \\
           ${document.location}`.trim();

document.getElementById("upload").addEventListener("click", (evt) => {
  const textToUpload = document.getElementById("upload-textarea").value;
  logpaste.uploadText(textToUpload).then((id) => {
    const resultUrl = `${document.location}${id}`;

    const paragraph = document.createElement("p");
    paragraph.innerText = resultUrl;

    const anchor = document.createElement("a");
    anchor.href = resultUrl;
    anchor.appendChild(paragraph);

    const resultDiv = document.getElementById("result");
    while (resultDiv.firstChild) {
      resultDiv.removeChild(resultDiv.lastChild);
    }
    resultDiv.appendChild(anchor);
  });
});
