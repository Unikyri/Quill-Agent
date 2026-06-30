package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/quill/backend/internal/models"
)

// GraphNode represents a node returned from graph queries.
type GraphNode struct {
	ID           string                 `json:"id"`
	Labels       []string               `json:"labels"`
	Properties   map[string]interface{} `json:"properties"`
}

// GraphEdge represents an edge returned from graph queries.
type GraphEdge struct {
	ID         string                 `json:"id"`
	Source     string                 `json:"source"`
	Target     string                 `json:"target"`
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
}

type GraphRepo struct {
	pool *pgxpool.Pool
}

func NewGraphRepo(pool *pgxpool.Pool) *GraphRepo {
	return &GraphRepo{pool: pool}
}

func (r *GraphRepo) CreateGraph(ctx context.Context, universeID string) error {
	graphName := "universe_" + universeID
	query := `SELECT * FROM cypher($1, $$ CREATE (g:Graph {name: $2}) RETURN g $$) AS (g agtype)`
	_, err := r.pool.Exec(ctx, query, graphName, graphName)
	if err != nil {
		return fmt.Errorf("create graph: %w", err)
	}
	return nil
}

func (r *GraphRepo) CreateNode(ctx context.Context, graphName, label string, properties map[string]interface{}) error {
	query := fmt.Sprintf(`SELECT * FROM cypher($1, $$ CREATE (n:%s {entity_id: '%s', name: '%s', status: '%s', relevance_score: %v}) RETURN n $$) AS (n agtype)`,
		label,
		properties["entity_id"],
		properties["name"],
		properties["status"],
		properties["relevance_score"],
	)
	_, err := r.pool.Exec(ctx, query, graphName)
	if err != nil {
		return fmt.Errorf("create graph node: %w", err)
	}
	return nil
}

func (r *GraphRepo) CreateEdge(ctx context.Context, graphName, sourceEntityID, targetEntityID, relType string, properties map[string]interface{}) error {
	query := fmt.Sprintf(`SELECT * FROM cypher($1, $$ MATCH (x {entity_id: '%s'}), (y {entity_id: '%s'}) CREATE (x)-[:%s]->(y) $$) AS (r agtype)`,
		sourceEntityID, targetEntityID, relType)
	_, err := r.pool.Exec(ctx, query, graphName)
	if err != nil {
		return fmt.Errorf("create graph edge: %w", err)
	}
	return nil
}

func (r *GraphRepo) UpdateNodeRelevance(ctx context.Context, graphName, entityID string, score float64) error {
	query := fmt.Sprintf(`SELECT * FROM cypher($1, $$ MATCH (n {entity_id: '%s'}) SET n.relevance_score = %v RETURN n $$) AS (n agtype)`,
		entityID, score)
	_, err := r.pool.Exec(ctx, query, graphName)
	if err != nil {
		return fmt.Errorf("update node relevance: %w", err)
	}
	return nil
}

func (r *GraphRepo) GetNeighbors(ctx context.Context, graphName, entityID string) ([]models.GraphNeighbor, error) {
	query := fmt.Sprintf(`SELECT * FROM cypher($1, $$ MATCH (n {entity_id: '%s'})-[r]-(m) RETURN type(r) AS rel_type, properties(r) AS rel_props, m $$) AS (rel_type agtype, rel_props agtype, m agtype)`,
		entityID)
	rows, err := r.pool.Query(ctx, query, graphName)
	if err != nil {
		return nil, fmt.Errorf("get neighbors: %w", err)
	}
	defer rows.Close()

	var neighbors []models.GraphNeighbor
	for rows.Next() {
		var n models.GraphNeighbor
		if err := rows.Scan(&n.RelType, &n.RelProps, &n.Node); err != nil {
			return nil, fmt.Errorf("scan neighbor: %w", err)
		}
		neighbors = append(neighbors, n)
	}
	return neighbors, nil
}

// ponytail: improved FullQuery returns structured data instead of "graph data".
func (r *GraphRepo) FullQuery(ctx context.Context, graphName string) ([]GraphNode, []GraphEdge, error) {
	query := `SELECT * FROM cypher($1, $$ MATCH (n) OPTIONAL MATCH (n)-[r]->(m) RETURN n, r, m $$) AS (n agtype, r agtype, m agtype)`
	rows, err := r.pool.Query(ctx, query, graphName)
	if err != nil {
		return nil, nil, fmt.Errorf("full query: %w", err)
	}
	defer rows.Close()

	return collectGraphRows(rows)
}

