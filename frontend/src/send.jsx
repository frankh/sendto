import { generateKey, encryptText } from './crypto'

(() => {
  function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
  }

  let files = {}
  if (!document.startViewTransition) {
    document.startViewTransition = (f) => f();
  }
  let handleEncrypt = async () => {
  	const plaintext = document.querySelector("#secret-field").value
  	const key = generateKey()
    const spinner = document.querySelector("#spinner").cloneNode(true)
    spinner.classList.remove("hidden")

    let viewTransition = document.startViewTransition(() => {
      document.querySelector("#encrypt").textContent = "Encrypting..."
      document.querySelector("#encrypt").prepend(spinner)
    })
    await sleep(100)
    const ciphertext = await encryptText(JSON.stringify({text: plaintext, files: files}), key)

    const response = await fetch("/api/send", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        to: document.querySelector("#to-field").value,
        cipherText: ciphertext,
      }),
    })

    const message = await response.json()
    const shareUrl = "http://sendto/r/"+message.ID+"#"+key
    if (viewTransition) {
      await viewTransition.finished
    }
    document.startViewTransition(() => {
      document.querySelector("#encrypt").textContent = "Encrypt"
      document.querySelector("#to").textContent = document.querySelector("#to-field").value
      document.querySelector("#share-link").textContent = shareUrl
      document.querySelector("#share-link").href = shareUrl
      document.querySelector("#share-box").classList.remove("hidden")
      document.querySelector("#share-box").onclick = () => {
        const selection = window.getSelection()
        const range = document.createRange()
        range.selectNodeContents(document.querySelector("#share-link"))
        selection.removeAllRanges()
        selection.addRange(range)
      }
    })

    console.log(shareUrl)
  }

  const fileInput = document.getElementById('file_upload');
  const fileDisplay = document.getElementById('file_display');

  fileInput.addEventListener('change', function(e) {
    Array.from(e.target.files).forEach((f) => {
      const file = f;
      const reader = new FileReader();

      reader.onload = function(event) {
        const key = `file_${Date.now()}`;

        files[key] = {
          name: file.name,
          size: file.size,
          data: reader.result,
        }

        // Create a new UI element representing the file
        const fileElement = document.querySelector(".file_preview_tmpl").cloneNode(true);
        fileElement.classList.remove("hidden", "file_preview_tmpl")

        const fileInfoElement = document.createElement("div");
        fileInfoElement.classList.add("flex", "items-center");

        const fileNameElement = document.createElement("div");
        fileNameElement.innerText = file.name;
        fileNameElement.classList.add("font-medium", "text-sm", "mr-2");

        const fileSizeElement = document.createElement("div");
        fileSizeElement.classList.add("text-gray-500", "text-sm", "mr-2");
        fileSizeElement.innerText = `${(file.size / 1024).toFixed(2)} KB`;

        const fileRemoveElement = document.createElement("button");
        fileRemoveElement.classList.add("ml-auto", "px-2", "py-1", "bg-red-500", "text-white", "rounded");
        fileRemoveElement.innerText = "Remove";
        fileRemoveElement.addEventListener("click", function() {
          // Remove the file from local storage
          delete files[key];

          // Remove the UI element
          fileElement.remove();

          // Clear custom error (for when file size was too big)
          e.target.setCustomValidity("")
        });

        // Add the UI elements to the file preview
        fileInfoElement.appendChild(fileNameElement);
        fileInfoElement.appendChild(fileSizeElement);
        fileInfoElement.appendChild(fileRemoveElement);

        fileElement.appendChild(fileInfoElement);
        fileDisplay.appendChild(fileElement);
      };

      reader.readAsDataURL(file);
    })
    e.target.value = null
  });
  document.querySelector("#form").onsubmit = (e) => {
    e.preventDefault()
    if (JSON.stringify(files).length > 40*1024*1025) {
      e.target.querySelector("#file_upload").setCustomValidity("Files too big!")
      e.target.reportValidity()
      return
    }

    handleEncrypt()
  }
  // Focus the to field if it's not pre-set, otherwise secret field
  if (!document.querySelector("#to-field").value) {
    document.querySelector("#to-field").focus()
  } else {
    document.querySelector("#secret-field").focus()
  }

  document.querySelector("#secret-field").oninput = () => {
    document.startViewTransition(() => {
      document.querySelector("#share-box").classList.add("hidden")
    })
  }
})()
