package ingestion

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/rphmauriciodev/goflow/internal/processing"
)

type BatchRequest struct {
	Source  string `json:"source"`
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

type Handler struct {
	repo         processing.BatchRepository
	orchestrator *processing.Orchestrator
}

func NewHandler(r processing.BatchRepository, o *processing.Orchestrator) *Handler {
	return &Handler{
		repo:         r,
		orchestrator: o,
	}
}

func (r *BatchRequest) Validate() error {
	if strings.TrimSpace(r.Source) == "" {
		return errors.New("o campo 'source' é obrigatório")
	}

	validTypes := map[string]bool{"json": true, "csv": true, "xml": true}
	if !validTypes[strings.ToLower(r.Type)] {
		return fmt.Errorf("tipo de arquivo '%s' inválido. Use: json, csv ou xml", r.Type)
	}
	if strings.TrimSpace(r.Payload) == "" {
		return errors.New("o campo 'payload' não pode estar vazio")
	}
	return nil
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/ping", h.Ping)
	mux.HandleFunc("POST /batches", h.UploadBatch)
}

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "pong (vindo do módulo ingestion)\n")
}

func (h *Handler) UploadBatch(w http.ResponseWriter, r *http.Request) {
	var req BatchRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "JSON Inválido", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	batch := &processing.Batch{
		ID:        fmt.Sprintf("batch_%d", time.Now().Unix()),
		Source:    req.Source,
		Type:      req.Type,
		RawData:   []byte(req.Payload),
		CreatedAt: time.Now(),
	}

	if err := h.repo.Save(batch); err != nil {
		http.Error(w, "Erro ao salvar lote", http.StatusInternalServerError)
		return
	}
	go h.orchestrator.Start(batch, h.repo)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"batch_id": batch.ID,
		"status":   "accepted",
		"message":  "Processamento iniciado em background",
	})
}
