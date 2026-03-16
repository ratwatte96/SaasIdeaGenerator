package server

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"saasideagenerator/backend/internal/store"
)

type Server struct{ store *store.Store }

func New(db *sql.DB) *Server { return &Server{store: store.New(db)} }

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", s.health)
	mux.HandleFunc("/api/ideas/", s.ideaByID)
	mux.HandleFunc("/api/ideas", s.listIdeas)
	mux.HandleFunc("/api/products", s.listProducts)
	return mux
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, 200, map[string]string{"status": "ok"})
}

func (s *Server) listIdeas(w http.ResponseWriter, r *http.Request) {
	minDemand, _ := strconv.ParseFloat(r.URL.Query().Get("min_demand_score"), 64)
	limit, offset := store.ParseLimitOffset(r.URL.Query().Get("limit"), r.URL.Query().Get("offset"))
	ideas, err := s.store.ListIdeas(r.Context(), r.URL.Query().Get("category"), r.URL.Query().Get("competition_level"), minDemand, limit, offset)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]any{"data": ideas, "limit": limit, "offset": offset})
}

func (s *Server) ideaByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/ideas/")
	if id == "" || strings.Contains(id, "/") {
		writeJSON(w, 400, map[string]string{"error": "invalid idea id"})
		return
	}
	idea, err := s.store.GetIdea(r.Context(), id)
	if err == sql.ErrNoRows {
		writeJSON(w, 404, map[string]string{"error": "idea not found"})
		return
	}
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	related, err := s.store.RelatedIdeas(r.Context(), idea.SourceProductID, idea.ID, 5)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]any{"idea": idea, "related_ideas": related})
}

func (s *Server) listProducts(w http.ResponseWriter, r *http.Request) {
	limit, offset := store.ParseLimitOffset(r.URL.Query().Get("limit"), r.URL.Query().Get("offset"))
	rows, err := s.store.ListProducts(r.Context(), r.URL.Query().Get("category"), limit, offset)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]any{"data": rows, "limit": limit, "offset": offset})
}
