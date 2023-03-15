package api

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/frankh/sendto/pkg/db"
)

type SendRequest struct {
	To         string `json:"to"`
	CipherText string `json:"cipherText"`
}

type SendResponse struct {
	ID string
}

func generateID() string {
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		panic("Failed to generate ID:" + err.Error())
	}

	return base64.URLEncoding.EncodeToString(b)
}

type SendHandler struct {
	Database db.DB
}

func (sh *SendHandler) Serve(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "HTTP Method Unsupported", http.StatusBadRequest)
		return
	}

	login, ok := r.Context().Value("login").(string)
	if !ok {
		http.Error(w, "Bad login provided", http.StatusBadRequest)
		return
	}

	var sr SendRequest
	err := json.NewDecoder(r.Body).Decode(&sr)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if sr.To == "" {
		http.Error(w, "Message must have \"To\" field", http.StatusBadRequest)
		return
	}

	message := db.Message{
		ID:         generateID(),
		From:       login,
		To:         sr.To,
		CipherText: sr.CipherText,
	}

	if err := sh.Database.Save(message); err != nil {
		http.Error(w, "Failed to save message", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}
