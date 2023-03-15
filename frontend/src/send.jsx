import { generateKey, encryptText } from './crypto'

(() => {
  if (!document.startViewTransition) {
    document.startViewTransition = (f) => f();
  }
  let handleEncrypt = async () => {
  	const plaintext = document.querySelector("#secret-field").value
  	const key = generateKey()
  	const ciphertext = encryptText(plaintext, key)
  	document.querySelector("#ciphertext-field").value = ciphertext
    const spinner = document.querySelector("#spinner").cloneNode(true)
    spinner.classList.remove("hidden")

    let viewTransition = document.startViewTransition(() => {
      document.querySelector("#encrypt").textContent = "Encrypting..."
      document.querySelector("#encrypt").prepend(spinner)
    })

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

  document.querySelector("#encrypt").onclick = handleEncrypt
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
