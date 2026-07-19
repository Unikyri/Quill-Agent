package repositories

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/quill/backend/internal/testutil"
)

// maliciousRelType/Label breaks out of the Cypher `[:%s]`/`n:%s` interpolation
// slot to attempt an unrelated DETACH DELETE — the injection class this guard
// closes.
const maliciousCypherIdentifier = `x}]->(n) DETACH DELETE n //`

// TestGraphRepoRejectsInvalidCypherIdentifiers verifies that CreateNode,
// CreateEdge, and DeleteEdge reject relType/label values that are not valid
// bare Cypher identifiers, returning ErrInvalidIdentifier and creating no
// graph row — and that legitimate identifiers still work unchanged.
func TestGraphRepoRejectsInvalidCypherIdentifiers(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	testutil.RunMigrationsUpTo(t, pool, "011")
	if !testutil.CheckAGE(t, pool) {
		t.Skip("Apache AGE extension not available; skipping graph-dependent test")
	}
	ctx := context.Background()
	repo := NewGraphRepo(pool)

	t.Run("malicious label rejected, no node created", func(t *testing.T) {
		uid := uuid.NewString()
		if err := repo.CreateGraph(ctx, uid); err != nil {
			t.Fatalf("CreateGraph: %v", err)
		}
		graphName := "universe_" + uid

		err := repo.CreateNode(ctx, graphName, maliciousCypherIdentifier, map[string]interface{}{
			"entity_id": uuid.NewString(), "name": "X", "status": "active", "relevance_score": 0.5,
		})
		if !errors.Is(err, ErrInvalidIdentifier) {
			t.Fatalf("CreateNode with malicious label: got err=%v, want ErrInvalidIdentifier", err)
		}

		nodes, _, qerr := repo.FullQuery(ctx, graphName)
		if qerr != nil {
			t.Fatalf("FullQuery: %v", qerr)
		}
		if len(nodes) != 0 {
			t.Errorf("expected 0 nodes created for rejected label, got %d", len(nodes))
		}
	})

	t.Run("malicious relType rejected, no edge created", func(t *testing.T) {
		graphName, e1, e2 := setupGraphTest(t, pool)

		err := repo.CreateEdge(ctx, graphName, e1, e2, maliciousCypherIdentifier, nil)
		if !errors.Is(err, ErrInvalidIdentifier) {
			t.Fatalf("CreateEdge with malicious relType: got err=%v, want ErrInvalidIdentifier", err)
		}

		// Only the ALLY_OF edge from setupGraphTest should exist — no node
		// deleted, no extra edge created by the injected fragment.
		_, edges, qerr := repo.NHopTraversal(ctx, graphName, e1, 1)
		if qerr != nil {
			t.Fatalf("NHopTraversal: %v", qerr)
		}
		if len(edges) != 1 {
			t.Errorf("expected 1 pre-existing edge (no injection side effects), got %d", len(edges))
		}
		nodes, _, qerr := repo.FullQuery(ctx, graphName)
		if qerr != nil {
			t.Fatalf("FullQuery: %v", qerr)
		}
		if len(nodes) != 2 {
			t.Errorf("expected 2 pre-existing nodes (DETACH DELETE must not have run), got %d", len(nodes))
		}
	})

	t.Run("malicious relType rejected by DeleteEdge, edge untouched", func(t *testing.T) {
		graphName, e1, e2 := setupGraphTest(t, pool)

		err := repo.DeleteEdge(ctx, graphName, e1, e2, maliciousCypherIdentifier)
		if !errors.Is(err, ErrInvalidIdentifier) {
			t.Fatalf("DeleteEdge with malicious relType: got err=%v, want ErrInvalidIdentifier", err)
		}

		_, edges, qerr := repo.NHopTraversal(ctx, graphName, e1, 1)
		if qerr != nil {
			t.Fatalf("NHopTraversal: %v", qerr)
		}
		if len(edges) != 1 {
			t.Errorf("expected existing ALLY_OF edge to survive rejected DeleteEdge, got %d edges", len(edges))
		}
	})

	t.Run("valid identifiers still succeed (regression)", func(t *testing.T) {
		uid := uuid.NewString()
		if err := repo.CreateGraph(ctx, uid); err != nil {
			t.Fatalf("CreateGraph: %v", err)
		}
		graphName := "universe_" + uid
		e1 := uuid.NewString()
		e2 := uuid.NewString()

		if err := repo.CreateNode(ctx, graphName, "Character", map[string]interface{}{
			"entity_id": e1, "name": "A", "status": "active", "relevance_score": 0.5,
		}); err != nil {
			t.Fatalf("CreateNode(Character): %v", err)
		}
		if err := repo.CreateNode(ctx, graphName, "Character", map[string]interface{}{
			"entity_id": e2, "name": "B", "status": "active", "relevance_score": 0.5,
		}); err != nil {
			t.Fatalf("CreateNode(Character): %v", err)
		}
		if err := repo.CreateEdge(ctx, graphName, e1, e2, "ALLY_OF", nil); err != nil {
			t.Fatalf("CreateEdge(ALLY_OF): %v", err)
		}

		nodes, edges, err := repo.FullQuery(ctx, graphName)
		if err != nil {
			t.Fatalf("FullQuery: %v", err)
		}
		if len(nodes) != 2 {
			t.Errorf("expected 2 nodes, got %d", len(nodes))
		}
		if len(edges) != 1 {
			t.Errorf("expected 1 edge, got %d", len(edges))
		}
	})
}

