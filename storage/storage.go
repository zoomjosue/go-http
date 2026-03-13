package storage

import (
	"encoding/json"
	"os"
	"sync"

	"go-http/models"
)

const rutaArchivo = "data/dragons.json"

var mutex sync.Mutex

func Cargar() ([]models.Dragon, error) {
	mutex.Lock()
	defer mutex.Unlock()

	datos, err := os.ReadFile(rutaArchivo)
	if err != nil {
		return nil, err
	}

	var dragones []models.Dragon
	if err := json.Unmarshal(datos, &dragones); err != nil {
		return nil, err
	}
	return dragones, nil
}

func Guardar(dragones []models.Dragon) error {
	mutex.Lock()
	defer mutex.Unlock()

	datos, err := json.MarshalIndent(dragones, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(rutaArchivo, datos, 0644)
}
