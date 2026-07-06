# Quill — Documentación de la API (Backend)

Este documento detalla todas las APIs (REST y WebSocket) expuestas por el backend de Quill. Explica para qué sirve cada endpoint, sus parámetros, los cuerpos de las solicitudes (`Request Body`) y las respuestas (`Response Body`).

---

## 📌 Aspectos Generales y Arquitectura

- **Framework Backend**: Go 1.22 con [Fiber v2](https://gofiber.io/).
- **Base de Datos**: PostgreSQL 16 con extensiones `pgvector` (búsqueda semántica) y `Apache AGE` (base de datos de grafos).
- **Ruta Base (Base URL)**: `/api/v1`
- **Autenticación**:
  - Los endpoints protegidos requieren un token JWT en la cabecera HTTP:
    ```http
    Authorization: Bearer <tu_token_jwt>
    ```
  - Si el token falta, ha expirado o es inválido, el servidor responderá con un error `401 Unauthorized`.

---

## 📋 Resumen de Endpoints

### 🩺 Sistema y Utilidades
* `GET /api/v1/health` — Verificación del estado de salud del sistema y dependencias.

### 🔑 Autenticación
* `POST /api/v1/auth/register` — Registro de un nuevo usuario.
* `POST /api/v1/auth/login` — Inicio de sesión y obtención del token JWT.
* `GET /api/v1/auth/me` — [Protegido] Obtención de datos del usuario autenticado actual.

### 🌌 Universos (Universes)
* `POST /api/v1/universes` — [Protegido] Crear un nuevo universo narrativo.
* `GET /api/v1/universes` — [Protegido] Listar universos con paginación.
* `GET /api/v1/universes/:id` — [Protegido] Obtener detalles de un universo por ID.
* `PUT /api/v1/universes/:id` — [Protegido] Actualizar información de un universo.
* `DELETE /api/v1/universes/:id` — [Protegido] Eliminar un universo (y todos sus datos relacionados).

### 📚 Obras (Works)
* `POST /api/v1/universes/:universe_id/works` — [Protegido] Crear una nueva obra (novela, guión) en un universo.
* `GET /api/v1/universes/:universe_id/works` — [Protegido] Listar todas las obras de un universo.
* `GET /api/v1/works/:id` — [Protegido] Obtener detalles de una obra por ID.
* `PUT /api/v1/works/:id` — [Protegido] Actualizar una obra.
* `DELETE /api/v1/works/:id` — [Protegido] Eliminar una obra.

### 📝 Capítulos (Chapters)
* `POST /api/v1/works/:work_id/chapters` — [Protegido] Crear un capítulo en una obra.
* `GET /api/v1/works/:work_id/chapters` — [Protegido] Listar capítulos de una obra.
* `GET /api/v1/chapters/:id` — [Protegido] Obtener detalles y contenido de un capítulo.
* `PUT /api/v1/chapters/:id` — [Protegido] Actualizar el título, contenido o texto sin formato de un capítulo.
* `DELETE /api/v1/chapters/:id` — [Protegido] Eliminar un capítulo.

### 👤 Entidades y Lore (Entities)
* `GET /api/v1/universes/:universe_id/entities` — [Protegido] Listar entidades (personajes, lugares, objetos) con filtros.
* `GET /api/v1/entities/:id` — [Protegido] Obtener detalles de una entidad.
* `PUT /api/v1/entities/:id` — [Protegido] Actualizar datos o propiedades de una entidad (lore, estatus, alias).

### ⚠️ Contradicciones (Contradictions)
* `GET /api/v1/universes/:universe_id/contradictions` — [Protegido] Obtener las contradicciones de lore detectadas en un universo.
* `PUT /api/v1/universes/:universe_id/contradictions/:id/resolve` — [Protegido] Marcar una contradicción como resuelta.
* `PUT /api/v1/universes/:universe_id/contradictions/:id/dismiss` — [Protegido] Descartar una contradicción sin resolverla en el texto.

### 📅 Línea de Tiempo (Timeline)
* `GET /api/v1/universes/:universe_id/timeline` — [Protegido] Obtener los eventos cronológicos del universo.
* `POST /api/v1/universes/:universe_id/timeline` — [Protegido] Agregar manualmente un evento a la línea de tiempo.

### 🕳️ Agujeros de Guión (Plot Holes)
* `GET /api/v1/universes/:universe_id/plot-holes` — [Protegido] Listar inconsistencias o tramas abiertas (agujeros de guión) detectadas por la IA.

### 🕸️ Grafo de Conocimiento y Memoria (Graph & Memory)
* `GET /api/v1/universes/:universe_id/graph` — [Protegido] Obtener todos los nodos y relaciones del grafo de conocimiento.
* `GET /api/v1/entities/:id/neighbors` — [Protegido] Obtener los vecinos y relaciones directas de una entidad específica (hasta N saltos).
* `POST /api/v1/universes/:id/recall` — [Protegido] Recuperar información contextual basada en búsqueda semántica vectorial.

### 📥 Ingesta de Documentos (Ingestion)
* `POST /api/v1/universes/:id/ingest` — [Protegido] Subir un archivo de texto/Markdown para su procesamiento asíncrono y extracción de entidades.

### 🔌 Canal de WebSocket (Real-time WebSocket)
* `GET /api/v1/ws` — [Protegido] Conexión bidireccional en tiempo real para análisis de texto interactivo.

### 🧪 Utilidades de Demostración (Demo)
* `POST /api/v1/demo/clone` — Clonar un universo de demostración plantilla asignándolo a una sesión.
* `POST /api/v1/demo/reset` — Restablecer los datos del universo de demostración asociado a una sesión.

---

## 🔎 Detalle de cada API

### 🩺 Sistema y Utilidades

#### `GET /api/v1/health`
* **Descripción**: Verifica el estado y conectividad del servidor, la base de datos PostgreSQL, las extensiones de base de datos (`age` y `vector`) y la API de Qwen Cloud.
* **Autenticación**: Pública (No requiere token).
* **Parámetros**: Ninguno.
* **Respuestas**:
  * **200 OK** (Si todos los servicios esenciales están funcionando correctamente):
    ```json
    {
      "status": "healthy",
      "db": "connected",
      "age": "available",
      "pgvector": "available",
      "qwen_api": "reachable",
      "disk_free_mb": 45120,
      "uptime_seconds": 3600
    }
    ```
  * **503 Service Unavailable** (Si la base de datos o extensiones críticas fallan):
    * El campo `"status"` cambiará a `"unhealthy"` o `"degraded"` (si el fallo es solo de la API de Qwen).

---

### 🔑 Autenticación

#### `POST /api/v1/auth/register`
* **Descripción**: Registra una nueva cuenta de escritor en la plataforma.
* **Autenticación**: Pública.
* **Cuerpo de la Solicitud (JSON)**:
  ```json
  {
    "email": "escritor@ejemplo.com",
    "password": "contrasena_segura_de_minimo_8_caracteres",
    "display_name": "Nombre de Autor"
  }
  ```
* **Respuestas**:
  * **201 Created**:
    ```json
    {
      "user": {
        "id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
        "email": "escritor@ejemplo.com",
        "display_name": "Nombre de Autor",
        "created_at": "2026-07-06T00:00:00Z",
        "updated_at": "2026-07-06T00:00:00Z"
      },
      "token": "eyJhbGciOiJIUzI1NiIsIn..."
    }
    ```
  * **400 Bad Request**: Si faltan campos obligatorios o la contraseña es menor a 8 caracteres.
  * **409 Conflict**: Si el correo ya se encuentra registrado.

#### `POST /api/v1/auth/login`
* **Descripción**: Valida las credenciales del escritor y devuelve un token JWT utilizable para consumir el resto de APIs protegidas.
* **Autenticación**: Pública.
* **Cuerpo de la Solicitud (JSON)**:
  ```json
  {
    "email": "escritor@ejemplo.com",
    "password": "contrasena_segura"
  }
  ```
* **Respuestas**:
  * **200 OK**:
    ```json
    {
      "user": {
        "id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
        "email": "escritor@ejemplo.com",
        "display_name": "Nombre de Autor",
        "created_at": "2026-07-06T00:00:00Z",
        "updated_at": "2026-07-06T00:00:00Z"
      },
      "token": "eyJhbGciOiJIUzI1NiIsIn..."
    }
    ```
  * **401 Unauthorized**: Credenciales de inicio de sesión inválidas.

#### `GET /api/v1/auth/me`
* **Descripción**: Retorna la información de perfil del usuario que está realizando la solicitud, identificándolo a través del token JWT provisto en la cabecera `Authorization`.
* **Autenticación**: Protegida.
* **Respuestas**:
  * **200 OK**:
    ```json
    {
      "user": {
        "id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
        "email": "escritor@ejemplo.com",
        "display_name": "Nombre de Autor",
        "created_at": "2026-07-06T00:00:00Z",
        "updated_at": "2026-07-06T00:00:00Z"
      }
    }
    ```

---

### 🌌 Universos (Universes)
Un *Universo* representa el macrocosmos de lore o el mundo donde se desarrollan las historias del autor (ej. "Tierra Media", "Cosmere").

#### `POST /api/v1/universes`
* **Descripción**: Crea un nuevo universo de escritura. Al crearse, el backend inicializa automáticamente un grafo en Apache AGE exclusivo para el universo, nombrado `universe_<universe_uuid>`.
* **Autenticación**: Protegida.
* **Cuerpo de la Solicitud (JSON)**:
  ```json
  {
    "name": "Cronicas de Aethelgard",
    "description": "Un mundo de fantasia epica con magia basada en runas y constelaciones.",
    "genre": "Fantasy",
    "format": "novel"
  }
  ```
* **Respuestas**:
  * **201 Created**: Retorna el objeto del universo creado con su ID único generado por el sistema.

#### `GET /api/v1/universes`
* **Descripción**: Lista los universos pertenecientes al usuario autenticado. Soporta paginación.
* **Autenticación**: Protegida.
* **Parámetros Query**:
  - `page`: Número de página (por defecto `1`).
  - `limit`: Cantidad de registros por página (por defecto `20`).
* **Respuestas**:
  * **200 OK**:
    ```json
    {
      "universes": [ ... ],
      "pagination": {
        "page": 1,
        "limit": 20,
        "total": 3,
        "total_pages": 1
      }
    }
    ```

#### `GET /api/v1/universes/:id`
* **Descripción**: Obtiene la información detallada de un universo mediante su UUID.
* **Autenticación**: Protegida.
* **Respuestas**:
  * **200 OK**: Objeto de universo.
  * **404 Not Found**: Si el universo no existe o no pertenece al usuario actual.

#### `PUT /api/v1/universes/:id`
* **Descripción**: Modifica los metadatos de un universo existente (nombre, descripción, género, formato).
* **Autenticación**: Protegida.
* **Cuerpo de la Solicitud (JSON)**: Igual que `CreateUniverseRequest`.
* **Respuestas**:
  * **200 OK**: Objeto del universo actualizado.

#### `DELETE /api/v1/universes/:id`
* **Descripción**: Borra de manera definitiva un universo. Esto eliminará en cascada todas las obras, capítulos, entidades, contradicciones, líneas de tiempo y el propio grafo en Apache AGE.
* **Autenticación**: Protegida.
* **Respuestas**:
  * **204 No Content**: Operación realizada con éxito.

---

### 📚 Obras (Works)
Una *Obra* (Work) representa una novela, libro o guión específico que transcurre dentro de un determinado universo.

#### `POST /api/v1/universes/:universe_id/works`
* **Descripción**: Crea una nueva obra asociada a un universo.
* **Autenticación**: Protegida.
* **Cuerpo de la Solicitud (JSON)**:
  ```json
  {
    "title": "El Despertar de la Runa",
    "type": "novel",
    "synopsis": "Un joven aprendiz descubre una runa prohibida capaz de rasgar el velo estelar."
  }
  ```
* **Respuestas**:
  * **201 Created**: Objeto de obra creado.

#### `GET /api/v1/universes/:universe_id/works`
* **Descripción**: Obtiene una lista de todas las obras creadas dentro de un universo.
* **Autenticación**: Protegida.
* **Respuestas**:
  * **200 OK**: `{"works": [...]}`

#### `GET /api/v1/works/:id`
* **Descripción**: Obtiene los detalles de una obra por su ID.
* **Autenticación**: Protegida.

#### `PUT /api/v1/works/:id`
* **Descripción**: Modifica la información de una obra (título, sinopsis, tipo).
* **Autenticación**: Protegida.
* **Cuerpo de la Solicitud (JSON)**: Estructura idéntica a la creación de obras.

#### `DELETE /api/v1/works/:id`
* **Descripción**: Elimina una obra y sus capítulos asociados.
* **Autenticación**: Protegida.
* **Respuestas**:
  * **204 No Content**

---

### 📝 Capítulos (Chapters)

#### `POST /api/v1/works/:work_id/chapters`
* **Descripción**: Crea un nuevo capítulo vacío dentro de una obra. El backend calcula automáticamente el índice de orden (`order_index`) al final de los capítulos existentes.
* **Autenticación**: Protegida.
* **Cuerpo de la Solicitud (JSON)**:
  ```json
  {
    "title": "Capitulo 1: El Guardian del Faro"
  }
  ```
* **Respuestas**:
  * **201 Created**: Objeto de capítulo.

#### `GET /api/v1/works/:work_id/chapters`
* **Descripción**: Obtiene la lista de todos los capítulos que componen una obra en particular.
* **Autenticación**: Protegida.

#### `GET /api/v1/chapters/:id`
* **Descripción**: Obtiene el contenido completo y los metadatos de un capítulo específico (incluyendo texto raw, recuento de palabras y fecha de análisis).
* **Autenticación**: Protegida.

#### `PUT /api/v1/chapters/:id`
* **Descripción**: Actualiza el contenido escrito del capítulo. **Nota:** Esta API REST se usa para actualizaciones síncronas o globales del capítulo. El guardado en tiempo real al escribir se gestiona mediante el canal de WebSocket.
* **Autenticación**: Protegida.
* **Cuerpo de la Solicitud (JSON)**:
  ```json
  {
    "title": "Capitulo 1: El Guardian del Faro (Editado)",
    "content": "<p>El viento aullaba sobre los acantilados de Aethelgard...</p>",
    "raw_text": "El viento aullaba sobre los acantilados de Aethelgard..."
  }
  ```
* **Respuestas**:
  * **200 OK**: Objeto del capítulo actualizado.

#### `DELETE /api/v1/chapters/:id`
* **Descripción**: Elimina permanentemente el capítulo de la obra.
* **Autenticación**: Protegida.

---

### 👤 Entidades y Lore (Entities)
Las *Entidades* son los componentes del lore del universo extraídos del texto por la IA o creados por el usuario (Personajes, Lugares, Objetos, Eventos o Conceptos).

#### `GET /api/v1/universes/:universe_id/entities`
* **Descripción**: Retorna las entidades registradas en el universo. Permite filtrar y buscar por diferentes campos para poblar vistas como la enciclopedia de lore.
* **Autenticación**: Protegida.
* **Parámetros Query**:
  - `type`: Filtrar por tipo (ej. `character`, `location`, `item`).
  - `status`: Filtrar por estatus (ej. `active`, `archived`).
  - `min_relevance`: Filtrar por una puntuación mínima de relevancia de lore (para descartar personajes muy secundarios).
  - `search`: Cadena de búsqueda textual para coincidencia con nombres y alias de entidades.
  - `page`: Número de página (defecto `1`).
  - `limit`: Límite por página (defecto `50`).
* **Respuestas**:
  * **200 OK**:
    ```json
    {
      "entities": [
        {
          "id": "e458ff62-602c-473d-82d2-8b64a2f8b50f",
          "universe_id": "8bb2f15b-...",
          "type": "character",
          "name": "Kaelen",
          "aliases": ["Kael", "El Guardian de las Runas"],
          "description": "Un joven aprendiz de escribano en el faro de Aethelgard.",
          "properties": {
            "age": "19",
            "hair_color": "silver",
            "affiliation": "Faro de Aethelgard"
          },
          "status": "active",
          "relevance_score": 0.95,
          "created_at": "2026-07-06T00:00:00Z",
          "updated_at": "2026-07-06T00:00:00Z"
        }
      ],
      "pagination": { ... }
    }
    ```

#### `GET /api/v1/entities/:id`
* **Descripción**: Devuelve los detalles completos de una entidad de lore concreta.
* **Autenticación**: Protegida.

#### `PUT /api/v1/entities/:id`
* **Descripción**: Permite a un autor editar o corregir manualmente la ficha técnica o de lore de una entidad.
* **Autenticación**: Protegida.
* **Cuerpo de la Solicitud (JSON)**:
  ```json
  {
    "name": "Kaelen Rion",
    "aliases": ["Kael", "El Guardian de las Runas", "Kaelen de Aethelgard"],
    "description": "El aprendiz principal del faro de Aethelgard, ahora portador de la runa estelar.",
    "status": "active",
    "properties": {
      "age": "20",
      "hair_color": "silver",
      "affiliation": "Orden Estelar"
    }
  }
  ```
* **Respuestas**:
  * **200 OK**: Ficha de la entidad modificada.

---

### ⚠️ Contradicciones (Contradictions)
El backend cuenta con un validador asíncrono basado en IA y memoria que contrasta los párrafos escritos contra el lore establecido en busca de contradicciones (por ejemplo, si un personaje muerto reaparece o cambia el color de sus ojos sin explicación).

#### `GET /api/v1/universes/:universe_id/contradictions`
* **Descripción**: Retorna la lista de todas las contradicciones de lore detectadas en el universo que aún no han sido descartadas.
* **Autenticación**: Protegida.
* **Respuestas**:
  * **200 OK**:
    ```json
    {
      "contradictions": [
        {
          "id": "c112aa52-0941-4775-bebe-5abfcf5522e8",
          "universe_id": "8bb2f15b-...",
          "entity_id": "e458ff62-602c-473d-82d2-8b64a2f8b50f",
          "severity": "high",
          "description": "Se menciona que Kaelen esta en el faro de Aethelgard, pero en el Capitulo 2 se establecio que la Orden Estelar lo capturo en los calabozos de Veridia.",
          "suggestion": "Aclarar como Kaelen escapo de Veridia o ajustar su ubicacion en el faro.",
          "evidence_a": "Kaelen se asomo a la ventana del faro sintiendo el viento salado.",
          "evidence_a_chapter_id": "3bb62aa2-...",
          "evidence_b": "Los guardias encadenaron a Kaelen en la fosa mas profunda de la prision de Veridia.",
          "evidence_b_chapter_id": "2ee62aa2-...",
          "status": "detected",
          "resolved_at": null,
          "created_at": "2026-07-06T01:30:00Z"
        }
      ]
    }
    ```

#### `PUT /api/v1/universes/:universe_id/contradictions/:id/resolve`
* **Descripción**: Marca una contradicción de lore como resuelta. El autor utiliza esto para indicar que ha editado el texto o que la inconsistencia ya no aplica.
* **Autenticación**: Protegida.
* **Respuestas**:
  * **200 OK**: `{"status": "resolved"}`

#### `PUT /api/v1/universes/:universe_id/contradictions/:id/dismiss`
* **Descripción**: Descarta o ignora una contradicción detectada por el sistema. Sirve para cuando la inconsistencia es intencional (por ejemplo, el personaje está mintiendo o hay un misterio planificado).
* **Autenticación**: Protegida.
* **Respuestas**:
  * **200 OK**: `{"status": "dismissed"}`

---

### 📅 Línea de Tiempo (Timeline)
La línea de tiempo consolida acontecimientos cronológicos ordenados para mantener la coherencia temporal de la narrativa.

#### `GET /api/v1/universes/:universe_id/timeline`
* **Descripción**: Retorna la lista de todos los eventos del timeline de un universo, ordenados cronológicamente por su posición numérica en la línea de tiempo.
* **Autenticación**: Protegida.
* **Respuestas**:
  * **200 OK**:
    ```json
    {
      "events": [
        {
          "id": "fd90a162-...",
          "universe_id": "8bb2f15b-...",
          "event_entity_id": "e458ff62-602c-473d-82d2-8b64a2f8b50f",
          "title": "Nacimiento de Kaelen",
          "description": "Nace bajo la conjuncion de la Constelacion del Fenix.",
          "timeline_position": 820.5,
          "timeline_label": "Año 820 de la Era Estelar",
          "chapter_id": null,
          "participants": ["e458ff62-..."],
          "created_at": "2026-07-06T01:00:00Z"
        }
      ]
    }
    ```

#### `POST /api/v1/universes/:universe_id/timeline`
* **Descripción**: Agrega manualmente un nuevo evento a la línea de tiempo del universo.
* **Autenticación**: Protegida.
* **Cuerpo de la Solicitud (JSON)**:
  ```json
  {
    "title": "La Caida de Aethelgard",
    "description": "Las tropas de Veridia asedian y destruyen el faro de runas.",
    "timeline_position": 840.1,
    "timeline_label": "Otoño del Año 840",
    "event_entity_id": "3bb62aa2-...",
    "participants": ["e458ff62-...", "9aa5231c-..."]
  }
  ```
* **Respuestas**:
  * **201 Created**: Objeto de evento de línea de tiempo con su ID asignado.

---

### 🕳️ Agujeros de Guión (Plot Holes)

#### `GET /api/v1/universes/:universe_id/plot-holes`
* **Descripción**: Retorna los agujeros de guión detectados por la IA en un universo. Un agujero de guión representa cabos sueltos graves, tales como un personaje con un rol clave que desaparece de repente por varios capítulos sin explicarse adónde fue.
* **Autenticación**: Protegida.
* **Respuestas**:
  * **200 OK**:
    ```json
    {
      "plot_holes": [
        {
          "id": "7aa90f11-...",
          "universe_id": "8bb2f15b-...",
          "title": "Desaparicion inexplicable de Elara",
          "description": "Elara tenia un rol prominente hasta el Capitulo 3, pero no ha vuelto a mencionarse ni aparecer en los ultimos 4 capitulos sin justificacion.",
          "related_entity_ids": ["9aa5231c-..."],
          "first_mentioned_chapter_id": "1ab62aa2-...",
          "status": "open",
          "created_at": "2026-07-06T01:45:00Z"
        }
      ]
    }
    ```

---

### 🕸️ Grafo de Conocimiento y Memoria (Graph & Memory)
Permite visualizar y consultar las relaciones semánticas entre entidades dentro de la base de datos de grafos Apache AGE.

#### `GET /api/v1/universes/:universe_id/graph`
* **Descripción**: Consulta y retorna el grafo completo (todos los nodos y aristas/relaciones) del universo para renderizar mapas visuales interactivos (como la red de relaciones del lore).
* **Autenticación**: Protegida.
* **Respuestas**:
  * **200 OK**:
    ```json
    {
      "nodes": [
        { "id": "e458ff62-...", "label": "character", "properties": { "name": "Kaelen" } }
      ],
      "edges": [
        { "id": "rel123", "label": "APPRENTICE_OF", "source": "e458ff62-...", "target": "3bb62aa2-...", "properties": {} }
      ]
    }
    ```

#### `GET /api/v1/entities/:id/neighbors`
* **Descripción**: Obtiene los vecinos de una entidad dada a través de un recorrido de N saltos en el grafo de relaciones.
* **Autenticación**: Protegida.
* **Parámetros Query**:
  - `universe_id` (Obligatorio): UUID del universo al que pertenece la entidad.
  - `hops`: Saltos máximos de profundidad a recorrer (por defecto `1`, máximo permitido `5` para proteger el rendimiento).
* **Respuestas**:
  * **200 OK**: JSON con estructura `{"nodes": [...], "edges": [...]}` filtrado en torno a la entidad inicial.

#### `POST /api/v1/universes/:id/recall`
* **Descripción**: Endpoint de búsqueda semántica (memoria asociativa) usado para consultar hechos o lore relevante del universo a través de vectores y embeddings.
* **Autenticación**: Protegida.
* **Cuerpo de la Solicitud (JSON)**:
  ```json
  {
    "query": "Como funciona la magia de runas?",
    "k": 5
  }
  ```
* **Respuestas**:
  * **200 OK**: Retorna una lista de hechos ordenados por relevancia semántica (calculada con distancia de coseno sobre los embeddings del texto).
    ```json
    {
      "items": [
        {
          "entity_id": "bf89a162-...",
          "fact": "La magia requiere dibujar runas fisicas y cargarlas mediante la luz de las constelaciones durante la medianoche.",
          "score": 0.895,
          "source": "Ficha de Lore: Magia de Runas"
        }
      ]
    }
    ```

---

### 📥 Ingesta de Documentos (Ingestion)

#### `POST /api/v1/universes/:id/ingest`
* **Descripción**: Permite subir un archivo de texto o Markdown (`.md`, `.txt`) que contiene una novela o borrador completo previamente escrito. El servidor recibe el archivo e inicia una tarea asíncrona en segundo plano que:
  1. Divide el texto en capítulos (basándose en cabeceras de Markdown).
  2. Segmenta en párrafos.
  3. Extrae entidades y sus relaciones semánticas con IA.
  4. Genera embeddings vectoriales de los párrafos.
  5. Popula la base de datos relacional, vectorial y el grafo de Apache AGE.
* **Autenticación**: Protegida.
* **Cuerpo de la Solicitud**: Formulario Multipart (`multipart/form-data`) que debe incluir un campo `file` con el archivo seleccionado.
* **Respuestas**:
  * **202 Accepted**: Retorna un ID de trabajo (`job_id`) para monitorear el progreso del análisis asíncrono.
    ```json
    {
      "job_id": "8aa12ff2-0941-4775-bebe-5abfcf5522e8",
      "status": "accepted"
    }
    ```
  * El progreso de la ingesta se transmite en tiempo real al usuario mediante eventos de WebSocket (`ingestion_progress`).

---

### 🔌 Canal de WebSocket (Real-time WebSocket)

#### `GET /api/v1/ws`
* **Descripción**: Abre una conexión persistente por WebSocket para habilitar el análisis de texto dinámico en tiempo real mientras el autor escribe en el editor de texto.
* **Autenticación**: Protegida (Habilitada pasando el token JWT en el protocolo de inicialización o mediante la URL de conexión).
* **Protocolo de Comunicación (WSMessage)**:
  Toda comunicación se empaqueta en una estructura común con un tipo y un payload:
  ```json
  {
    "type": "NOMBRE_DEL_EVENTO",
    "payload": { ... }
  }
  ```

#### Mensajes Cliente → Servidor (Enviados por la App Web)
1. **`auth_init`**: Autentica la sesión de WebSocket tras conectar.
   ```json
   {
     "type": "auth_init",
     "payload": { "token": "<JWT_TOKEN>" }
   }
   ```
2. **`paragraph_submit`**: Envía un fragmento de texto o párrafo recién editado para su análisis de lore en tiempo real.
   ```json
   {
     "type": "paragraph_submit",
     "payload": {
       "universe_id": "8bb2f15b-...",
       "work_id": "5aa90f11-...",
       "chapter_id": "3bb62aa2-...",
       "text": "Kaelen dibujo la runa estelar en el aire con pulso tembloroso, ignorando que las leyes del faro lo prohibian."
     }
   }
   ```

#### Mensajes Servidor → Cliente (Recibidos por la App Web)
1. **`analysis_result`**: Enviado por la IA tras completar el análisis de un párrafo enviado.
   * Contiene una lista breve de entidades reconocidas en el párrafo, las contradicciones detectadas y los agujeros de guión generados o modificados.
2. **`contradiction_alert`**: Notificación push inmediata si la IA detecta que el último párrafo contradice una regla crítica del lore establecido.
3. **`entity_discovered`**: Se envía cuando la IA detecta la mención de una nueva entidad de lore y la ha creado, o cuando una existente ha variado significativamente de estado.
4. **`graph_updated`**: Alerta que el grafo del universo ha sufrido cambios estructurales para que el cliente refresque la vista visual del mapa.
5. **`ingestion_progress`**: Informa del progreso del trabajo de ingesta asíncrona de un libro completo (ej. `"capítulos procesados: 5/12"`).

---

### 🧪 Utilidades de Demostración (Demo)
*Estas APIs están diseñadas para permitir a los usuarios del sitio probar el editor y el sistema de grafos de forma interactiva con datos semilla preconfigurados.*

#### `POST /api/v1/demo/clone`
* **Descripción**: Clona el universo plantilla de demostración por defecto en la base de datos y lo asocia con la sesión temporal del navegador identificada en la cabecera `X-Session-ID`.
* **Autenticación**: Pública (Identificada por cabecera de sesión).
* **Cabecera HTTP**: `X-Session-ID: <session_uuid>` (si no se envía, el backend genera uno nuevo).
* **Respuestas**:
  * **200 OK**:
    ```json
    {
      "status": "success",
      "universe_id": "9bb2f15b-...",
      "message": "Demo universe cloned successfully"
    }
    ```

#### `POST /api/v1/demo/reset`
* **Descripción**: Elimina los datos de prueba modificados por el usuario para su sesión y vuelve a clonar el universo semilla original de demostración.
* **Autenticación**: Pública.
* **Cabecera HTTP**: `X-Session-ID: <session_uuid>` (Obligatoria).