// TestGraphRepoWithAgeTxRestoresSearchPath is a regression test for
// search_path poisoning: withAgeTx sets search_path to ag_catalog first so
// AGE functions resolve, but must restore the pre-existing value before
// returning control to the caller — otherwise the pooled connection keeps
// resolving unqualified queries (e.g. "entities") against ag_catalog's
// internal tables instead of public.
func TestGraphRepoWithAgeTxRestoresSearchPath(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	testutil.RunMigrationsUpTo(t, pool, "011")
	if !testutil.CheckAGE(t, pool) {
		t.Skip("Apache AGE extension not available; skipping graph-dependent test")
	}

	ctx := context.Background()
	repo := NewGraphRepo(pool)

	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatalf("begin tx: %v", err)
	}
	defer tx.Rollback(ctx)

	c := tx.Conn()

	var before string
	if err := c.QueryRow(ctx, "SHOW search_path").Scan(&before); err != nil {
		t.Fatalf("show search_path (before): %v", err)
	}

	if err := repo.CreateGraphTx(ctx, tx, uuid.NewString()); err != nil {
		t.Fatalf("CreateGraphTx: %v", err)
	}

	var after string
	if err := c.QueryRow(ctx, "SHOW search_path").Scan(&after); err != nil {
		t.Fatalf("show search_path (after): %v", err)
	}

	if after != before {
		t.Errorf("search_path not restored: before=%q after=%q", before, after)
	}
	if strings.Contains(after, "ag_catalog") {
		t.Errorf("search_path still poisoned with ag_catalog: %q", after)
	}
}

func setupGraphTest(t *testing.T, pool *pgxpool.Pool) (string, string, string) {
	t.Helper()
	ctx := context.Background()
	repo := NewGraphRepo(pool)

	uid := uuid.NewString()
	if err := repo.CreateGraph(ctx, uid); err != nil {
		t.Fatalf("CreateGraph: %v", err)
	}
	graphName := "universe_" + uid

	e1 := uuid.NewString()
	e2 := uuid.NewString()

	// create nodes
	if err := repo.CreateNode(ctx, graphName, "Character", map[string]interface{}{
		"entity_id":       e1,
		"name":            "Alice",
		"status":          "active",
		"relevance_score": 0.8,
	}); err != nil {
		t.Fatalf("create node 1: %v", err)
	}
	if err := repo.CreateNode(ctx, graphName, "Character", map[string]interface{}{
		"entity_id":       e2,
		"name":            "Bob",
		"status":          "active",
		"relevance_score": 0.5,
	}); err != nil {
		t.Fatalf("create node 2: %v", err)
	}

	// create edge
	if err := repo.CreateEdge(ctx, graphName, e1, e2, "ALLY_OF", nil); err != nil {
		t.Fatalf("create edge: %v", err)
	}

	return graphName, e1, e2
}

// TestGraphRepoNHopTraversal verifies traversal up to N hops.
func TestGraphRepoNHopTraversal(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	testutil.RunMigrationsUpTo(t, pool, "011")
	if !testutil.CheckAGE(t, pool) {
		t.Skip("Apache AGE extension not available; skipping graph-dependent test")
	}
	graphName, e1, _ := setupGraphTest(t, pool)

	ctx := context.Background()
	repo := NewGraphRepo(pool)

	nodes, edges, err := repo.NHopTraversal(ctx, graphName, e1, 2)
	if err != nil {
		t.Fatalf("NHopTraversal: %v", err)
	}

	if len(nodes) < 2 {
		t.Errorf("NHopTraversal should return at least 2 nodes (start + neighbor), got %d", len(nodes))
	}
	if len(edges) < 1 {
		t.Errorf("NHopTraversal should return at least 1 edge, got %d", len(edges))
	}
	for _, edge := range edges {
		if edge.Type != "ALLY_OF" {
			t.Errorf("edge.Type = %q, want AGE relationship type ALLY_OF", edge.Type)
		}
	}
}

