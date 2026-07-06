# Plan Maestro: Quill → Track 1 MemoryAgent Winner

## Resumen Ejecutivo

Este documento detalla **6 soluciones técnicas** que atacan directamente los criterios de evaluación donde el proyecto es débil. Cada solución incluye: el problema exacto, la solución arquitectónica completa, los archivos a modificar/crear, los algoritmos y estructuras de datos, y finalmente un desglose en tareas ejecutables.

---

## SOLUCIÓN 1: Context-Window Budget Manager
### Criterio atacado: Technical Depth (30%) — *"Recalling critical memories within limited context windows"*

> [!CAUTION]
> Esta es la debilidad más fatal. Las reglas **literalmente** piden esto como capacidad core. Sin esto, el proyecto no puede ganar.

### 1.1 El Problema

Actualmente los prompts en [qwen_service.go](file:///home/daikyri/Workspace/Hackathon-QwenCloud/backend/internal/services/qwen_service.go) concatenan texto + entidades + tool results sin límite alguno. No hay conteo de tokens, no hay selección por presupuesto, no hay compresión. Si un universo tiene 500 entidades y 2000 párrafos, el prompt puede exceder la ventana de contexto y fallar silenciosamente o truncarse.

### 1.2 La Solución Completa

Implementar un **ContextBudgetManager** — un servicio que:

1. **Cuenta tokens** de cualquier string usando un tokenizer compatible con Qwen (tiktoken cl100k_base, que es el tokenizer de la familia Qwen)
2. **Asigna presupuestos** a cada sección del prompt (system, context, entities, tool results, user text)
3. **Selecciona memorias** por ranking hasta llenar el presupuesto
4. **Comprime** memorias que no caben completas usando summarization via Qwen-turbo
5. **Reporta métricas** de uso del context window al frontend (para visualización)

#### Arquitectura del Budget Manager

```
┌──────────────────────────────────────────┐
│           ContextBudgetManager           │
│                                          │
│  MaxTokens: 30000 (qwen-max = 32k)      │
│  ReservedForResponse: 2000               │
│  Available: 28000                        │
│                                          │
│  ┌──────────────────────────────────┐    │
│  │ Budget Allocation:               │    │
│  │  SystemPrompt:    2000 (fixed)   │    │
│  │  UserText:        3000 (dynamic) │    │
│  │  EntityContext:   8000 (ranked)  │    │
│  │  VectorMemories: 10000 (ranked) │    │
│  │  ToolResults:     5000 (capped)  │    │
│  └──────────────────────────────────┘    │
│                                          │
│  Methods:                                │
│   - CountTokens(text) → int             │
│   - AllocateBudget(sections) → Budget   │
│   - FitMemories(items, budget) → fitted │
│   - CompressIfNeeded(text, max) → text  │
│   - Report() → BudgetReport             │
└──────────────────────────────────────────┘
```

#### 1.2.1 Token Counting

Crear un módulo Go que use un tokenizer. Dado que no existe un port nativo perfecto de tiktoken en Go, usaremos la librería `github.com/pkoukk/tiktoken-go` que implementa el encoding cl100k_base (compatible con la tokenización de Qwen).

**Archivo nuevo:** `backend/internal/services/tokenizer.go`

```go
package services

import (
    "sync"
    "github.com/pkoukk/tiktoken-go"
)

// Tokenizer provides token counting for context window budget management.
// Uses cl100k_base encoding (compatible with Qwen model family).
type Tokenizer struct {
    enc  *tiktoken.Tiktoken
    once sync.Once
}

var globalTokenizer = &Tokenizer{}

// CountTokens returns the number of tokens in the given text.
// Thread-safe via sync.Once initialization.
func (t *Tokenizer) CountTokens(text string) int {
    t.once.Do(func() {
        enc, err := tiktoken.GetEncoding("cl100k_base")
        if err != nil {
            // Fallback: approximate 1 token ≈ 4 chars
            return
        }
        t.enc = enc
    })
    if t.enc == nil {
        // Fallback approximation
        return len(text) / 4
    }
    return len(t.enc.Encode(text, nil, nil))
}

// CountTokensForMessages counts tokens for a full message array,
// including the per-message overhead (role tokens, separators).
func (t *Tokenizer) CountTokensForMessages(messages []QwenMessage) int {
    total := 0
    for _, msg := range messages {
        total += 4 // every message has: <|im_start|>{role}\n ... <|im_end|>\n
        total += t.CountTokens(msg.Role)
        total += t.CountTokens(msg.Content)
        for _, tc := range msg.ToolCalls {
            total += t.CountTokens(tc.Function.Name)
            total += t.CountTokens(tc.Function.Arguments)
        }
    }
    total += 2 // priming
    return total
}
```

#### 1.2.2 Budget Manager Service

**Archivo nuevo:** `backend/internal/services/context_budget.go`

```go
package services

// BudgetAllocation defines the token limits for each prompt section.
type BudgetAllocation struct {
    SystemPrompt   int
    UserText       int
    EntityContext  int
    VectorMemories int
    ToolResults    int
    Total          int
}

// BudgetReport is sent to the frontend to visualize context window usage.
type BudgetReport struct {
    MaxTokens       int     `json:"max_tokens"`
    UsedTokens      int     `json:"used_tokens"`
    UtilizationPct  float64 `json:"utilization_pct"`
    Sections        map[string]SectionReport `json:"sections"`
    CompressedCount int     `json:"compressed_count"`
    DroppedCount    int     `json:"dropped_count"`
}

type SectionReport struct {
    Budget    int `json:"budget"`
    Used      int `json:"used"`
    ItemCount int `json:"item_count"`
}

// ContextBudgetManager manages token allocation across prompt sections.
type ContextBudgetManager struct {
    tokenizer *Tokenizer
    maxTokens int // Model's context window (e.g., 30000 for qwen-max)
    responseReserve int // Tokens reserved for response generation
}

// NewContextBudgetManager creates a budget manager.
// maxContextTokens is the model's total window; responseReserve is
// how many tokens to leave for the model's response.
func NewContextBudgetManager(maxContextTokens, responseReserve int) *ContextBudgetManager {
    return &ContextBudgetManager{
        tokenizer:       globalTokenizer,
        maxTokens:       maxContextTokens,
        responseReserve: responseReserve,
    }
}

// ComputeBudget distributes available tokens across sections using
// proportional allocation with fixed minimums for system/user.
func (m *ContextBudgetManager) ComputeBudget(systemPromptTokens, userTextTokens int) BudgetAllocation {
    available := m.maxTokens - m.responseReserve - systemPromptTokens - userTextTokens
    if available < 0 {
        available = 0
    }
    // Proportional: entities 35%, vector 40%, tools 25%
    return BudgetAllocation{
        SystemPrompt:   systemPromptTokens,
        UserText:       userTextTokens,
        EntityContext:  int(float64(available) * 0.35),
        VectorMemories: int(float64(available) * 0.40),
        ToolResults:    int(float64(available) * 0.25),
        Total:          m.maxTokens - m.responseReserve,
    }
}

// RankedItem is any item that can be fitted into a budget section.
type RankedItem struct {
    ID       string
    Text     string
    Score    float64
    Tokens   int // pre-computed
}

// FitToBudget selects top-scored items that fit within the token budget.
// Items must be pre-sorted by score descending.
// Returns: fitted items, items that were dropped, total tokens used.
func (m *ContextBudgetManager) FitToBudget(items []RankedItem, budgetTokens int) (fitted []RankedItem, dropped int, tokensUsed int) {
    for _, item := range items {
        if item.Tokens == 0 {
            item.Tokens = m.tokenizer.CountTokens(item.Text)
        }
        if tokensUsed + item.Tokens <= budgetTokens {
            fitted = append(fitted, item)
            tokensUsed += item.Tokens
        } else {
            dropped++
        }
    }
    return fitted, dropped, tokensUsed
}

// CountTokens delegates to the internal tokenizer.
func (m *ContextBudgetManager) CountTokens(text string) int {
    return m.tokenizer.CountTokens(text)
}
```

#### 1.2.3 Integrar el Budget Manager en el Analysis Pipeline

**Archivo a modificar:** [analysis_service.go](file:///home/daikyri/Workspace/Hackathon-QwenCloud/backend/internal/services/analysis_service.go)

Añadir `budgetMgr *ContextBudgetManager` al struct `AnalysisService`. Inyectar en el constructor desde `main.go`.

**Archivo a modificar:** [contradiction_service.go](file:///home/daikyri/Workspace/Hackathon-QwenCloud/backend/internal/services/contradiction_service.go)

En `CheckSemantic`, antes de construir el `userMessage`:
1. Contar tokens del system prompt
2. Contar tokens del user text
3. Calcular budget restante para entity context
4. Seleccionar solo las top-k entidades que caben en el budget
5. Si alguna entity description es muy larga, comprimirla vía Qwen-turbo

**Archivo a modificar:** [qwen_service.go](file:///home/daikyri/Workspace/Hackathon-QwenCloud/backend/internal/services/qwen_service.go)

En `RunAgentLoop`, después de cada iteración del loop:
1. Contar tokens acumulados en `msgs`
2. Si se acerca al límite, resumir los tool results anteriores en un mensaje comprimido
3. Emitir `BudgetReport` al final para el frontend

**Archivo a modificar:** [main.go](file:///home/daikyri/Workspace/Hackathon-QwenCloud/backend/cmd/server/main.go)

Crear el `ContextBudgetManager` con los valores del modelo:
```go
budgetMgr := services.NewContextBudgetManager(30000, 2000) // qwen-max: ~32k window
```

Inyectarlo en `AnalysisService`, `ContradictionService`, `PlotHoleService`, y `QwenService`.

#### 1.2.4 Config

**Archivo a modificar:** [config.go](file:///home/daikyri/Workspace/Hackathon-QwenCloud/backend/internal/config/config.go)

Añadir:
```go
MaxContextTokens    int  // QWEN_MAX_CONTEXT_TOKENS, default 30000
ResponseReserve     int  // QWEN_RESPONSE_RESERVE, default 2000
```

---

## SOLUCIÓN 2: Índices HNSW en pgvector
### Criterio atacado: Technical Depth (30%) — *"Efficient storage/retrieval at scale"*

### 2.1 El Problema

Las tablas `entity_embeddings` y `paragraph_embeddings` tienen columnas `vector(1024)` pero **ningún índice vectorial**. Las queries `<=>` (cosine distance) hacen full table scan O(n). Con 10k+ embeddings esto es inaceptablemente lento.

### 2.2 La Solución

Crear una nueva migración que añada índices HNSW a ambas tablas. HNSW (Hierarchical Navigable Small Worlds) da búsquedas aproximadas en O(log n) con recall >95%.

**Archivo nuevo:** `backend/migrations/015_add_hnsw_indexes.up.sql`

```sql
-- HNSW indexes for approximate nearest neighbor search.
-- vector_cosine_ops matches the <=> (cosine distance) operator used in queries.
-- m=16, ef_construction=64 are balanced defaults for 1024-dim vectors.

CREATE INDEX IF NOT EXISTS idx_entity_embeddings_hnsw 
ON entity_embeddings 
USING hnsw (description_embedding vector_cosine_ops)
WITH (m = 16, ef_construction = 64);

CREATE INDEX IF NOT EXISTS idx_paragraph_embeddings_hnsw 
ON paragraph_embeddings 
USING hnsw (embedding vector_cosine_ops)
WITH (m = 16, ef_construction = 64);
```

**Archivo nuevo:** `backend/migrations/015_add_hnsw_indexes.down.sql`

```sql
DROP INDEX IF EXISTS idx_entity_embeddings_hnsw;
DROP INDEX IF EXISTS idx_paragraph_embeddings_hnsw;
```

**Archivo a modificar:** [vector_repo.go](file:///home/daikyri/Workspace/Hackathon-QwenCloud/backend/internal/repositories/vector_repo.go)

Cambiar `FindSimilarEntity` para que retorne **top-k** en lugar de solo 1 resultado. Esto también es necesario para la Solución 4 (Hybrid Retrieval):

```go
// FindSimilarEntities returns top-k entities by cosine similarity.
// Uses HNSW index for O(log n) approximate nearest neighbor search.
func (r *VectorRepo) FindSimilarEntities(ctx context.Context, universeID uuid.UUID, 
    embedding []float32, threshold float64, limit int) ([]SimilarEntity, error) {
    
    query := `
        SELECT e.id, e.name, ee.description_embedding <=> $1 AS distance
        FROM entities e
        JOIN entity_embeddings ee ON e.id = ee.entity_id
        WHERE e.universe_id = $2 AND ee.description_embedding <=> $1 < $3
        ORDER BY distance ASC
        LIMIT $4
    `
    rows, err := r.pool.Query(ctx, query, embedding, universeID, threshold, limit)
    // ... scan rows into []SimilarEntity ...
}

type SimilarEntity struct {
    ID       uuid.UUID
    Name     string
    Distance float64
}
```

También añadir un parámetro `ef_search` al inicio de las queries para controlar la precisión vs velocidad del HNSW en runtime:

```go
// SetHNSWSearchParams sets the ef_search parameter for the current session.
// Higher values = more accurate but slower. Default: 40, range: 1-1000.
func (r *VectorRepo) SetHNSWSearchParams(ctx context.Context, efSearch int) error {
    _, err := r.pool.Exec(ctx, fmt.Sprintf("SET hnsw.ef_search = %d", efSearch))
    return err
}
```

---

## SOLUCIÓN 3: Memory Consolidation & Summarization
### Criterio atacado: Technical Depth (30%) — *"Timely forgetting of outdated information"*

### 3.1 El Problema

El decay actual (`score *= e^(-λ)`) marca entidades como "archived" cuando bajan de un threshold. Pero esas entidades archivadas y sus párrafos asociados siguen ocupando la misma cantidad de tokens en la DB y potencialmente en los prompts. No hay **compactación ni consolidación**. Las reglas piden: *"mechanisms to decay, prune, or overwrite stale memories"*.

### 3.2 La Solución

Implementar un **Memory Consolidation Pipeline** que, al archivar entidades, genere un resumen comprimido de su historia y lo almacene como una "compressed memory". Esto convierte 20 párrafos de una entidad en 1-2 oraciones que preservan los hechos clave.

#### 3.2.1 Nuevo modelo y tabla

**Archivo a modificar:** [models.go](file:///home/daikyri/Workspace/Hackathon-QwenCloud/backend/internal/models/models.go)

Añadir:
```go
// ConsolidatedMemory is a compressed summary of an archived entity's history.
// Created when an entity's relevance score drops below the archive threshold.
// Replaces the need to scan all individual paragraphs for that entity.
type ConsolidatedMemory struct {
    ID           uuid.UUID  `json:"id"`
    UniverseID   uuid.UUID  `json:"universe_id"`
    EntityID     uuid.UUID  `json:"entity_id"`
    Summary      string     `json:"summary"`       // Compressed narrative summary
    KeyFacts     []string   `json:"key_facts"`      // Extracted bullet points
    TokenCount   int        `json:"token_count"`    // Pre-computed token count
    MentionCount int        `json:"mention_count"`  // Original number of mentions
    Embedding    []float32  `json:"embedding"`      // Embedding of the summary
    CreatedAt    time.Time  `json:"created_at"`
}
```

**Archivo nuevo:** `backend/migrations/016_create_consolidated_memories.up.sql`

```sql
CREATE TABLE consolidated_memories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    universe_id UUID NOT NULL REFERENCES universes(id) ON DELETE CASCADE,
    entity_id UUID NOT NULL REFERENCES entities(id) ON DELETE CASCADE,
    summary TEXT NOT NULL,
    key_facts TEXT[] DEFAULT '{}',
    token_count INTEGER NOT NULL DEFAULT 0,
    mention_count INTEGER NOT NULL DEFAULT 0,
    embedding vector(1024),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(entity_id)
);

CREATE INDEX idx_consolidated_memories_universe ON consolidated_memories(universe_id);
CREATE INDEX idx_consolidated_memories_embedding_hnsw 
ON consolidated_memories USING hnsw (embedding vector_cosine_ops)
WITH (m = 16, ef_construction = 64);
```

**Archivo nuevo:** `backend/migrations/016_create_consolidated_memories.down.sql`
```sql
DROP TABLE IF EXISTS consolidated_memories;
```

#### 3.2.2 Consolidation Repository

**Archivo nuevo:** `backend/internal/repositories/consolidation_repo.go`

```go
package repositories

type ConsolidationRepo struct {
    pool *pgxpool.Pool
}

func NewConsolidationRepo(pool *pgxpool.Pool) *ConsolidationRepo { ... }

func (r *ConsolidationRepo) Create(ctx context.Context, cm *models.ConsolidatedMemory) error { ... }
func (r *ConsolidationRepo) FindByEntityID(ctx context.Context, entityID uuid.UUID) (*models.ConsolidatedMemory, error) { ... }
func (r *ConsolidationRepo) FindByUniverse(ctx context.Context, universeID uuid.UUID) ([]models.ConsolidatedMemory, error) { ... }
func (r *ConsolidationRepo) Delete(ctx context.Context, entityID uuid.UUID) error { ... } // For reactivation
```

#### 3.2.3 Consolidation Service

**Archivo nuevo:** `backend/internal/services/consolidation_service.go`

```go
package services

// ConsolidationService creates compressed memory summaries when entities
// are archived, and deletes them when entities are reactivated.
type ConsolidationService struct {
    consolidationRepo *repositories.ConsolidationRepo
    entityRepo        *repositories.EntityRepo
    vectorRepo        *repositories.VectorRepo
    qwenSvc           *QwenService
    budgetMgr         *ContextBudgetManager
}

// ConsolidateEntity is called when an entity transitions to "archived".
// It:
// 1. Fetches all mentions of the entity
// 2. Sends them to Qwen-turbo for summarization
// 3. Extracts key facts as bullet points
// 4. Generates an embedding of the summary
// 5. Stores the consolidated memory
// 6. Optionally prunes old paragraph embeddings for this entity
func (s *ConsolidationService) ConsolidateEntity(ctx context.Context, entityID, universeID uuid.UUID) error {
    // 1. Fetch all mentions
    mentions, _ := s.entityRepo.GetMentionsByEntity(ctx, entityID, 100)
    entity, _ := s.entityRepo.FindByID(ctx, entityID)
    
    // 2. Build context from mentions
    var mentionTexts []string
    for _, m := range mentions {
        mentionTexts = append(mentionTexts, m.ContextSnippet)
    }
    
    // 3. Call Qwen-turbo to summarize
    summaryPrompt := fmt.Sprintf(`Summarize the complete narrative history of entity "%s" (%s) 
    based on these mentions. Include: current status, key relationships, important events, 
    and any unresolved plot threads. Be concise but preserve ALL factual details.
    
    Mentions:
    %s
    
    Respond with JSON: {"summary": "...", "key_facts": ["fact1", "fact2", ...]}`,
        entity.Name, entity.Type, strings.Join(mentionTexts, "\n---\n"))
    
    // 4. Parse response, generate embedding, store
    // ...
}

// DeconsolidateEntity removes the consolidated memory when an entity
// is reactivated (mentioned again after being archived).
func (s *ConsolidationService) DeconsolidateEntity(ctx context.Context, entityID uuid.UUID) error {
    return s.consolidationRepo.Delete(ctx, entityID)
}
```

#### 3.2.4 Integrar en el flujo de Decay

**Archivo a modificar:** [relevance_service.go](file:///home/daikyri/Workspace/Hackathon-QwenCloud/backend/internal/services/relevance_service.go)

Modificar `DecayAll` para que, después de archivar entidades, llame al `ConsolidationService` para cada entidad recién archivada:

```go
func (s *RelevanceService) DecayAll(ctx context.Context, universeID uuid.UUID) error {
    // 1. Decay scores (existing)
    if err := s.entityRepo.DecayAll(ctx, universeID, s.lambda); err != nil {
        return err
    }
    
    // 2. Find entities that just crossed the threshold
    newlyArchived, _ := s.entityRepo.FindNewlyArchivable(ctx, universeID, s.archiveThreshold)
    
    // 3. Archive them
    _, err := s.pool.Exec(ctx, `
        UPDATE entities SET status = 'archived', updated_at = NOW()
        WHERE universe_id = $1 AND status = 'active' AND relevance_score <= $2
    `, universeID, s.archiveThreshold)
    
    // 4. Consolidate each newly archived entity (async, best-effort)
    if s.consolidationSvc != nil {
        for _, entityID := range newlyArchived {
            go func(eid uuid.UUID) {
                if err := s.consolidationSvc.ConsolidateEntity(context.Background(), eid, universeID); err != nil {
                    log.Printf("[relevance] consolidate entity %s: %v", eid, err)
                }
            }(entityID)
        }
    }
    
    return err
}
```

**Archivo a modificar:** [relevance_service.go](file:///home/daikyri/Workspace/Hackathon-QwenCloud/backend/internal/services/relevance_service.go)

Modificar `Reactivate` para que elimine la consolidated memory:

```go
func (s *RelevanceService) Reactivate(ctx context.Context, entityID uuid.UUID) error {
    // ... existing reactivation logic ...
    
    // Remove consolidated memory since entity is active again
    if s.consolidationSvc != nil {
        s.consolidationSvc.DeconsolidateEntity(ctx, entityID)
    }
    
    return tx.Commit(ctx)
}
```

**Nuevo método en EntityRepo:** `FindNewlyArchivable`

```go
func (r *EntityRepo) FindNewlyArchivable(ctx context.Context, universeID uuid.UUID, threshold float64) ([]uuid.UUID, error) {
    query := `SELECT id FROM entities WHERE universe_id = $1 AND status = 'active' AND relevance_score <= $2`
    // ... return IDs ...
}
```

---

## SOLUCIÓN 4: Hybrid Retrieval con Reciprocal Rank Fusion
### Criterio atacado: Technical Depth (30%) — *"hybrid retrieval, context-window budgeting techniques"*

### 4.1 El Problema

El [MemoryService.Recall](file:///home/daikyri/Workspace/Hackathon-QwenCloud/backend/internal/services/memory_service.go#L50-L160) actual tiene múltiples problemas:
1. Hace N graph queries (una por entidad activa) — O(n²)
2. Solo retorna 1 resultado del vector search
3. Los pesos de combinación (0.4/0.3/0.3) están hardcoded
4. No hay keyword search
5. No es "hybrid retrieval" real — es un merge ad-hoc

### 4.2 La Solución

Reescribir `MemoryService.Recall` usando **Reciprocal Rank Fusion (RRF)** — un algoritmo estándar de information retrieval que combina rankings de múltiples fuentes sin necesidad de normalizar scores.

#### Algoritmo RRF

```
RRF_score(doc) = Σ (1 / (k + rank_i(doc)))
```

Donde `k` es una constante (típicamente 60) y `rank_i(doc)` es la posición del documento en el ranking `i`.

#### 4.2.1 Reescribir MemoryService

**Archivo a reescribir:** [memory_service.go](file:///home/daikyri/Workspace/Hackathon-QwenCloud/backend/internal/services/memory_service.go)

```go
package services

const rrfK = 60 // Standard RRF constant

// HybridRecallSource identifies which retrieval pipeline produced a result.
type HybridRecallSource string

const (
    SourceVector   HybridRecallSource = "vector"
    SourceGraph    HybridRecallSource = "graph"
    SourceKeyword  HybridRecallSource = "keyword"
    SourceRecency  HybridRecallSource = "recency"
    SourceConsolidated HybridRecallSource = "consolidated"
)

// HybridRecallItem extends RecallItem with source tracking and RRF scoring.
type HybridRecallItem struct {
    EntityID    uuid.UUID
    EntityName  string
    Fact        string
    RRFScore    float64
    Sources     []HybridRecallSource
    TokenCount  int // Pre-computed for budget fitting
}

// MemoryService provides hybrid retrieval by combining multiple retrieval
// pipelines using Reciprocal Rank Fusion (RRF).
type MemoryService struct {
    graphRepo         *repositories.GraphRepo
    entityRepo        *repositories.EntityRepo
    vectorRepo        *repositories.VectorRepo
    consolidationRepo *repositories.ConsolidationRepo
    budgetMgr         *ContextBudgetManager
}

// Recall executes 4 parallel retrieval pipelines and fuses their results
// using RRF to produce a single ranked list of memory items.
//
// Pipelines:
//  1. Vector similarity: embed query → top-k paragraphs by cosine distance
//  2. Graph context: find entities connected to query entities → 1-hop neighbors (batch)
//  3. Keyword/recency: active entities sorted by relevance_score (already computed by decay)
//  4. Consolidated: archived entity summaries matching by embedding similarity
//
// The result is then fitted to the context window budget via FitToBudget.
func (s *MemoryService) Recall(ctx context.Context, universeID uuid.UUID, 
    queryEmbedding []float32, queryText string, k int) ([]models.RecallItem, *BudgetReport, error) {
    
    // ── Pipeline 1: Vector Similarity ──
    vectorResults := s.vectorPipeline(ctx, universeID, queryEmbedding, k*2)
    
    // ── Pipeline 2: Graph Neighbors (batch, not per-entity) ──
    graphResults := s.graphPipeline(ctx, universeID, vectorResults)
    
    // ── Pipeline 3: Recency/Keyword ──
    recencyResults := s.recencyPipeline(ctx, universeID, queryText, k*2)
    
    // ── Pipeline 4: Consolidated Memories ──
    consolidatedResults := s.consolidatedPipeline(ctx, universeID, queryEmbedding, k)
    
    // ── Reciprocal Rank Fusion ──
    fused := s.fuseRRF(vectorResults, graphResults, recencyResults, consolidatedResults)
    
    // ── Budget Fitting ──
    if s.budgetMgr != nil {
        budget := s.budgetMgr.ComputeBudget(0, 0) // just for memory section
        items := make([]RankedItem, len(fused))
        for i, f := range fused {
            items[i] = RankedItem{ID: f.EntityID.String(), Text: f.Fact, Score: f.RRFScore}
        }
        fitted, dropped, used := s.budgetMgr.FitToBudget(items, budget.VectorMemories)
        // Convert back to RecallItems...
        report := &BudgetReport{
            MaxTokens: budget.VectorMemories,
            UsedTokens: used,
            DroppedCount: dropped,
        }
        return recallItems, report, nil
    }
    
    // No budget manager: just return top-k
    if len(fused) > k {
        fused = fused[:k]
    }
    return toRecallItems(fused), nil, nil
}

// vectorPipeline: semantic search via paragraph embeddings
func (s *MemoryService) vectorPipeline(ctx context.Context, universeID uuid.UUID, 
    embedding []float32, limit int) []rankedEntry {
    
    paragraphs, err := s.vectorRepo.FindSimilarParagraphs(ctx, universeID, embedding, uuid.Nil, limit)
    if err != nil { return nil }
    
    entries := make([]rankedEntry, len(paragraphs))
    for i, p := range paragraphs {
        entries[i] = rankedEntry{
            id:   p.ChapterID.String() + ":" + p.Content[:min(50, len(p.Content))],
            fact: p.Content,
            score: 1.0 - p.Distance,
        }
    }
    // Already sorted by distance ASC (= score DESC)
    return entries
}

// graphPipeline: BATCH graph neighbor lookup for top entities from vector results
func (s *MemoryService) graphPipeline(ctx context.Context, universeID uuid.UUID, 
    vectorResults []rankedEntry) []rankedEntry {
    
    graphName := "universe_" + universeID.String()
    // Get top 5 entity IDs from vector results (avoid N queries for all entities)
    seen := map[string]bool{}
    var entries []rankedEntry
    
    // Only traverse neighbors of top-5 entities (not ALL active entities)
    entityIDs := s.extractTopEntityIDs(ctx, universeID, vectorResults, 5)
    
    for _, eid := range entityIDs {
        neighbors, err := s.graphRepo.GetNeighbors(ctx, graphName, eid)
        if err != nil { continue }
        for _, n := range neighbors {
            key := n.Node + ":" + n.RelType
            if seen[key] { continue }
            seen[key] = true
            entries = append(entries, rankedEntry{
                id:   key,
                fact: fmt.Sprintf("%s → %s → %s", eid, n.RelType, n.Node),
                score: 1.0,
            })
        }
    }
    return entries
}

// recencyPipeline: active entities sorted by relevance_score
func (s *MemoryService) recencyPipeline(ctx context.Context, universeID uuid.UUID,
    queryText string, limit int) []rankedEntry {
    
    entities, _ := s.entityRepo.ListByUniverseActive(ctx, universeID)
    entries := make([]rankedEntry, 0, min(len(entities), limit))
    for i, e := range entities {
        if i >= limit { break }
        entries = append(entries, rankedEntry{
            id:    e.ID.String(),
            fact:  fmt.Sprintf("%s (%s): %s", e.Name, e.Type, truncate(e.Description, 100)),
            score: e.RelevanceScore,
        })
    }
    return entries
}

// consolidatedPipeline: search compressed memories of archived entities
func (s *MemoryService) consolidatedPipeline(ctx context.Context, universeID uuid.UUID,
    embedding []float32, limit int) []rankedEntry {
    
    if s.consolidationRepo == nil || len(embedding) == 0 { return nil }
    
    memories, _ := s.consolidationRepo.FindSimilarByEmbedding(ctx, universeID, embedding, limit)
    entries := make([]rankedEntry, len(memories))
    for i, m := range memories {
        entries[i] = rankedEntry{
            id:   m.EntityID.String() + ":consolidated",
            fact: m.Summary,
            score: 1.0 - m.Distance,
        }
    }
    return entries
}

// fuseRRF combines multiple ranked lists using Reciprocal Rank Fusion.
func (s *MemoryService) fuseRRF(lists ...[]rankedEntry) []HybridRecallItem {
    scores := map[string]*HybridRecallItem{}
    
    for listIdx, list := range lists {
        source := []HybridRecallSource{SourceVector, SourceGraph, SourceRecency, SourceConsolidated}[listIdx]
        for rank, entry := range list {
            rrfScore := 1.0 / float64(rrfK + rank + 1)
            
            if existing, ok := scores[entry.id]; ok {
                existing.RRFScore += rrfScore
                existing.Sources = append(existing.Sources, source)
            } else {
                scores[entry.id] = &HybridRecallItem{
                    Fact:     entry.fact,
                    RRFScore: rrfScore,
                    Sources:  []HybridRecallSource{source},
                }
            }
        }
    }
    
    // Sort by RRF score descending
    result := make([]HybridRecallItem, 0, len(scores))
    for _, item := range scores {
        result = append(result, *item)
    }
    sort.Slice(result, func(i, j int) bool {
        return result[i].RRFScore > result[j].RRFScore
    })
    
    return result
}
```

#### 4.2.2 Keyword Search (PostgreSQL Full-Text)

**Archivo nuevo:** `backend/migrations/017_add_fulltext_index.up.sql`

```sql
-- Full-text search index for keyword-based hybrid retrieval
ALTER TABLE paragraph_embeddings ADD COLUMN IF NOT EXISTS tsv tsvector 
    GENERATED ALWAYS AS (to_tsvector('english', coalesce(content, ''))) STORED;
    
CREATE INDEX IF NOT EXISTS idx_paragraph_embeddings_tsv 
ON paragraph_embeddings USING gin(tsv);
```

**Archivo a modificar:** [vector_repo.go](file:///home/daikyri/Workspace/Hackathon-QwenCloud/backend/internal/repositories/vector_repo.go)

Añadir:
```go
// KeywordSearch performs PostgreSQL full-text search on paragraph content.
func (r *VectorRepo) KeywordSearch(ctx context.Context, universeID uuid.UUID, 
    query string, limit int) ([]SimilarParagraph, error) {
    
    sql := `
        SELECT pe.content, pe.chapter_id, c.title, 
               ts_rank(pe.tsv, websearch_to_tsquery('english', $1)) AS rank
        FROM paragraph_embeddings pe
        JOIN chapters c ON pe.chapter_id = c.id
        JOIN works w ON c.work_id = w.id
        WHERE w.universe_id = $2 AND pe.tsv @@ websearch_to_tsquery('english', $1)
        ORDER BY rank DESC
        LIMIT $3
    `
    // ... scan rows ...
}
```

---

## SOLUCIÓN 5: Qwen API — Structured Output + Streaming
### Criterio atacado: Innovation & AI Creativity (30%) — *"Sophisticated use of Qwen Cloud APIs"*

### 5.1 El Problema

1. No se usa `response_format` para JSON mode — hay un hack de "strip markdown fences" como fallback
2. No hay streaming — el usuario espera sin feedback visual hasta que toda la respuesta llega
3. Los modelos están hardcoded como strings — no se usan los config values `QwenMaxModel` y `QwenTurboModel` que ya existen en config

### 5.2 La Solución

#### 5.2.1 Structured Output (JSON Mode)

**Archivo a modificar:** [qwen_service.go](file:///home/daikyri/Workspace/Hackathon-QwenCloud/backend/internal/services/qwen_service.go)

Cambiar el `QwenRequest` para soportar `response_format`:

```go
type QwenRequest struct {
    Model          string            `json:"model"`
    Messages       []QwenMessage     `json:"messages"`
    ResponseFormat *ResponseFormat   `json:"response_format,omitempty"`
    Tools          []QwenTool        `json:"tools,omitempty"`
    ToolChoice     interface{}       `json:"tool_choice,omitempty"`
    Stream         bool              `json:"stream,omitempty"`
}

type ResponseFormat struct {
    Type string `json:"type"` // "json_object" or "text"
}
```

Modificar **todas** las funciones que esperan JSON de vuelta para añadir `ResponseFormat: &ResponseFormat{Type: "json_object"}`:

- `ExtractEntities` — línea 242-248
- `AnalyzeRelationships` — línea 338-344
- `CheckContradictions` — línea 395-401

Esto elimina la necesidad del hack "strip markdown fences" en `CheckSemantic`.

#### 5.2.2 Usar Config Models en lugar de hardcoded strings

**Archivo a modificar:** [qwen_service.go](file:///home/daikyri/Workspace/Hackathon-QwenCloud/backend/internal/services/qwen_service.go)

Añadir campos al struct:
```go
type QwenService struct {
    client     *http.Client
    baseURL    string
    apiKey     string
    maxModel   string // from cfg.QwenMaxModel
    turboModel string // from cfg.QwenTurboModel
    embModel   string // from cfg.QwenEmbeddingModel
    maxSem     chan struct{}
    turboSem   chan struct{}
}
```

Reemplazar todos los `"qwen-turbo"` y `"qwen-max"` hardcoded con `s.turboModel` y `s.maxModel`.

#### 5.2.3 Streaming para Analysis Pipeline

Implementar streaming SSE para que el usuario vea el progreso de la analysis en tiempo real en lugar de esperar el resultado completo.

**Archivo nuevo:** `backend/internal/services/streaming.go`

```go
package services

// StreamChunk represents a partial response from a streaming Qwen call.
type StreamChunk struct {
    Delta   string `json:"delta"`   // Incremental text
    Done    bool   `json:"done"`    // True if this is the final chunk
    ToolCall *QwenToolCall `json:"tool_call,omitempty"`
}

// StreamingResponse wraps a channel that receives chunks.
type StreamingResponse struct {
    Chunks <-chan StreamChunk
    Err    error
}
```

**Archivo a modificar:** [qwen_service.go](file:///home/daikyri/Workspace/Hackathon-QwenCloud/backend/internal/services/qwen_service.go)

Añadir método de streaming:

```go
// ChatCompletionStream sends a streaming request to Qwen and returns 
// a channel of chunks. The caller must drain the channel.
func (s *QwenService) ChatCompletionStream(ctx context.Context, payload QwenRequest) (<-chan StreamChunk, error) {
    payload.Stream = true
    
    // ... setup HTTP request ...
    
    ch := make(chan StreamChunk, 100)
    go func() {
        defer close(ch)
        scanner := bufio.NewScanner(resp.Body)
        for scanner.Scan() {
            line := scanner.Text()
            if !strings.HasPrefix(line, "data: ") { continue }
            data := strings.TrimPrefix(line, "data: ")
            if data == "[DONE]" { 
                ch <- StreamChunk{Done: true}
                return 
            }
            // Parse SSE delta...
            var chunk streamResponse
            json.Unmarshal([]byte(data), &chunk)
            if len(chunk.Choices) > 0 {
                ch <- StreamChunk{Delta: chunk.Choices[0].Delta.Content}
            }
        }
    }()
    
    return ch, nil
}
```

**Archivo a modificar:** [ws/protocol.go](file:///home/daikyri/Workspace/Hackathon-QwenCloud/backend/internal/ws/protocol.go)

Añadir tipo de mensaje para streaming progress:
```go
const TypeAnalysisProgress = "analysis_progress"
```

**Archivo a modificar:** [analysis_service.go](file:///home/daikyri/Workspace/Hackathon-QwenCloud/backend/internal/services/analysis_service.go)

En `processJob`, enviar mensajes de progreso intermedios:
```go
// After entity extraction
s.sendProgress(job.UserID, "entities_extracted", len(resolvedEntities))

// Before contradiction check
s.sendProgress(job.UserID, "checking_contradictions", 0)

// After semantic check
s.sendProgress(job.UserID, "contradictions_checked", len(result.Contradictions))

// After plot hole scan
s.sendProgress(job.UserID, "plot_holes_scanned", len(result.PlotHoles))

// With budget report
s.sendProgress(job.UserID, "context_budget", budgetReport)
```

---

## SOLUCIÓN 6: Visualización del Ciclo de Memoria en el Frontend
### Criterio atacado: Presentation & Documentation (15%) — *"Is the memory logic clearly visualized?"*

### 6.1 El Problema

El frontend no muestra:
- El decay de relevance scores en el tiempo
- El context window budget y cómo se seleccionan las memorias
- El ciclo de vida de las entidades (active → decaying → archived → consolidated → reactivated)
- Las fuentes de retrieval (qué vino de vector, qué de graph, qué de keyword)

### 6.2 La Solución

Crear un **Memory Inspector Panel** — un componente lateral que aparece junto al editor y muestra en tiempo real el estado del sistema de memoria.

#### 6.2.1 Nuevo componente: MemoryInspectorPanel

**Archivo nuevo:** `frontend/src/components/memory-inspector/MemoryInspectorPanel.tsx`

Este panel tiene 4 tabs:

**Tab 1: Context Budget** — Visualiza cómo se distribuyó el presupuesto de tokens:
- Barra horizontal segmentada mostrando system/user/entities/vector/tools
- Porcentaje de utilización
- Número de memorias incluidas vs excluidas

**Tab 2: Entity Lifecycle** — Timeline vertical de cada entidad:
- Sparkline de relevance_score a lo largo de los capítulos
- Indicadores de estado (🟢 active, 🟡 decaying, 🔴 archived, 📦 consolidated, 🔄 reactivated)
- Click para expandir y ver la consolidated memory summary

**Tab 3: Retrieval Sources** — Muestra de dónde vino cada memoria recallada:
- Badges de color por fuente (vector=azul, graph=verde, keyword=naranja, consolidated=púrpura)
- RRF score de cada item
- Highlight de items que aparecen en múltiples fuentes

**Tab 4: Live Pipeline** — Vista de las etapas del analysis pipeline:
- Steps: Entity Extraction → Contradiction Check → Plot Hole Scan → Memory Recall
- Estado de cada paso (pending/running/done)
- Tiempos de ejecución

#### 6.2.2 Nuevo store: memoryStore

**Archivo nuevo:** `frontend/src/stores/memoryStore.ts`

```typescript
import { create } from 'zustand'

interface EntityLifecycle {
  entityId: string
  entityName: string
  type: string
  status: 'active' | 'decaying' | 'archived' | 'consolidated' | 'reactivated'
  relevanceScore: number
  scoreHistory: { chapter: number; score: number }[]
  consolidatedSummary?: string
}

interface BudgetSnapshot {
  maxTokens: number
  usedTokens: number
  utilizationPct: number
  sections: Record<string, { budget: number; used: number; itemCount: number }>
  compressedCount: number
  droppedCount: number
}

interface RecallSource {
  fact: string
  rrfScore: number
  sources: string[]
  tokenCount: number
}

interface PipelineStep {
  name: string
  status: 'pending' | 'running' | 'done' | 'error'
  durationMs?: number
  resultCount?: number
}

interface MemoryState {
  entityLifecycles: EntityLifecycle[]
  budgetSnapshot: BudgetSnapshot | null
  recallSources: RecallSource[]
  pipelineSteps: PipelineStep[]
  
  updateBudget: (budget: BudgetSnapshot) => void
  updateLifecycles: (entities: EntityLifecycle[]) => void
  updateRecallSources: (sources: RecallSource[]) => void
  updatePipelineStep: (step: PipelineStep) => void
}
```

#### 6.2.3 Conectar el wsStore al memoryStore

**Archivo a modificar:** [wsStore.ts](file:///home/daikyri/Workspace/Hackathon-QwenCloud/frontend/src/stores/wsStore.ts)

Añadir handlers para los nuevos message types:

```typescript
case 'analysis_progress':
    useMemoryStore.getState().updatePipelineStep(payload as PipelineStep)
    break
case 'budget_report':
    useMemoryStore.getState().updateBudget(payload as BudgetSnapshot)
    break
case 'recall_sources':
    useMemoryStore.getState().updateRecallSources(payload.items as RecallSource[])
    break
```

#### 6.2.4 API endpoint para entity lifecycle data

**Archivo a modificar:** [handlers/graph.go](file:///home/daikyri/Workspace/Hackathon-QwenCloud/backend/internal/handlers/graph.go)

Añadir endpoint `GET /api/v1/universes/:id/memory-status` que retorne:
- Todas las entidades con su status y relevance_score
- Memorias consolidadas
- Estadísticas de memoria (total activas, archivadas, consolidadas)

**Archivo a modificar:** [main.go](file:///home/daikyri/Workspace/Hackathon-QwenCloud/backend/cmd/server/main.go)

Registrar la ruta:
```go
api.Get("/universes/:id/memory-status", graphH.MemoryStatus)
```

#### 6.2.5 Componentes UI específicos

**Archivo nuevo:** `frontend/src/components/memory-inspector/BudgetBar.tsx`
- Barra horizontal segmentada con colores por sección
- Animación al actualizar

**Archivo nuevo:** `frontend/src/components/memory-inspector/EntityLifecycleTimeline.tsx`
- Sparkline SVG del score decay
- Iconos de status con tooltip

**Archivo nuevo:** `frontend/src/components/memory-inspector/RetrievalSourceBadges.tsx`
- Badges con colores: `--vector-blue`, `--graph-green`, `--keyword-orange`, `--consolidated-purple`

**Archivo nuevo:** `frontend/src/components/memory-inspector/PipelineProgress.tsx`
- Stepper vertical con estados animados

**Archivo nuevo:** `frontend/src/components/memory-inspector/MemoryInspector.module.css`
- Estilos siguiendo el design system existente (ink/paper palette, serif typography)

---

## SPRINT: Desglose en Tareas Ejecutables

> [!IMPORTANT]
> Las tareas están ordenadas por dependencia. Cada tarea es autocontenida y testeable individualmente. La numeración refleja el orden de ejecución óptimo.

### FASE 1: Infraestructura Base (Tareas 1-5)

Estas tareas crean la base que las demás soluciones necesitan.

#### Tarea 1: Migración HNSW
- **Archivo:** Crear `backend/migrations/015_add_hnsw_indexes.up.sql` y `.down.sql`
- **Contenido:** Los `CREATE INDEX USING hnsw` exactos de la Solución 2
- **Verificación:** `docker compose up postgres`, ejecutar migraciones, confirmar que los índices existen con `\di` en psql
- **Tiempo estimado:** 15 min

#### Tarea 2: Dependencia tiktoken-go
- **Comando:** `cd backend && go get github.com/pkoukk/tiktoken-go`
- **Verificación:** `go build ./...` compila sin errores
- **Tiempo estimado:** 5 min

#### Tarea 3: Tokenizer
- **Archivo:** Crear `backend/internal/services/tokenizer.go`
- **Contenido:** Exactamente el código de la sección 1.2.1
- **Test:** Crear `backend/internal/services/tokenizer_test.go` — test que "Hello world" ≈ 2 tokens, string vacío = 0, string largo cuenta más
- **Verificación:** `go test ./internal/services/ -run TestTokenizer`
- **Tiempo estimado:** 30 min

#### Tarea 4: Context Budget Manager
- **Archivo:** Crear `backend/internal/services/context_budget.go`
- **Contenido:** Exactamente el código de la sección 1.2.2
- **Test:** Crear `backend/internal/services/context_budget_test.go` — test ComputeBudget con varios tamaños, FitToBudget con items que exceden budget
- **Verificación:** `go test ./internal/services/ -run TestContextBudget`
- **Tiempo estimado:** 45 min

#### Tarea 5: Config updates
- **Archivo:** Modificar [config.go](file:///home/daikyri/Workspace/Hackathon-QwenCloud/backend/internal/config/config.go)
- **Cambios:** Añadir `MaxContextTokens` (default 30000) y `ResponseReserve` (default 2000)
- **Archivo:** Modificar [.env.example](file:///home/daikyri/Workspace/Hackathon-QwenCloud/.env.example) — añadir las nuevas variables
- **Verificación:** `go build ./...`
- **Tiempo estimado:** 10 min

---

### FASE 2: Vector Search Mejorado (Tareas 6-8)

#### Tarea 6: VectorRepo — FindSimilarEntities (top-k)
- **Archivo:** Modificar [vector_repo.go](file:///home/daikyri/Workspace/Hackathon-QwenCloud/backend/internal/repositories/vector_repo.go)
- **Cambios:** Renombrar `FindSimilarEntity` → mantener por backward compat, añadir nuevo `FindSimilarEntities` que retorna `[]SimilarEntity` con limit parametrizable
- **Añadir:** `SetHNSWSearchParams` method
- **Añadir:** `KeywordSearch` method (usa tsvector)
- **Test:** `backend/internal/repositories/vector_repo_test.go` — test con múltiples embeddings insertados
- **Verificación:** `go test ./internal/repositories/ -run TestVectorRepo`
- **Tiempo estimado:** 1h

#### Tarea 7: Migración full-text search
- **Archivo:** Crear `backend/migrations/017_add_fulltext_index.up.sql` y `.down.sql`
- **Contenido:** ALTER TABLE + CREATE INDEX GIN de la sección 4.2.2
- **Verificación:** ejecutar migraciones, confirmar índice con `\di`
- **Tiempo estimado:** 15 min

#### Tarea 8: EntityService — usar FindSimilarEntities
- **Archivo:** Modificar [entity_service.go](file:///home/daikyri/Workspace/Hackathon-QwenCloud/backend/internal/services/entity_service.go)
- **Cambios:** En `ResolveOrCreate` step 3, usar `FindSimilarEntities` con limit 3 en lugar de 1, para mejor entity dedup. Comparar por score y si hay empate, usar Levenshtein distance
- **Verificación:** `go test ./internal/services/ -run TestEntityService`
- **Tiempo estimado:** 45 min

---

### FASE 3: Memory Consolidation (Tareas 9-13)

#### Tarea 9: Migración consolidated_memories
- **Archivo:** Crear `backend/migrations/016_create_consolidated_memories.up.sql` y `.down.sql`
- **Contenido:** Exactamente el SQL de la sección 3.2.1
- **Verificación:** ejecutar migraciones
- **Tiempo estimado:** 15 min

#### Tarea 10: ConsolidatedMemory model
- **Archivo:** Modificar [models.go](file:///home/daikyri/Workspace/Hackathon-QwenCloud/backend/internal/models/models.go)
- **Cambios:** Añadir struct `ConsolidatedMemory` de la sección 3.2.1
- **Verificación:** `go build ./...`
- **Tiempo estimado:** 10 min

#### Tarea 11: ConsolidationRepo
- **Archivo:** Crear `backend/internal/repositories/consolidation_repo.go`
- **Contenido:** CRUD + `FindSimilarByEmbedding` para búsqueda por cosine en consolidated memories
- **Test:** `backend/internal/repositories/consolidation_repo_test.go`
- **Verificación:** `go test ./internal/repositories/ -run TestConsolidationRepo`
- **Tiempo estimado:** 1h

#### Tarea 12: ConsolidationService
- **Archivo:** Crear `backend/internal/services/consolidation_service.go`
- **Contenido:** `ConsolidateEntity` y `DeconsolidateEntity` de la sección 3.2.3
- **Test:** `backend/internal/services/consolidation_service_test.go` — mock QwenService, verify consolidated memory is created
- **Verificación:** `go test ./internal/services/ -run TestConsolidation`
- **Tiempo estimado:** 1.5h

#### Tarea 13: Integrar consolidation en RelevanceService
- **Archivo:** Modificar [relevance_service.go](file:///home/daikyri/Workspace/Hackathon-QwenCloud/backend/internal/services/relevance_service.go)
- **Cambios:**
  - Añadir `consolidationSvc *ConsolidationService` al struct
  - Modificar `DecayAll` para llamar `ConsolidateEntity` para cada entidad recién archivada
  - Modificar `Reactivate` para llamar `DeconsolidateEntity`
  - Añadir `FindNewlyArchivable` al EntityRepo
- **Archivo:** Modificar [entity_repo.go](file:///home/daikyri/Workspace/Hackathon-QwenCloud/backend/internal/repositories/entity_repo.go) — añadir `FindNewlyArchivable`
- **Test:** Actualizar `relevance_service_test.go` y `entity_repo_test.go`
- **Verificación:** `go test ./internal/services/ -run TestRelevance`
- **Tiempo estimado:** 1h

---

### FASE 4: Hybrid Retrieval (Tareas 14-16)

#### Tarea 14: Reescribir MemoryService con RRF
- **Archivo:** Reescribir [memory_service.go](file:///home/daikyri/Workspace/Hackathon-QwenCloud/backend/internal/services/memory_service.go)
- **Contenido:** El MemoryService completo de la sección 4.2.1 con los 4 pipelines + RRF fusion
- **Nota:** La firma de `Recall` cambia — ahora acepta `queryText string` adicional y retorna `*BudgetReport`
- **Test:** Reescribir `memory_service_test.go` con mocks para cada pipeline
- **Verificación:** `go test ./internal/services/ -run TestMemoryService`
- **Tiempo estimado:** 2h

#### Tarea 15: Actualizar callers de Recall
- **Archivos a modificar:**
  - [analysis_service.go](file:///home/daikyri/Workspace/Hackathon-QwenCloud/backend/internal/services/analysis_service.go) línea 353: pasar `queryText` al `Recall`
  - [ws/hub.go](file:///home/daikyri/Workspace/Hackathon-QwenCloud/backend/internal/ws/hub.go) línea 326: actualizar firma de llamada a `Recall`
  - [handlers/graph.go](file:///home/daikyri/Workspace/Hackathon-QwenCloud/backend/internal/handlers/graph.go): actualizar handler de `/recall`
- **Verificación:** `go build ./...` y `go test ./...`
- **Tiempo estimado:** 1h

#### Tarea 16: Wiring en main.go
- **Archivo:** Modificar [main.go](file:///home/daikyri/Workspace/Hackathon-QwenCloud/backend/cmd/server/main.go)
- **Cambios:**
  - Crear `budgetMgr`
  - Crear `consolidationRepo` y `consolidationSvc`
  - Inyectar en `MemoryService`, `RelevanceService`, `AnalysisService`, `ContradictionService`
  - Añadir ruta `/universes/:id/memory-status`
- **Verificación:** `go build ./...` y el servidor arranca sin errores
- **Tiempo estimado:** 30 min

---

### FASE 5: Qwen API Avanzado (Tareas 17-20)

#### Tarea 17: Structured Output (response_format)
- **Archivo:** Modificar [qwen_service.go](file:///home/daikyri/Workspace/Hackathon-QwenCloud/backend/internal/services/qwen_service.go)
- **Cambios:**
  - Añadir `ResponseFormat` struct
  - Añadir campo `ResponseFormat` a `QwenRequest`
  - Añadir `ResponseFormat: &ResponseFormat{Type: "json_object"}` en `ExtractEntities`, `AnalyzeRelationships`, `CheckContradictions`
  - Usar `s.turboModel` y `s.maxModel` en lugar de strings hardcoded
- **Test:** Actualizar `qwen_service_test.go` — verificar que response_format aparece en el request body
- **Verificación:** `go test ./internal/services/ -run TestQwen`
- **Tiempo estimado:** 45 min

#### Tarea 18: Eliminar hack de "strip markdown fences"
- **Archivo:** Modificar [contradiction_service.go](file:///home/daikyri/Workspace/Hackathon-QwenCloud/backend/internal/services/contradiction_service.go)
- **Cambios:** Eliminar el bloque de fallback de strip markdown fences (líneas 219-228). Con structured output, Qwen siempre devuelve JSON válido
- **Verificación:** `go test ./internal/services/ -run TestContradiction`
- **Tiempo estimado:** 15 min

#### Tarea 19: Budget-aware Agent Loop
- **Archivo:** Modificar [qwen_service.go](file:///home/daikyri/Workspace/Hackathon-QwenCloud/backend/internal/services/qwen_service.go)
- **Cambios en `RunAgentLoop`:**
  - Aceptar `budgetMgr *ContextBudgetManager` como parámetro opcional
  - Después de cada iteración, contar tokens en `msgs`
  - Si `tokensUsed > 0.8 * maxTokens`: resumir tool results anteriores en un único mensaje comprimido
  - Retornar `BudgetReport` junto con la respuesta
- **Test:** Test con mock que simula 5 tool calls acumulando tokens
- **Verificación:** `go test ./internal/services/ -run TestRunAgentLoop`
- **Tiempo estimado:** 1.5h

#### Tarea 20: Integrar Budget en ContradictionService y PlotHoleService
- **Archivo:** Modificar [contradiction_service.go](file:///home/daikyri/Workspace/Hackathon-QwenCloud/backend/internal/services/contradiction_service.go)
- **Cambios en `CheckSemantic`:**
  - Usar budget manager para seleccionar top-k entidades que caben
  - Contar tokens del system prompt + user message antes de llamar al agent loop
  - Emitir BudgetReport
- **Archivo:** Modificar [plot_hole_service.go](file:///home/daikyri/Workspace/Hackathon-QwenCloud/backend/internal/services/plot_hole_service.go)
- **Cambios en `evaluatePlotHole`:**
  - Usar budget manager para limitar contexto
- **Verificación:** `go test ./internal/services/ -run TestCheckSemantic`
- **Tiempo estimado:** 1h

---

### FASE 6: Frontend — Memory Inspector (Tareas 21-27)

#### Tarea 21: memoryStore
- **Archivo:** Crear `frontend/src/stores/memoryStore.ts`
- **Contenido:** El store completo de la sección 6.2.2
- **Test:** `frontend/src/stores/__tests__/memoryStore.test.ts`
- **Verificación:** `npm run test -- --run src/stores/__tests__/memoryStore.test.ts`
- **Tiempo estimado:** 30 min

#### Tarea 22: wsStore — nuevos message handlers
- **Archivo:** Modificar [wsStore.ts](file:///home/daikyri/Workspace/Hackathon-QwenCloud/frontend/src/stores/wsStore.ts)
- **Cambios:** Añadir cases para `analysis_progress`, `budget_report`, `recall_sources` en el dispatch switch
- **Verificación:** `npm run test`
- **Tiempo estimado:** 20 min

#### Tarea 23: BudgetBar component
- **Archivo:** Crear `frontend/src/components/memory-inspector/BudgetBar.tsx`
- **Archivo:** Crear `frontend/src/components/memory-inspector/BudgetBar.module.css`
- **Diseño:** Barra horizontal segmentada con colores del design system:
  - System: `--ink-80`
  - User: `--gold`
  - Entities: `--character-purple` (#7A5C86)
  - Vector: `--place-green` (#6E8B4E)
  - Tools: `--event-orange` (#C07B3A)
- **Animación:** GSAP tween al actualizar anchos de segmento
- **Verificación:** Render manual en Storybook o en una página de test
- **Tiempo estimado:** 1h

#### Tarea 24: EntityLifecycleTimeline component
- **Archivo:** Crear `frontend/src/components/memory-inspector/EntityLifecycleTimeline.tsx`
- **Archivo:** Crear `frontend/src/components/memory-inspector/EntityLifecycleTimeline.module.css`
- **Diseño:**
  - Lista vertical de entidades, cada una con:
    - Nombre + tipo badge
    - Status icon (🟢/🟡/🔴/📦/🔄)
    - Mini sparkline SVG (inline, 80px wide) del score history
    - Click para expandir: consolidated summary si archived
  - Ordenar por status: active first, then decaying, then archived
- **Verificación:** Visual check
- **Tiempo estimado:** 1.5h

#### Tarea 25: RetrievalSourceBadges component
- **Archivo:** Crear `frontend/src/components/memory-inspector/RetrievalSourceBadges.tsx`
- **Diseño:** Para cada recall item, mostrar badges de color por fuente + RRF score
- **Verificación:** Visual check
- **Tiempo estimado:** 30 min

#### Tarea 26: PipelineProgress component
- **Archivo:** Crear `frontend/src/components/memory-inspector/PipelineProgress.tsx`
- **Diseño:** Stepper vertical con:
  - 5 steps: Extract → Contradict → Plot Holes → Recall → Complete
  - Cada step tiene spinner cuando running, check cuando done, cross cuando error
  - Duration en ms junto a cada step completado
- **Verificación:** Visual check
- **Tiempo estimado:** 45 min

#### Tarea 27: MemoryInspectorPanel (contenedor)
- **Archivo:** Crear `frontend/src/components/memory-inspector/MemoryInspectorPanel.tsx`
- **Archivo:** Crear `frontend/src/components/memory-inspector/MemoryInspector.module.css`
- **Diseño:** Panel lateral con 4 tabs (Budget / Entities / Sources / Pipeline)
- **Integración:** Importar en [EditorPage.tsx](file:///home/daikyri/Workspace/Hackathon-QwenCloud/frontend/src/pages/EditorPage.tsx) como sidebar derecho
- **Integración:** Importar en [WorkPage.tsx](file:///home/daikyri/Workspace/Hackathon-QwenCloud/frontend/src/pages/WorkPage.tsx) como panel colapsable
- **Verificación:** `npm run build` compila, visual check en browser
- **Tiempo estimado:** 1.5h

---

### FASE 7: Integration & Broadcast (Tareas 28-30)

#### Tarea 28: WS protocol updates (backend)
- **Archivo:** Modificar [ws/protocol.go](file:///home/daikyri/Workspace/Hackathon-QwenCloud/backend/internal/ws/protocol.go)
- **Cambios:** Añadir constantes: `TypeAnalysisProgress`, `TypeBudgetReport`, `TypeRecallSources`, `TypeMemoryStatus`
- **Verificación:** `go build ./...`
- **Tiempo estimado:** 10 min

#### Tarea 29: Analysis pipeline — emit progress + budget
- **Archivo:** Modificar [analysis_service.go](file:///home/daikyri/Workspace/Hackathon-QwenCloud/backend/internal/services/analysis_service.go)
- **Cambios en `processJob`:**
  - Emitir `analysis_progress` después de cada paso
  - Emitir `budget_report` con el BudgetReport del Recall
  - Emitir `recall_sources` con las fuentes RRF
- **Verificación:** `go test ./internal/services/ -run TestAnalysis`
- **Tiempo estimado:** 1h

#### Tarea 30: Memory status endpoint
- **Archivo:** Modificar [handlers/graph.go](file:///home/daikyri/Workspace/Hackathon-QwenCloud/backend/internal/handlers/graph.go)
- **Cambios:** Añadir handler `MemoryStatus` que retorne:
  ```json
  {
    "entities": [{"id": "...", "name": "...", "status": "...", "relevance_score": 0.75, ...}],
    "consolidated_count": 5,
    "active_count": 23,
    "archived_count": 12,
    "total_paragraph_embeddings": 450,
    "total_entity_embeddings": 35
  }
  ```
- **Verificación:** `curl http://localhost:8080/api/v1/universes/{id}/memory-status`
- **Tiempo estimado:** 45 min

---

## Resumen del Sprint

| Fase | Tareas | Tiempo Est. | Archivos Nuevos | Archivos Modificados |
|---|---|---|---|---|
| 1. Infraestructura | 1-5 | ~1.5h | 4 | 2 |
| 2. Vector Search | 6-8 | ~2h | 2 | 2 |
| 3. Consolidation | 9-13 | ~4h | 5 | 3 |
| 4. Hybrid Retrieval | 14-16 | ~3.5h | 0 | 4 |
| 5. Qwen API | 17-20 | ~3.5h | 1 | 3 |
| 6. Frontend | 21-27 | ~6h | 10 | 2 |
| 7. Integration | 28-30 | ~2h | 0 | 3 |
| **TOTAL** | **30 tareas** | **~22.5h** | **22 archivos** | **19 archivos** |

> [!TIP]
> Las fases 1-2 son fundacionales y deben hacerse primero. Las fases 3-5 pueden paralelizarse parcialmente. La fase 6 puede empezar en paralelo una vez la fase 1 esté lista (usa mocks). La fase 7 es integración final.

---

## Criterios atacados por fase

| Fase | Innovation (30%) | Technical Depth (30%) | Impact (25%) | Presentation (15%) |
|---|---|---|---|---|
| 1. Infraestructura | | ✅ | | |
| 2. Vector Search | | ✅✅ | | |
| 3. Consolidation | ✅ | ✅✅ | | |
| 4. Hybrid Retrieval | ✅ | ✅✅✅ | | |
| 5. Qwen API | ✅✅✅ | ✅ | | |
| 6. Frontend | | | | ✅✅✅ |
| 7. Integration | | ✅ | | ✅ |
