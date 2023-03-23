import { decrypt } from './crypto'

(() => {
  const fileDisplay = document.getElementById('file_display');

  if (window.CIPHERTEXT !== undefined && window.location.hash !== "") {
    let key = window.location.hash.substr(1)
    let decrypted = JSON.parse(decrypt(CIPHERTEXT, key))
    document.querySelector("#message").textContent = decrypted.text
    console.log(decrypted.files)

    for(const fileId in decrypted.files) {
      const file = decrypted.files[fileId]
      fileDisplay.classList.remove("hidden")
      const fileElement = document.querySelector(".file_download_tmpl").cloneNode(true);
      fileElement.classList.remove("hidden", "file_download_tmpl")

      const fileInfoElement = document.createElement("div");
      fileInfoElement.classList.add("flex", "items-center");

      const fileNameElement = document.createElement("div");
      fileNameElement.innerText = file.name;
      fileNameElement.classList.add("font-medium", "text-sm", "mr-2");

      const fileSizeElement = document.createElement("div");
      fileSizeElement.classList.add("text-gray-500", "text-sm", "mr-2");
      fileSizeElement.innerText = `${(file.size / 1024).toFixed(2)} KB`;

      const fileDownloadElement = document.createElement("a");
      fileDownloadElement.classList.add("ml-auto", "px-2", "py-1", "bg-orange-500", "text-white", "rounded");
      fileDownloadElement.innerText = "Download";
      fileDownloadElement.href = file.data;
      fileDownloadElement.download = file.name;

      // Add the UI elements to the file preview
      fileInfoElement.appendChild(fileNameElement);
      fileInfoElement.appendChild(fileSizeElement);
      fileInfoElement.appendChild(fileDownloadElement);

      fileElement.appendChild(fileInfoElement);
      fileDisplay.appendChild(fileElement);
    }
  }
})()