func TestGraphRepoBoundedNHopTraversalIncludesSecondHop(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	testutil.RunMigrationsUpTo(t, pool, "011")
	if !testutil.CheckAGE(t, pool) {
		t.Skip("Apache AGE extension not available; skipping graph-dependent test")
	}

	ctx := context.Background()
	repo := NewGraphRepo(pool)
	graphID := uuid.NewString()
	if err := repo.CreateGraph(ctx, graphID); err != nil {
		t.Fatalf("CreateGraph: %v", err)
	}
	graphName := "universe_" + graphID
	entityIDs := []string{uuid.NewString(), uuid.NewString(), uuid.NewString()}
	for _, entityID := range entityIDs {
		if err := repo.CreateNode(ctx, graphName, "Character", map[string]interface{}{
			"entity_id": entityID, "name": entityID, "status": "active", "relevance_score": 0.5,
		}); err != nil {
			t.Fatalf("CreateNode(%s): %v", entityID, err)
		}
	}
	if err := repo.CreateEdge(ctx, graphName, entityIDs[0], entityIDs[1], "KNOWS", nil); err != nil {
		t.Fatalf("CreateEdge(first hop): %v", err)
	}
	if err := repo.CreateEdge(ctx, graphName, entityIDs[1], entityIDs[2], "KNOWS", nil); err != nil {
		t.Fatalf("CreateEdge(second hop): %v", err)
	}

	result, err := repo.BoundedNHopTraversal(ctx, graphName, entityIDs[0], 2)
	if err != nil {
		t.Fatalf("BoundedNHopTraversal: %v", err)
	}
	if result.Truncated {
		t.Fatalf("small two-hop graph should not be truncated: %#v", result)
	}
	if result.Limits.Hops != 2 {
		t.Fatalf("applied hops = %d, want 2", result.Limits.Hops)
	}

	returnedIDs := make(map[string]struct{}, len(result.Nodes))
	for _, node := range result.Nodes {
		returnedIDs[node.ID] = struct{}{}
	}
	if _, ok := returnedIDs[entityIDs[2]]; !ok {
		t.Fatalf("second-hop entity %q missing from %#v", entityIDs[2], result.Nodes)
	}
}

