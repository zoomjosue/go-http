package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"go-http/models"
	"go-http/storage"
)

func RutaDragones(w http.ResponseWriter, r *http.Request) {
	ruta := strings.TrimPrefix(r.URL.Path, "/api/dragones")
	ruta = strings.TrimSuffix(ruta, "/")

	if ruta == "" {
		switch r.Method {
		case http.MethodGet:
			obtenerDragones(w, r)
		case http.MethodPost:
			crearDragon(w, r)
		default:
			escribirError(w, http.StatusMethodNotAllowed, "Método no permitido", "Use GET o POST en /api/dragones")
		}
		return
	}

	idStr := strings.TrimPrefix(ruta, "/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		escribirError(w, http.StatusBadRequest, "Solicitud incorrecta", "El ID debe ser un número entero")
		return
	}

	switch r.Method {
	case http.MethodGet:
		obtenerDragonPorRuta(w, r, id)
	case http.MethodPut:
		actualizarDragon(w, r, id)
	case http.MethodPatch:
		actualizarDragonParcial(w, r, id)
	case http.MethodDelete:
		eliminarDragon(w, r, id)
	default:
		escribirError(w, http.StatusMethodNotAllowed, "Método no permitido", "Use GET, PUT, PATCH o DELETE en /api/dragones/{id}")
	}
}

func obtenerDragones(w http.ResponseWriter, r *http.Request) {
	dragones, err := storage.Cargar()
	if err != nil {
		escribirError(w, http.StatusInternalServerError, "Error", "No se pudieron cargar los dragones")
		return
	}

	q := r.URL.Query()

	if idStr := q.Get("id"); idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			escribirError(w, http.StatusBadRequest, "Solicitud incorrecta", "El parámetro 'id' debe ser un número entero")
			return
		}
		for _, d := range dragones {
			if d.ID == id {
				total := 1
				escribirJSON(w, http.StatusOK, RespuestaExito{Estado: 200, Datos: d, Total: &total})
				return
			}
		}
		escribirError(w, http.StatusNotFound, "No se encontró", "No existe un dragón con ese el id")
		return
	}

	filtrados := dragones

	if especie := q.Get("especie"); especie != "" {
		var tmp []models.Dragon
		for _, d := range filtrados {
			if strings.EqualFold(d.Especie, especie) {
				tmp = append(tmp, d)
			}
		}
		filtrados = tmp
	}

	if jinete := q.Get("jinete"); jinete != "" {
		var tmp []models.Dragon
		for _, d := range filtrados {
			if strings.EqualFold(d.Jinete, jinete) {
				tmp = append(tmp, d)
			}
		}
		filtrados = tmp
	}

	if alfaStr := q.Get("es_alfa"); alfaStr != "" {
		esAlfa, err := strconv.ParseBool(alfaStr)
		if err != nil {
			escribirError(w, http.StatusBadRequest, "Solicitud incorrecta", "El parámetro 'es_alfa' debe ser true o false")
			return
		}
		var tmp []models.Dragon
		for _, d := range filtrados {
			if d.EsAlfa == esAlfa {
				tmp = append(tmp, d)
			}
		}
		filtrados = tmp
	}

	if ubicacion := q.Get("ubicacion"); ubicacion != "" {
		var tmp []models.Dragon
		for _, d := range filtrados {
			if strings.EqualFold(d.Ubicacion, ubicacion) {
				tmp = append(tmp, d)
			}
		}
		filtrados = tmp
	}

	if filtrados == nil {
		filtrados = []models.Dragon{}
	}

	total := len(filtrados)
	escribirJSON(w, http.StatusOK, RespuestaExito{Estado: 200, Datos: filtrados, Total: &total})
}

func obtenerDragonPorRuta(w http.ResponseWriter, r *http.Request, id int) {
	dragones, err := storage.Cargar()
	if err != nil {
		escribirError(w, http.StatusInternalServerError, "Error", "No se pudieron cargar los dragones")
		return
	}
	for _, d := range dragones {
		if d.ID == id {
			escribirJSON(w, http.StatusOK, RespuestaExito{Estado: 200, Datos: d})
			return
		}
	}
	escribirError(w, http.StatusNotFound, "No se encontró", "No existe un dragón con el id "+strconv.Itoa(id))
}

func crearDragon(w http.ResponseWriter, r *http.Request) {
	var entrada models.Dragon
	if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
		escribirError(w, http.StatusBadRequest, "Solicitud incorrecta", "El body JSON no es válido")
		return
	}

	if msgError := validarDragon(entrada); msgError != "" {
		escribirError(w, http.StatusUnprocessableEntity, "Error de Validación", msgError)
		return
	}

	dragones, err := storage.Cargar()
	if err != nil {
		escribirError(w, http.StatusInternalServerError, "Error", "No se pudieron cargar los dragones")
		return
	}

	maxID := 0
	for _, d := range dragones {
		if d.ID > maxID {
			maxID = d.ID
		}
	}
	entrada.ID = maxID + 1

	dragones = append(dragones, entrada)

	if err := storage.Guardar(dragones); err != nil {
		escribirError(w, http.StatusInternalServerError, "Error", "No se pudo guardar el dragón")
		return
	}

	escribirJSON(w, http.StatusCreated, RespuestaExito{
		Estado:  201,
		Mensaje: "Dragón creado",
		Datos:   entrada,
	})
}

