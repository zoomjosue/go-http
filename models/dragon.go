package models

type Dragon struct {
	ID          int     `json:"id"`
	Nombre      string  `json:"nombre"`
	Especie     string  `json:"especie"`
	Jinete      string  `json:"jinete"`
	Color       string  `json:"color"`
	Envergadura float64 `json:"envergadura_m"`
	Peso        float64 `json:"peso_kg"`
	Habilidad   string  `json:"habilidad"`
	EsAlfa      bool    `json:"es_alfa"`
	Ubicacion   string  `json:"ubicacion"`
}