// TestGraphRepoBoundedNHopTraversalBuildsRankedEgoGraph is an end-to-end
// check against real AGE for the ego-graph shape: only the top
// EgoGraphDegree1Limit direct neighbors by relevance_score survive at hop 1
// (a lower-relevance direct neighbor must be excluded), only the top
// EgoGraphDegree2PerSeedLimit of a hop-1 seed's own neighbors survive at hop
// 2, and every returned node carries the correct Properties["hop"].
func TestGraphRepoBoundedNHopTraversalBuildsRankedEgoGraph(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	testutil.RunMigrationsUpTo(t, pool, "011")
	if !testutil.CheckAGE(t, pool) {
		t.Skip("Apache AGE extension not available; skipping graph-dependent test")
	}

	ctx := context.Background()
	repo := NewGraphRepo(pool)
	graphID := uuid.NewString()
	if err := repo.CreateGraph(ctx, graphID); err != nil {
		t.Fatalf("CreateGraph: %v", err)
	}
	graphName := "universe_" + graphID

	focalID := uuid.NewString()
	if err := repo.CreateNode(ctx, graphName, "Character", map[string]interface{}{
		"entity_id": focalID, "name": "focal", "status": "active", "relevance_score": 1.0,
	}); err != nil {
		t.Fatalf("CreateNode(focal): %v", err)
	}

	// One more direct neighbor than EgoGraphDegree1Limit allows, ranked by a
	// distinct relevance_score per neighbor so the excluded one is unambiguous.
	directCount := EgoGraphDegree1Limit + 1
	directIDs := make([]string, directCount)
	for i := 0; i < directCount; i++ {
		id := uuid.NewString()
		directIDs[i] = id
		score := 0.9 - float64(i)*0.05 // strictly decreasing: directIDs[0] is most relevant
		if err := repo.CreateNode(ctx, graphName, "Character", map[string]interface{}{
			"entity_id": id, "name": fmt.Sprintf("direct-%d", i), "status": "active", "relevance_score": score,
		}); err != nil {
			t.Fatalf("CreateNode(direct-%d): %v", i, err)
		}
		if err := repo.CreateEdge(ctx, graphName, focalID, id, "KNOWS", nil); err != nil {
			t.Fatalf("CreateEdge(focal->direct-%d): %v", i, err)
		}
	}
	lowestRankedDirectID := directIDs[directCount-1]

	// Give the single top-ranked direct neighbor more second-hop neighbors
	// than EgoGraphDegree2PerSeedLimit allows.
	secondHopCount := EgoGraphDegree2PerSeedLimit + 1
	secondHopIDs := make([]string, secondHopCount)
	for i := 0; i < secondHopCount; i++ {
		id := uuid.NewString()
		secondHopIDs[i] = id
		score := 0.5 - float64(i)*0.05
		if err := repo.CreateNode(ctx, graphName, "Character", map[string]interface{}{
			"entity_id": id, "name": fmt.Sprintf("second-%d", i), "status": "active", "relevance_score": score,
		}); err != nil {
			t.Fatalf("CreateNode(second-%d): %v", i, err)
		}
		if err := repo.CreateEdge(ctx, graphName, directIDs[0], id, "KNOWS", nil); err != nil {
			t.Fatalf("CreateEdge(direct-0->second-%d): %v", i, err)
		}
	}
	lowestRankedSecondHopID := secondHopIDs[secondHopCount-1]

	result, err := repo.BoundedNHopTraversal(ctx, graphName, focalID, 2)
	if err != nil {
		t.Fatalf("BoundedNHopTraversal: %v", err)
	}

	hopByID := make(map[string]interface{}, len(result.Nodes))
	for _, node := range result.Nodes {
		hopByID[node.ID] = node.Properties["hop"]
	}

	if hop, ok := hopByID[focalID]; !ok || hop != 0 {
		t.Errorf("focal hop = %v (present=%v), want 0", hop, ok)
	}
	if hop, ok := hopByID[directIDs[0]]; !ok || hop != 1 {
		t.Errorf("top-ranked direct neighbor hop = %v (present=%v), want 1", hop, ok)
	}
	if _, ok := hopByID[lowestRankedDirectID]; ok {
		t.Errorf("lowest-ranked direct neighbor (rank %d of %d) should have been excluded by EgoGraphDegree1Limit=%d, but was present with hop=%v",
			directCount, directCount, EgoGraphDegree1Limit, hopByID[lowestRankedDirectID])
	}
	if hop, ok := hopByID[secondHopIDs[0]]; !ok || hop != 2 {
		t.Errorf("top-ranked second-hop neighbor hop = %v (present=%v), want 2", hop, ok)
	}
	if _, ok := hopByID[lowestRankedSecondHopID]; ok {
		t.Errorf("lowest-ranked second-hop neighbor should have been excluded by EgoGraphDegree2PerSeedLimit=%d, but was present with hop=%v",
			EgoGraphDegree2PerSeedLimit, hopByID[lowestRankedSecondHopID])
	}

	wantDegree1 := 0
	wantDegree2 := 0
	for _, hop := range hopByID {
		switch hop {
		case 1:
			wantDegree1++
		case 2:
			wantDegree2++
		}
	}
	if wantDegree1 != EgoGraphDegree1Limit {
		t.Errorf("degree-1 node count = %d, want %d", wantDegree1, EgoGraphDegree1Limit)
	}
	if wantDegree2 != EgoGraphDegree2PerSeedLimit { // only directIDs[0] was seeded with second-hop neighbors
		t.Errorf("degree-2 node count = %d, want %d", wantDegree2, EgoGraphDegree2PerSeedLimit)
	}
}

func TestBoundedGraphCollectorReportsNodeAndEdgeTruncation(t *testing.T) {
	t.Run("node limit", func(t *testing.T) {
		collector := newBoundedGraphCollector(2, 2)
		collector.addRows([]graphRow{
			graphRowForTest("focal", "first", "edge-1"),
			graphRowForTest("focal", "second", "edge-2"),
		})

		nodes, edges := collector.nodesAndEdges()
		if !collector.truncated {
			t.Fatal("expected node cap to mark the neighborhood as truncated")
		}
		if len(nodes) != 2 || len(edges) != 1 {
			t.Errorf("node-limited graph = %d nodes, %d edges; want 2 nodes, 1 edge", len(nodes), len(edges))
		}
	})

	t.Run("edge limit", func(t *testing.T) {
		collector := newBoundedGraphCollector(3, 1)
		collector.addRows([]graphRow{
			graphRowForTest("focal", "first", "edge-1"),
			graphRowForTest("focal", "second", "edge-2"),
		})

		nodes, edges := collector.nodesAndEdges()
		if !collector.truncated {
			t.Fatal("expected edge cap to mark the neighborhood as truncated")
		}
		if len(nodes) != 2 || len(edges) != 1 {
			t.Errorf("edge-limited graph = %d nodes, %d edges; want 2 nodes, 1 edge", len(nodes), len(edges))
		}
		if edges[0].Source != "focal" || edges[0].Target != "first" {
			t.Errorf("retained edge = %#v, want focal -> first", edges[0])
		}
	})
}

