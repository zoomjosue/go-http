#  API REST — Cómo Entrenar a Tu Dragón

API RESTful construida en Go utilizando únicamente la librería estándar (`net/http`, `encoding/json`).
El tema es la película **Cómo Entrenar a Tu Dragón (2010)**, con un catálogo completo de dragones del universo de la película.


---

## Estructura del proyecto
```
GO-HTTP/
├── data/
│   └── dragons.json     # Base de datos en archivo JSON (persistente)
├── handlers/
│   ├── dragones.go      # Lógica de todos los endpoints
│   └── respuesta.go     # Helpers para respuestas JSON estandarizadas
├── models/
│   └── dragon.go        # Struct del modelo Dragon
├── storage/
│   └── storage.go       # Lectura y escritura del archivo JSON
├── .gitignore
├── Dockerfile
├── go.mod
├── main.go
└── README.md
```

---

##  Cómo ejecutar

### Sin Docker
```bash
go mod init GO-HTTP
go run .
```

### Con Docker
```bash
docker build -t httyd-api .
docker run -p 24918:24918 httyd-api
```

El servidor corre en el puerto **24918**.

---


---

##  Endpoints

### GET `/api/dragones`
Retorna todos los dragones registrados.
```
GET http://localhost:24918/api/dragones
```

**Response `200 OK`:**
```json
{
  "estado": 200,
  "datos": [ ... ],
  "total": 12
}
```

---

### GET `/api/dragones?id=1`
Filtra por ID usando query parameter.
```
GET http://localhost:24918/api/dragones?id=1
```

**Response `200 OK`:**
```json
{
  "estado": 200,
  "datos": {
    "id": 1,
    "nombre": "Chimuelo",
    "especie": "Furia Nocturna",
    "jinete": "Hipo",
    "color": "Negro",
    "envergadura_m": 14.5,
    "peso_kg": 180,
    "habilidad": "Explosión de Plasma",
    "es_alfa": true,
    "ubicacion": "Berk"
  },
  "total": 1
}
```

---

### GET `/api/dragones?especie=Gronckle&ubicacion=Berk`
Soporta múltiples filtros combinados.

**Query parameters disponibles:**

| Parámetro   | Tipo   | Ejemplo                  |
|-------------|--------|--------------------------|
| `id`        | int    | `?id=3`                  |
| `especie`   | string | `?especie=Gronckle`      |
| `jinete`    | string | `?jinete=Hipo`           |
| `es_alfa`   | bool   | `?es_alfa=true`          |
| `ubicacion` | string | `?ubicacion=Berk`        |

**Ejemplo combinado:**
```
GET http://localhost:24918/api/dragones?es_alfa=false&ubicacion=Berk
```

---

### GET `/api/dragones/{id}`
Obtiene un dragón por path parameter.
```
GET http://localhost:24918/api/dragones/3
```

**Response `200 OK`:**
```json
{
  "estado": 200,
  "datos": {
    "id": 3,
    "nombre": "Garfio",
    "especie": "Pesadilla Monstruosa",
    "jinete": "Patán",
    "color": "Rojo y Naranja",
    "envergadura_m": 16,
    "peso_kg": 200,
    "habilidad": "Auto Ignición",
    "es_alfa": false,
    "ubicacion": "Berk"
  }
}
```

---

### POST `/api/dragones`
Crea un nuevo dragón. El `id` se genera automáticamente.
```
POST http://localhost:24918/api/dragones
Content-Type: application/json
```

**Body:**
```json
{
  "nombre": "Sombra de Luna",
  "especie": "Furia Lunar",
  "jinete": "Luna",
  "color": "Blanco Perla",
  "envergadura_m": 13.5,
  "peso_kg": 170.0,
  "habilidad": "Rayo de Luna",
  "es_alfa": false,
  "ubicacion": "Santuario de Dragones"
}
```

**Response `201 Created`:**
```json
{
  "estado": 201,
  "mensaje": "Dragón creado exitosamente",
  "datos": { "id": 13, "nombre": "Sombra de Luna", "..." }
}
```

---

### PUT `/api/dragones/{id}`
Reemplaza completamente un dragón existente. Requiere todos los campos.
```
PUT http://localhost:24918/api/dragones/1
Content-Type: application/json
```

**Body:**
```json
{
  "nombre": "Chimuelo",
  "especie": "Furia Nocturna",
  "jinete": "Hipo Horrendo Abrazo III",
  "color": "Negro Brillante",
  "envergadura_m": 14.5,
  "peso_kg": 185.0,
  "habilidad": "Explosión de Plasma",
  "es_alfa": true,
  "ubicacion": "Berk"
}
```

**Response `200 OK`:**
```json
{
  "estado": 200,
  "mensaje": "Dragón actualizado exitosamente",
  "datos": { ... }
}
```

---

### PATCH `/api/dragones/{id}`
Actualiza parcialmente un dragón. Solo se envían los campos a modificar.
```
PATCH http://localhost:24918/api/dragones/2
Content-Type: application/json
```

**Body:**
```json
{
  "ubicacion": "Santuario de Dragones",
  "es_alfa": true
}
```

**Response `200 OK`:**
```json
{
  "estado": 200,
  "mensaje": "Dragón actualizado parcialmente",
  "datos": { ... }
}
```

---

### DELETE `/api/dragones/{id}`
Elimina un dragón por su ID.
```
DELETE http://localhost:24918/api/dragones/5
```

**Response `200 OK`:**
```json
{
  "estado": 200,
  "mensaje": "Dragón con id 5 eliminado exitosamente"
}
```

---

## ❌ Manejo de errores

Todos los errores devuelven JSON estructurado:
```json
{
  "estado": 404,
  "error": "No Encontrado",
  "mensaje": "No existe un dragón con id 999"
}
```

| Código | Situación                                        |
|--------|--------------------------------------------------|
| `400`  | JSON inválido o parámetro con tipo incorrecto    |
| `404`  | Dragón no encontrado                             |
| `405`  | Método HTTP no permitido en ese endpoint         |
| `422`  | Campo requerido faltante o valor inválido        |
| `500`  | Error interno al leer o escribir el archivo      |

---

## 🔒 Validaciones

Los siguientes campos son obligatorios en POST y PUT:

- `nombre`, `especie`, `jinete`, `color`, `habilidad`, `ubicacion` — no pueden estar vacíos
- `envergadura_m`, `peso_kg` — deben ser mayores a 0

---

## 💾 Persistencia

Todos los cambios (POST, PUT, PATCH, DELETE) se guardan directamente en `data/dragons.json`, garantizando que los datos persisten entre reinicios del servidor.

---

## 🛠️ Tecnologías

- **Lenguaje:** Go 1.22
- **Librerías:** Solo librería estándar de Go
- **Persistencia:** Archivo JSON
- **Contenedor:** Docker (multi-stage build)

---

## 👤 Información

**Carnet:** 24918  
**Tema:** Cómo Entrenar a Tu Dragón — API REST en Go