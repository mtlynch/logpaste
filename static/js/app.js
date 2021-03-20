"use strict";

const baseUrl = document.location.toString().replace(/\/$/, "");

const curlCmd = document.getElementById("curl-cmd");
if (curlCmd) {
  curlCmd.innerHTML = Prism.highlight(
    `
echo "some text I want to upload" | \\
  curl \\
    --silent \\
    --show-error \\
    --form 'logpaste=<-' \\
    ${baseUrl}`.trim(),
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