func TestGraphTraversalLimitsNormalizePublicHops(t *testing.T) {
	tests := []struct {
		input int
		want  int
	}{
		{input: -1, want: 1},
		{input: 0, want: 1},
		{input: 1, want: 1},
		{input: 2, want: 2},
		{input: 99, want: GraphTraversalMaxHops},
	}

	for _, tt := range tests {
		t.Run("hops", func(t *testing.T) {
			result := NewGraphTraversalResult(tt.input)
			if result.Limits.Hops != tt.want {
				t.Errorf("applied hops = %d, want %d", result.Limits.Hops, tt.want)
			}
			if result.Limits.MaxHops != GraphTraversalMaxHops || result.Limits.NodeLimit != GraphTraversalNodeLimit || result.Limits.EdgeLimit != GraphTraversalEdgeLimit || result.Limits.ResultLimit != GraphTraversalResultLimit {
				t.Errorf("unexpected traversal limits: %#v", result.Limits)
			}
		})
	}
}

func graphRowForTest(source, target, edgeID string) graphRow {
	sourceRaw := `{"properties":{"entity_id": "` + source + `"}}`
	targetRaw := `{"properties":{"entity_id": "` + target + `"}}`
	relationshipRaw := `{"id":"` + edgeID + `"}`
	relationshipType := `"KNOWS"`
	return graphRow{
		node:             &sourceRaw,
		relationship:     &relationshipRaw,
		target:           &targetRaw,
		relationshipType: &relationshipType,
	}
}

func graphRowForTestWithScore(source, target, edgeID string, score float64) graphRow {
	sourceRaw := `{"properties":{"entity_id": "` + source + `"}}`
	targetRaw := fmt.Sprintf(`{"properties":{"entity_id": "%s", "relevance_score": %v}}`, target, score)
	relationshipRaw := `{"id":"` + edgeID + `"}`
	relationshipType := `"KNOWS"`
	return graphRow{
		node:             &sourceRaw,
		relationship:     &relationshipRaw,
		target:           &targetRaw,
		relationshipType: &relationshipType,
	}
}

func TestExtractNumericProp(t *testing.T) {
	tests := []struct {
		raw    string
		key    string
		want   float64
		wantOk bool
	}{
		{`{"relevance_score": 0.65, "status": "archived"}`, "relevance_score", 0.65, true},
		{`{"relevance_score": 1}`, "relevance_score", 1, true},
		{`{"relevance_score": -0.5}`, "relevance_score", -0.5, true},
		{`{"status": "active"}`, "relevance_score", 0, false},
	}
	for _, tt := range tests {
		got, ok := extractNumericProp(tt.raw, tt.key)
		if ok != tt.wantOk || got != tt.want {
			t.Errorf("extractNumericProp(%q, %q) = (%v, %v), want (%v, %v)", tt.raw, tt.key, got, ok, tt.want, tt.wantOk)
		}
	}
}

func TestRankedDirectNeighborsOrdersByRelevanceDescendingWithIDTiebreak(t *testing.T) {
	rows := []graphRow{
		graphRowForTestWithScore("focal", "low", "e1", 0.2),
		graphRowForTestWithScore("focal", "high", "e2", 0.9),
		graphRowForTestWithScore("focal", "tie-b", "e3", 0.5),
		graphRowForTestWithScore("focal", "tie-a", "e4", 0.5),
	}
	neighbors := rankedDirectNeighbors(rows, "focal")
	order := make([]string, len(neighbors))
	for i, n := range neighbors {
		order[i] = n.node.ID
	}
	want := []string{"high", "tie-a", "tie-b", "low"}
	if !reflect.DeepEqual(order, want) {
		t.Errorf("rankedDirectNeighbors order = %v, want %v", order, want)
	}
}

