package handler

import (
	"encoding/json"
	"net/http"
	"testLionParcell/usecase"
)

type Handlers interface {
	Check(w http.ResponseWriter, r *http.Request)
	Upload(w http.ResponseWriter, r *http.Request)
}

type handlers struct {
	usecase usecase.Usecases
}

func NewHandler(usecase usecase.Usecases) Handlers {
	return &handlers{
		usecase: usecase,
	}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (h *handlers) Check(w http.ResponseWriter, r *http.Request) {
	Message := "Chuck"
	healthStatus := struct {
		Message string `json:"Message"`
	}{Message}
	respondWithJSON(w, http.StatusOK, healthStatus)
}

func (h *handlers) Upload(w http.ResponseWriter, r *http.Request) {
	_, file, err := r.FormFile("upload-file")
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, "fail to upload")
	}

	err = h.usecase.Upload(*file)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, err)
	}

	respondWithJSON(w, http.StatusOK, "success")
}
