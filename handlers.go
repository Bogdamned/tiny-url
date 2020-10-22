package main

import (
	"encoding/json"
	"log"
	"net/http"
)

//Handler handles therequests
type Handler struct {
	cache Cache
	urlDB UrlRepository
}

//NewHandler creates handler instance
func NewHandler(cache Cache, urlRepo UrlRepository) *Handler {
	return &Handler{cache, urlRepo}
}

//Tinify shortens an URL
func (h *Handler) Tinify(w http.ResponseWriter, r *http.Request) {
	tiny := ""

	payload := struct {
		URL string `json:"url"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Println("Unable to decode tiny url")
		http.Error(w, "Unable to decode tiny url", http.StatusInternalServerError)
		return
	}
	url := payload.URL
	urlValid := urlIsValid(url)

	if urlValid {
		newId, err := h.urlDB.GetID()
		if err != nil {
			log.Println("DB error: ", err)
			http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		newId++
		tiny = Encode(uint64(newId))
		//Insert into DB
		if err := h.urlDB.Insert(newId, url); err != nil {
			log.Println("DB error: ", err)
			http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		//Update counter
		h.urlDB.SetID(newId)

		//Insert into cache
		if err := h.cache.Insert(newId, url); err != nil {
			log.Println("Cache error: ", err)
			http.Error(w, "Cache error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		log.Println("Url isnot valid")
		http.Error(w, "Url is not valid. Please enter a valid one.", http.StatusInternalServerError)
		return
	}

	tinyURL := struct {
		TinyURL string `json:"tinyURL"`
	}{tiny}

	js, err := json.Marshal(tinyURL)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error. "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// TinyRedirect redirects to an originalURL
func (h *Handler) TinyRedirect(w http.ResponseWriter, r *http.Request) {
	tinyURL := r.URL.String()[1:]
	idDecoded, err := Decode(tinyURL)
	if err != nil {
		http.Error(w, "Error: unknown url", 404)
	}

	//Get from cache
	if url, err := h.cache.Get(int(idDecoded)); err == nil {
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}

	// If not present in cache get from DB
	if url, err := h.urlDB.Get(int(idDecoded)); err == nil {
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	} else {
		log.Println(err)
		http.Error(w, "Unknown url error", http.StatusInternalServerError)
		return
	}
}