// DeleteEdge removes a relationship between two nodes in the graph.
func (r *GraphRepo) DeleteEdge(ctx context.Context, graphName, sourceEntityID, targetEntityID, relType string) error {
	query := fmt.Sprintf(
		`SELECT * FROM cypher($1, $$ MATCH (x {entity_id: '%s'})-[r:%s]->(y {entity_id: '%s'}) DELETE r $$) AS (a agtype)`,
		sourceEntityID, relType, targetEntityID,
	)
	_, err := r.pool.Exec(ctx, query, graphName)
	if err != nil {
		return fmt.Errorf("delete edge: %w", err)
	}
	return nil
}

// NHopTraversal performs a BFS traversal from a start node up to `hops` depth,
// returning all nodes and edges discovered.
//
// ponytail: use AGE's variable-length patterns MATCH (n)-[*1..hops]-(m).
func (r *GraphRepo) NHopTraversal(ctx context.Context, graphName, startEntityID string, hops int) ([]GraphNode, []GraphEdge, error) {
	query := fmt.Sprintf(
		`SELECT * FROM cypher($1, $$ MATCH (n {entity_id: '%s'})-[r*1..%d]-(m) RETURN n, r, m $$) AS (n agtype, r agtype, m agtype)`,
		startEntityID, hops,
	)
	rows, err := r.pool.Query(ctx, query, graphName)
	if err != nil {
		return nil, nil, fmt.Errorf("n-hop traversal: %w", err)
	}
	defer rows.Close()

	return collectGraphRows(rows)
}

// collectGraphRows extracts nodes and edges from AGE cypher result rows.
// ponytail: shared helper for FullQuery and NHopTraversal — deduplication by entity_id.
func collectGraphRows(rows pgx.Rows) ([]GraphNode, []GraphEdge, error) {
	nodeMap := make(map[string]GraphNode)
	edgeMap := make(map[string]GraphEdge)

	for rows.Next() {
		var nStr, rStr, mStr *string
		if err := rows.Scan(&nStr, &rStr, &mStr); err != nil {
			return nil, nil, fmt.Errorf("scan row: %w", err)
		}
		if nStr != nil {
			id := extractProp(*nStr, "entity_id")
			if id != "" {
				if _, exists := nodeMap[id]; !exists {
					nodeMap[id] = GraphNode{ID: id, Properties: map[string]interface{}{"raw": *nStr}}
				}
			}
		}
		if mStr != nil {
			id := extractProp(*mStr, "entity_id")
			if id != "" {
				if _, exists := nodeMap[id]; !exists {
					nodeMap[id] = GraphNode{ID: id, Properties: map[string]interface{}{"raw": *mStr}}
				}
			}
		}
		if rStr != nil {
			key := *rStr
			if _, exists := edgeMap[key]; !exists {
				edgeMap[key] = GraphEdge{ID: key, Type: "relationship", Properties: map[string]interface{}{"raw": *rStr}}
			}
		}
	}

	nodes := make([]GraphNode, 0, len(nodeMap))
	for _, n := range nodeMap {
		nodes = append(nodes, n)
	}
	edges := make([]GraphEdge, 0, len(edgeMap))
	for _, e := range edgeMap {
		edges = append(edges, e)
	}
	return nodes, edges, nil
}

// extractProp pulls the entity_id value from a raw agtype string.
// ponytail: simple string extraction instead of full JSON parsing for agtype.
func extractProp(agtypeStr, key string) string {
	// agtype looks like: {"entity_id": "abc-123", ...}
	search := fmt.Sprintf(`"%s": "`, key)
	idx := strings.Index(agtypeStr, search)
	if idx < 0 {
		return ""
	}
	start := idx + len(search)
	end := strings.Index(agtypeStr[start:], `"`)
	if end < 0 {
		return ""
	}
	return agtypeStr[start : start+end]
}

func (r *GraphRepo) DropGraph(ctx context.Context, graphName string) error {
	query := `SELECT * FROM cypher($1, $$ MATCH (n) DETACH DELETE n $$) AS (a agtype)`
	_, err := r.pool.Exec(ctx, query, graphName)
	if err != nil {
		return fmt.Errorf("drop graph: %w", err)
	}
	return nil
}