func TestRankedDirectNeighborsGroupsMultipleEdgesToSameNeighbor(t *testing.T) {
	rows := []graphRow{
		graphRowForTestWithScore("focal", "n1", "e1", 0.5),
		graphRowForTestWithScore("focal", "n1", "e2", 0.5),
	}
	neighbors := rankedDirectNeighbors(rows, "focal")
	if len(neighbors) != 1 {
		t.Fatalf("expected 1 grouped neighbor, got %d", len(neighbors))
	}
	if len(neighbors[0].rows) != 2 {
		t.Errorf("expected both edges retained on the grouped neighbor, got %d", len(neighbors[0].rows))
	}
}

func TestRankedDirectNeighborsExcludesTheReferenceNodeItself(t *testing.T) {
	rows := []graphRow{graphRowForTestWithScore("focal", "focal", "e1", 0.9)}
	neighbors := rankedDirectNeighbors(rows, "focal")
	if len(neighbors) != 0 {
		t.Errorf("expected a self-referential row to be excluded, got %v", neighbors)
	}
}

func TestAddRankedNeighborsCapsAtLimitAndReturnsAddedIDsInRankOrder(t *testing.T) {
	collector := newBoundedGraphCollector(0, 0)
	rows := []graphRow{
		graphRowForTestWithScore("focal", "c", "e3", 0.7),
		graphRowForTestWithScore("focal", "a", "e1", 0.9),
		graphRowForTestWithScore("focal", "b", "e2", 0.8),
	}
	neighbors := rankedDirectNeighbors(rows, "focal")

	added := addRankedNeighbors(collector, neighbors, 2)
	if !reflect.DeepEqual(added, []string{"a", "b"}) {
		t.Errorf("added = %v, want [a b]", added)
	}

	nodes, edges := collector.nodesAndEdges()
	if len(nodes) != 3 { // focal (as the reference node on every added row) + a + b
		t.Errorf("expected 3 collected nodes (focal, a, b), got %d: %#v", len(nodes), nodes)
	}
	if len(edges) != 2 {
		t.Errorf("expected 2 collected edges, got %d", len(edges))
	}
	if collector.hasNode("c") {
		t.Error("the third-ranked neighbor should not have been added past the limit")
	}
}

func TestGraphRepoNHopTraversalIncludesIsolatedFocal(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	testutil.RunMigrationsUpTo(t, pool, "011")
	if !testutil.CheckAGE(t, pool) {
		t.Skip("Apache AGE extension not available; skipping graph-dependent test")
	}

	ctx := context.Background()
	repo := NewGraphRepo(pool)
	graphID := uuid.NewString()
	if err := repo.CreateGraph(ctx, graphID); err != nil {
		t.Fatalf("CreateGraph: %v", err)
	}

	entityID := uuid.NewString()
	graphName := "universe_" + graphID
	if err := repo.CreateNode(ctx, graphName, "Character", map[string]interface{}{
		"entity_id":       entityID,
		"name":            "Isolated",
		"status":          "active",
		"relevance_score": 0.8,
	}); err != nil {
		t.Fatalf("CreateNode: %v", err)
	}

	nodes, edges, err := repo.NHopTraversal(ctx, graphName, entityID, 2)
	if err != nil {
		t.Fatalf("NHopTraversal: %v", err)
	}
	if len(nodes) != 1 || nodes[0].ID != entityID {
		t.Fatalf("expected isolated focal node %q, got %#v", entityID, nodes)
	}
	if len(edges) != 0 {
		t.Fatalf("expected no edges for isolated focal node, got %#v", edges)
	}
}

// TestGraphRepoDeleteEdge removes an edge between two nodes.
func TestGraphRepoDeleteEdge(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	testutil.RunMigrationsUpTo(t, pool, "011")
	if !testutil.CheckAGE(t, pool) {
		t.Skip("Apache AGE extension not available; skipping graph-dependent test")
	}
	graphName, e1, e2 := setupGraphTest(t, pool)

	ctx := context.Background()
	repo := NewGraphRepo(pool)

	if err := repo.DeleteEdge(ctx, graphName, e1, e2, "ALLY_OF"); err != nil {
		t.Fatalf("DeleteEdge: %v", err)
	}

	// Traverse again — edge should be gone
	_, edges, err := repo.NHopTraversal(ctx, graphName, e1, 1)
	if err != nil {
		t.Fatalf("NHopTraversal after delete: %v", err)
	}
	if len(edges) != 0 {
		t.Errorf("expected 0 edges after DeleteEdge, got %d", len(edges))
	}
}

