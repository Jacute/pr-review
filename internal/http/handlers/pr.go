package handlers

import "net/http"

func (h *Handlers) CreatePR() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		panic("not implemented")
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
