package handlers

import (
	"encoding/json"
	"net/http"
)

type RespuestaError struct {
	Estado  int    `json:"estado"`
	Error   string `json:"error"`
	Mensaje string `json:"mensaje"`
}

type RespuestaExito struct {
	Estado  int         `json:"estado"`
	Mensaje string      `json:"mensaje,omitempty"`
	Datos   interface{} `json:"datos,omitempty"`
	Total   *int        `json:"total,omitempty"`
}

func escribirJSON(w http.ResponseWriter, codigoEstado int, contenido interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(codigoEstado)
	json.NewEncoder(w).Encode(contenido)
}

func escribirError(w http.ResponseWriter, codigoEstado int, tipoError, detalle string) {
	escribirJSON(w, codigoEstado, RespuestaError{
		Estado:  codigoEstado,
		Error:   tipoError,
		Mensaje: detalle,
	})
}