// TestGraphRepoUpdateEdge modifies an edge's relationship type.
func TestGraphRepoUpdateEdge(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	testutil.RunMigrationsUpTo(t, pool, "011")
	if !testutil.CheckAGE(t, pool) {
		t.Skip("Apache AGE extension not available; skipping graph-dependent test")
	}
	graphName, e1, e2 := setupGraphTest(t, pool)

	ctx := context.Background()
	repo := NewGraphRepo(pool)

	// The existing CreateEdge created ALLY_OF. Delete it and recreate as ENEMY_OF.
	// ponytail: update an edge = delete old + create new; AGE doesn't have SET on edges easily
	if err := repo.DeleteEdge(ctx, graphName, e1, e2, "ALLY_OF"); err != nil {
		t.Fatalf("delete old edge: %v", err)
	}
	if err := repo.CreateEdge(ctx, graphName, e1, e2, "ENEMY_OF", nil); err != nil {
		t.Fatalf("create new edge: %v", err)
	}

	_, edges, err := repo.NHopTraversal(ctx, graphName, e1, 1)
	if err != nil {
		t.Fatalf("traverse after update: %v", err)
	}
	if len(edges) != 1 {
		t.Fatalf("expected 1 edge after update, got %d", len(edges))
	}
}

// TestGraphRepoFullQuery returns structured graph data.
func TestGraphRepoFullQuery(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	testutil.RunMigrationsUpTo(t, pool, "011")
	if !testutil.CheckAGE(t, pool) {
		t.Skip("Apache AGE extension not available; skipping graph-dependent test")
	}
	graphName, _, _ := setupGraphTest(t, pool)

	ctx := context.Background()
	repo := NewGraphRepo(pool)

	nodes, edges, err := repo.FullQuery(ctx, graphName)
	if err != nil {
		t.Fatalf("FullQuery: %v", err)
	}

	if len(nodes) < 2 {
		t.Errorf("FullQuery should return at least 2 nodes, got %d", len(nodes))
	}
	if len(edges) < 1 {
		t.Errorf("FullQuery should return at least 1 edge, got %d", len(edges))
	}

	// Verify nodes have data
	for _, n := range nodes {
		if n.ID == "" {
			t.Error("node.ID should not be empty")
		}
	}

	// Verify edges have data
	for _, e := range edges {
		if e.Source == "" || e.Target == "" {
			t.Error("edge Source and Target should not be empty")
		}
	}
}

// TestGraphRepoUpdateNodeState keeps the graph-visible properties aligned
// with the canonical SQL relevance/lifecycle state.
func TestGraphRepoUpdateNodeState(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	testutil.RunMigrationsUpTo(t, pool, "011")
	if !testutil.CheckAGE(t, pool) {
		t.Skip("Apache AGE extension not available; skipping graph-dependent test")
	}
	graphName, entityID, _ := setupGraphTest(t, pool)
	repo := NewGraphRepo(pool)
	if err := repo.UpdateNodeState(context.Background(), graphName, entityID, 0.65, "archived"); err != nil {
		t.Fatalf("UpdateNodeState: %v", err)
	}
	nodes, _, err := repo.FullQuery(context.Background(), graphName)
	if err != nil {
		t.Fatalf("FullQuery: %v", err)
	}
	for _, node := range nodes {
		if node.ID != entityID {
			continue
		}
		raw, _ := node.Properties["raw"].(string)
		if !strings.Contains(raw, `"relevance_score": 0.65`) || !strings.Contains(raw, `"status": "archived"`) {
			t.Fatalf("graph node properties = %s, want current score and status", raw)
		}
		return
	}
	t.Fatalf("updated graph node %s not returned", entityID)
}

// TestGraphRepoCompose tests creating two edges between the same pair to verify compose works.
func TestGraphRepoEdgeCRUD(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	testutil.RunMigrationsUpTo(t, pool, "011")
	if !testutil.CheckAGE(t, pool) {
		t.Skip("Apache AGE extension not available; skipping graph-dependent test")
	}
	ctx := context.Background()
	repo := NewGraphRepo(pool)
	uid := uuid.NewString()
	repo.CreateGraph(ctx, uid)
	graphName := "universe_" + uid

	e1 := uuid.NewString()
	e2 := uuid.NewString()

	// Create nodes
	repo.CreateNode(ctx, graphName, "Character", map[string]interface{}{"entity_id": e1, "name": "X", "status": "active", "relevance_score": 0.5})
	repo.CreateNode(ctx, graphName, "Character", map[string]interface{}{"entity_id": e2, "name": "Y", "status": "active", "relevance_score": 0.5})

	// Create edge
	if err := repo.CreateEdge(ctx, graphName, e1, e2, "KNOWS", map[string]interface{}{"since": "ch1"}); err != nil {
		t.Fatalf("create edge: %v", err)
	}

	// Delete it
	if err := repo.DeleteEdge(ctx, graphName, e1, e2, "KNOWS"); err != nil {
		t.Fatalf("delete edge: %v", err)
	}

	// Verify deletion
	_, edges, _ := repo.NHopTraversal(ctx, graphName, e1, 1)
	if len(edges) != 0 {
		t.Errorf("expected 0 edges after delete, got %d", len(edges))
	}

	// Recreate with different type
	if err := repo.CreateEdge(ctx, graphName, e1, e2, "ENEMY_OF", nil); err != nil {
		t.Fatalf("recreate edge: %v", err)
	}

	_, edges2, _ := repo.NHopTraversal(ctx, graphName, e1, 1)
	if len(edges2) != 1 {
		t.Errorf("expected 1 edge after recreate, got %d", len(edges2))
	}
}

