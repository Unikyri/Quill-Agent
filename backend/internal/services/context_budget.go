package services

import (
	"sort"
	"strings"
)

// RankedItem is a piece of context text with a relevance score used to
// prioritize inclusion under a token budget.
type RankedItem struct {
	Text  string
	Score float64
}

// BudgetAllocation is the proportional token split computed by
// ContextBudgetManager.ComputeBudget.
type BudgetAllocation struct {
	EntitiesTokens int
	VectorTokens   int
	ToolsTokens    int
	Available      int
}

// BudgetReport is a serializable summary of a BudgetAllocation, suitable for
// sending over the WS progress channel (BudgetAllocation itself has no
// maxContext/percent-used view baked in).
type BudgetReport struct {
	MaxContextTokens int     `json:"max_context_tokens"`
	Available        int     `json:"available"`
	EntitiesTokens   int     `json:"entities_tokens"`
	VectorTokens     int     `json:"vector_tokens"`
	ToolsTokens      int     `json:"tools_tokens"`
	UsedPercent      float64 `json:"used_percent"`
	// VectorTokensUsed is how many tokens of the VectorTokens allocation the
	// fitted recall items actually consumed. Set by RecallExplain; zero for
	// reports built straight from an allocation.
	VectorTokensUsed int `json:"vector_tokens_used"`
}

// Report summarizes a into a BudgetReport against maxContext, the context
// window size the allocation was computed under. UsedPercent is the share of
// maxContext consumed by system/user/response overhead (maxContext -
// Available), floored at 0 to avoid division by zero when maxContext is 0.
func (a BudgetAllocation) Report(maxContext int) BudgetReport {
	var usedPercent float64
	if maxContext > 0 {
		used := maxContext - a.Available
		usedPercent = float64(used) / float64(maxContext) * 100
	}
	return BudgetReport{
		MaxContextTokens: maxContext,
		Available:        a.Available,
		EntitiesTokens:   a.EntitiesTokens,
		VectorTokens:     a.VectorTokens,
		ToolsTokens:      a.ToolsTokens,
		UsedPercent:      usedPercent,
	}
}

// ContextBudgetManager allocates a fixed context-window token budget across
// entities, vector memories, and tool results, and fits ranked items into a
// given budget.
type ContextBudgetManager struct {
	tok              *Tokenizer
	maxContextTokens int
	responseReserve  int
}

// NewContextBudgetManager builds a ContextBudgetManager bound to a tokenizer
// and the context window limits.
func NewContextBudgetManager(tok *Tokenizer, maxContextTokens, responseReserve int) *ContextBudgetManager {
	return &ContextBudgetManager{
		tok:              tok,
		maxContextTokens: maxContextTokens,
		responseReserve:  responseReserve,
	}
}

// ComputeBudget reserves systemTokens, userBaseTokens, and responseReserve
// from maxContextTokens, then splits the remainder 35% entities / 40%
// vector / 25% tools. Available and each split are floored at 0.
func (b *ContextBudgetManager) ComputeBudget(systemTokens, userBaseTokens int) BudgetAllocation {
	available := b.maxContextTokens - b.responseReserve - systemTokens - userBaseTokens
	if available < 0 {
		available = 0
	}

	return BudgetAllocation{
		EntitiesTokens: available * 35 / 100,
		VectorTokens:   available * 40 / 100,
		ToolsTokens:    available * 25 / 100,
		Available:      available,
	}
}

// FitToBudget sorts items by Score descending, then greedily includes items
// while staying within budget. It uses continue (not break) on overflow so
// a smaller later item can still fit after a larger one is skipped.
func (b *ContextBudgetManager) FitToBudget(items []RankedItem, budget int) (fitted []RankedItem, dropped, tokensUsed int) {
	if len(items) == 0 {
		return nil, 0, 0
	}

	sorted := make([]RankedItem, len(items))
	copy(sorted, items)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Score > sorted[j].Score })

	for _, item := range sorted {
		itemTokens := b.tok.CountTokens(item.Text)
		if tokensUsed+itemTokens > budget {
			dropped++
			continue
		}
		fitted = append(fitted, item)
		tokensUsed += itemTokens
	}
	return fitted, dropped, tokensUsed
}

// TruncateToTokens greedily accumulates text paragraph-by-paragraph
// (splitting on "\n\n") from the start until adding the next paragraph would
// exceed budget, preserving prose order (unlike FitToBudget, which sorts by
// score and would scramble a chapter's narrative flow).
//
// Pinned edge case: if the very first paragraph alone exceeds budget, it is
// still returned — truncated at the token level (not dropped to empty) —
// since giving the caller partial prose to analyze beats giving it nothing.
func (b *ContextBudgetManager) TruncateToTokens(text string, budget int) string {
	if budget <= 0 || text == "" {
		return ""
	}

	paragraphs := strings.Split(text, "\n\n")
	var buf strings.Builder
	used := 0
	for _, p := range paragraphs {
		pTokens := b.tok.CountTokens(p)
		if used == 0 && buf.Len() == 0 && pTokens > budget {
			return truncateParagraphToTokens(b.tok, p, budget)
		}
		if used+pTokens > budget {
			break
		}
		if buf.Len() > 0 {
			buf.WriteString("\n\n")
		}
		buf.WriteString(p)
		used += pTokens
	}
	return buf.String()
}

// truncateParagraphToTokens binary-searches the longest rune-safe prefix of
// text whose token count fits within budget.
//
// ponytail: O(log n) tokenizer calls — fine at chapter-text sizes; only
// reached for a single oversized paragraph, not the common path.
func truncateParagraphToTokens(tok *Tokenizer, text string, budget int) string {
	runes := []rune(text)
	lo, hi := 0, len(runes)
	for lo < hi {
		mid := (lo + hi + 1) / 2
		if tok.CountTokens(string(runes[:mid])) <= budget {
			lo = mid
		} else {
			hi = mid - 1
		}
	}
	return string(runes[:lo])
}
