package handlers

import (
	"net/http"
	"pr-review/internal/http/dto"

	"github.com/go-chi/render"
)

func (h *Handlers) CreatePR() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, dto.ErrContentTypeNotJson)
			return
		}

		var req dto.CreatePRRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, dto.ErrInvalidBody)
			return
		}
		if req.Validate() != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, req.Validate())
			return
		}

		err := h.uc.CreatePR(r.Context(), &req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, dto.ErrInternal)
			return
		}

		w.WriteHeader(http.StatusCreated)

	}
}

func (h *Handlers) MergePR() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		panic("not implemented")
	}
}

func (h *Handlers) ReassignPR() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		panic("not implemented")
	}
}