// TestGraphRepoGetNeighborsBatch verifies that neighbors for multiple seed
// entities are resolved via a single batched Cypher call (spec: "Graph
// Pipeline Uses Batched Neighbor Lookup"), not one GetNeighbors call per
// seed. Correctness of the returned per-seed map is what this test asserts;
// the "single call" property is enforced by GetNeighborsBatch's
// implementation issuing exactly one cypher() query regardless of len(seeds).
func TestGraphRepoGetNeighborsBatch(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	testutil.RunMigrationsUpTo(t, pool, "011")
	if !testutil.CheckAGE(t, pool) {
		t.Skip("Apache AGE extension not available; skipping graph-dependent test")
	}
	ctx := context.Background()
	repo := NewGraphRepo(pool)
	uid := uuid.NewString()
	if err := repo.CreateGraph(ctx, uid); err != nil {
		t.Fatalf("CreateGraph: %v", err)
	}
	graphName := "universe_" + uid

	e1 := uuid.NewString()
	e2 := uuid.NewString()
	e3 := uuid.NewString()

	for _, e := range []string{e1, e2, e3} {
		if err := repo.CreateNode(ctx, graphName, "Character", map[string]interface{}{
			"entity_id": e, "name": "N" + e[:4], "status": "active", "relevance_score": 0.5,
		}); err != nil {
			t.Fatalf("create node %s: %v", e, err)
		}
	}

	// e1 and e2 both know e3 — e3 should show up as a neighbor of both seeds.
	if err := repo.CreateEdge(ctx, graphName, e1, e3, "KNOWS", nil); err != nil {
		t.Fatalf("create edge e1-e3: %v", err)
	}
	if err := repo.CreateEdge(ctx, graphName, e2, e3, "KNOWS", nil); err != nil {
		t.Fatalf("create edge e2-e3: %v", err)
	}

	result, err := repo.GetNeighborsBatch(ctx, graphName, []string{e1, e2})
	if err != nil {
		t.Fatalf("GetNeighborsBatch: %v", err)
	}

	if len(result[e1]) != 1 {
		t.Errorf("expected 1 neighbor for e1, got %d", len(result[e1]))
	}
	if len(result[e2]) != 1 {
		t.Errorf("expected 1 neighbor for e2, got %d", len(result[e2]))
	}
}

// TestGraphRepoGetNeighborsBatchEmpty verifies the empty-input short-circuit
// (no seeds → no query, mirroring EntityIDsForParagraphs' convention).
func TestGraphRepoGetNeighborsBatchEmpty(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	testutil.RunMigrationsUpTo(t, pool, "011")
	repo := NewGraphRepo(pool)

	result, err := repo.GetNeighborsBatch(context.Background(), "universe_"+uuid.NewString(), nil)
	if err != nil {
		t.Fatalf("GetNeighborsBatch(empty): %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty map for empty seeds, got %d entries", len(result))
	}
}

// TestEscapeCypherString verifies that the escapeCypherString helper
// correctly escapes single quotes and backslashes for safe Cypher interpolation.
//
// RED: escapeCypherString does not exist yet — compilation will fail until
// the production code is added.
func TestEscapeCypherString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"single_quote", "O'Brien", "O\\'Brien"},
		{"backslash", "path\\to\\file", "path\\\\to\\\\file"},
		{"normal", "Alice", "Alice"},
		{"empty", "", ""},
		{"mixed", "O'Brien\\Jr", "O\\'Brien\\\\Jr"},
		{"only_quote", "'", "\\'"},
		{"only_backslash", "\\", "\\\\"},
		{"already_escaped", "already\\'safe", "already\\\\\\'safe"}, // ponytail: double-escape is harmless — better safe than injection
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := escapeCypherString(tt.input)
			if got != tt.want {
				t.Errorf("escapeCypherString(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
