import { decrypt } from './crypto'

(() => {
  if (window.CIPHERTEXT !== undefined && window.location.hash !== "") {
    let key = window.location.hash.substr(1)
    let decrypted = decrypt(CIPHERTEXT, key)
    document.querySelector("#message").textContent = decrypted
  }
})()
