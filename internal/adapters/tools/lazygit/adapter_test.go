package lazygit

import (
	"testing"

	"theme-engine/internal/domain/renderctx"
)

// TestReloadIsDocumentedNoOp locks the research
// conclusion: lazygit reads config.yml at startup
// only (issues #1158, #4602) and has no IPC or
// signal to re-read it, so Reload must succeed
// without doing anything — the engine's apply
// phase stays quiet instead of warning.
func TestReloadIsDocumentedNoOp(t *testing.T) {
	if err := (Processor{}).Reload("/nonexistent/config.yml", nil); err != nil {
		t.Fatalf("no-op Reload returned error: %v", err)
	}
}

// TestParseResolveNilSafe keeps the generic
// processor contract (Parse(nil) and Resolve(nil)
// are nil-safe) covered for this target too.
func TestParseResolveNilSafe(t *testing.T) {
	p := Processor{}

	if _, err := p.Parse(nil); err != nil {
		t.Errorf("Parse(nil) error: %v", err)
	}
	if _, err := p.Resolve(nil, &renderctx.Context{}); err != nil {
		t.Errorf("Resolve(nil) error: %v", err)
	}
}
