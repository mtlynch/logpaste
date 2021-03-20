"use strict";

const baseUrl = document.location.toString().replace(/\/$/, "");

document.getElementById("curl-cmd").innerText = `
echo "some text I want to upload" | \\
  curl \\
    --silent \\
    --show-error \\
    --form 'logpaste=<-' \\
    ${baseUrl}`.trim();

document.getElementById("js-example").innerText = `
<script src="${baseUrl}/js/logpaste.js"></script>
<script>
const text = "some text I want to upload";

logpaste.uploadText(text).then((id) => {
  console.log(\`uploaded to \${baseUrl}/\${id}\`);
});
</script>
`.trim();

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