func actualizarDragon(w http.ResponseWriter, r *http.Request, id int) {
	var entrada models.Dragon
	if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
		escribirError(w, http.StatusBadRequest, "Solicitud incorrecta", "El body JSON no es válido")
		return
	}

	if msgError := validarDragon(entrada); msgError != "" {
		escribirError(w, http.StatusUnprocessableEntity, "Error de Validación", msgError)
		return
	}

	dragones, err := storage.Cargar()
	if err != nil {
		escribirError(w, http.StatusInternalServerError, "Error", "No se pudieron cargar los dragones")
		return
	}

	encontrado := false
	for i, d := range dragones {
		if d.ID == id {
			entrada.ID = id
			dragones[i] = entrada
			encontrado = true
			break
		}
	}

	if !encontrado {
		escribirError(w, http.StatusNotFound, "No se encontró", "No existe un dragón con el id "+strconv.Itoa(id))
		return
	}

	if err := storage.Guardar(dragones); err != nil {
		escribirError(w, http.StatusInternalServerError, "Error", "No se pudieron guardar los cambios")
		return
	}

	escribirJSON(w, http.StatusOK, RespuestaExito{
		Estado:  200,
		Mensaje: "Dragón actualizado",
		Datos:   entrada,
	})
}

func actualizarDragonParcial(w http.ResponseWriter, r *http.Request, id int) {
	var parcial map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&parcial); err != nil {
		escribirError(w, http.StatusBadRequest, "Solicitud incorrecta", "El body JSON no es válido")
		return
	}

	dragones, err := storage.Cargar()
	if err != nil {
		escribirError(w, http.StatusInternalServerError, "Error", "No se pudieron cargar los dragones")
		return
	}

	encontrado := false
	var actualizado models.Dragon
	for i, d := range dragones {
		if d.ID == id {
			if v, ok := parcial["nombre"].(string); ok && v != "" {
				d.Nombre = v
			}
			if v, ok := parcial["especie"].(string); ok && v != "" {
				d.Especie = v
			}
			if v, ok := parcial["jinete"].(string); ok && v != "" {
				d.Jinete = v
			}
			if v, ok := parcial["color"].(string); ok && v != "" {
				d.Color = v
			}
			if v, ok := parcial["habilidad"].(string); ok && v != "" {
				d.Habilidad = v
			}
			if v, ok := parcial["ubicacion"].(string); ok && v != "" {
				d.Ubicacion = v
			}
			if v, ok := parcial["envergadura_m"].(float64); ok {
				d.Envergadura = v
			}
			if v, ok := parcial["peso_kg"].(float64); ok {
				d.Peso = v
			}
			if v, ok := parcial["es_alfa"].(bool); ok {
				d.EsAlfa = v
			}
			dragones[i] = d
			actualizado = d
			encontrado = true
			break
		}
	}

	if !encontrado {
		escribirError(w, http.StatusNotFound, "No se encontró", "No existe un dragón con el id "+strconv.Itoa(id))
		return
	}

	if err := storage.Guardar(dragones); err != nil {
		escribirError(w, http.StatusInternalServerError, "Error", "No se pudieron guardar los cambios")
		return
	}

	escribirJSON(w, http.StatusOK, RespuestaExito{
		Estado:  200,
		Mensaje: "Dragón actualizado",
		Datos:   actualizado,
	})
}

func eliminarDragon(w http.ResponseWriter, r *http.Request, id int) {
	dragones, err := storage.Cargar()
	if err != nil {
		escribirError(w, http.StatusInternalServerError, "Error", "No se pudieron cargar los dragones")
		return
	}

	nuevaLista := make([]models.Dragon, 0, len(dragones))
	encontrado := false
	for _, d := range dragones {
		if d.ID == id {
			encontrado = true
			continue
		}
		nuevaLista = append(nuevaLista, d)
	}

	if !encontrado {
		escribirError(w, http.StatusNotFound, "No se encontró", "No existe un dragón con el id "+strconv.Itoa(id))
		return
	}

	if err := storage.Guardar(nuevaLista); err != nil {
		escribirError(w, http.StatusInternalServerError, "Error", "No se pudieron guardar los cambios")
		return
	}

	escribirJSON(w, http.StatusOK, RespuestaExito{
		Estado:  200,
		Mensaje: "Dragón con id " + strconv.Itoa(id) + " eliminado",
	})
}

func validarDragon(d models.Dragon) string {
	if strings.TrimSpace(d.Nombre) == "" {
		return "El campo 'nombre' es necesario"
	}
	if strings.TrimSpace(d.Especie) == "" {
		return "El campo 'especie' es necesario"
	}
	if strings.TrimSpace(d.Jinete) == "" {
		return "El campo 'jinete' es necesario"
	}
	if strings.TrimSpace(d.Color) == "" {
		return "El campo 'color' es necesario"
	}
	if strings.TrimSpace(d.Habilidad) == "" {
		return "El campo 'habilidad' es necesario"
	}
	if strings.TrimSpace(d.Ubicacion) == "" {
		return "El campo 'ubicacion' es necesario"
	}
	if d.Envergadura <= 0 {
		return "El campo 'envergadura_m' debe ser mayor a 0"
	}
	if d.Peso <= 0 {
		return "El campo 'peso_kg' debe ser mayor a 0"
	}
	return ""
}
