"use strict";

// Make ESLint happy.
/* global Prism, logpaste */

const baseUrl = document.location.origin;

const curlCmd = document.getElementById("curl-cmd");
if (curlCmd) {
  curlCmd.innerHTML = Prism.highlight(
    `
echo "some text I want to upload" | \\
  curl -F '_=<-' ${baseUrl}`.trim(),
    Prism.languages.bash,
    "bash"
  );
}

const curlFileCmd = document.getElementById("curl-file-cmd");
if (curlFileCmd) {
  curlFileCmd.innerHTML = Prism.highlight(
    `
curl -F "_=@/path/to/file.txt" ${baseUrl}`.trim(),
    Prism.languages.bash,
    "bash"
  );
}

const jsExample = document.getElementById("js-example");
if (jsExample) {
  jsExample.innerHTML = Prism.highlight(
    `
<script src="${baseUrl}/js/logpaste.js"></script>
<script>
const text = "some text I want to upload";

logpaste.uploadText(text).then((id) => {
  console.log(\`uploaded to ${baseUrl}/\${id}\`);
});
</script>
    `.trim(),
    Prism.languages.javascript,
    "javascript"
  );
}

function displayResult(resultId) {
  clearError();
  clearResult();

  const resultDiv = document.getElementById("result");

  const resultUrl = `${document.location}${resultId}`;

  const header = document.createElement("h3");
  header.innerText = "Shareable link";
  resultDiv.appendChild(header);

  const anchor = document.createElement("a");
  anchor.href = `/${resultId}`;
  anchor.innerText = resultUrl;
  resultDiv.appendChild(anchor);

  resultDiv.style.visibility = "visible";
}

function clearResult() {
  const resultDiv = document.getElementById("result");
  while (resultDiv.firstChild) {
    resultDiv.removeChild(resultDiv.lastChild);
  }
  resultDiv.style.visibility = "hidden";
}

function clearError() {
  const uploadError = document.getElementById("form-upload-error");
  uploadError.innerText = " ";
  uploadError.style.visibility = "hidden";
}

function displayError(error) {
  const uploadError = document.getElementById("form-upload-error");
  uploadError.innerText = error;
  uploadError.style.visibility = "visible";
}

document.getElementById("upload").addEventListener("click", () => {
  const textToUpload = document.getElementById("upload-textarea").value;
  logpaste
    .uploadText(textToUpload)
    .then((id) => {
      displayResult(id);
    })
    .catch((error) => {
      clearResult();
      displayError(error);
    });
});
