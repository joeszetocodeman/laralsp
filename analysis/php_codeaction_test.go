package analysis

import (
	"educationalsp/lsp"
	"testing"
)

func TestPHPTraditionalFunctionToArrowCodeActions(t *testing.T) {
	uri := "file:///test.php"
	selected := "function ($x) { return $x + 1; }"
	rng := lsp.Range{
		Start: lsp.Position{Line: 0, Character: 0},
		End:   lsp.Position{Line: 0, Character: len(selected)},
	}

	actions := phpTraditionalFunctionToArrowCodeActions(uri, selected, rng)
	if len(actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(actions))
	}

	edits := actions[0].Edit.Changes[uri]
	if len(edits) != 1 {
		t.Fatalf("expected 1 edit, got %d", len(edits))
	}

	if edits[0].NewText != "fn($x) => $x + 1" {
		t.Fatalf("unexpected replacement text: %q", edits[0].NewText)
	}
}

func TestPHPTraditionalFunctionToArrowCodeActionsNoMatch(t *testing.T) {
	uri := "file:///test.php"
	selected := "function ($x) { $y = $x + 1; return $y; }"
	rng := lsp.Range{
		Start: lsp.Position{Line: 0, Character: 0},
		End:   lsp.Position{Line: 0, Character: len(selected)},
	}

	actions := phpTraditionalFunctionToArrowCodeActions(uri, selected, rng)
	if len(actions) != 0 {
		t.Fatalf("expected 0 actions, got %d", len(actions))
	}
}

func TestPHPArrowFunctionToTraditionalCodeActions(t *testing.T) {
	uri := "file:///test.php"
	selected := "fn($x): int => $x + 1"
	rng := lsp.Range{
		Start: lsp.Position{Line: 0, Character: 0},
		End:   lsp.Position{Line: 0, Character: len(selected)},
	}

	actions := phpArrowFunctionToTraditionalCodeActions(uri, selected, rng)
	if len(actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(actions))
	}

	edits := actions[0].Edit.Changes[uri]
	if len(edits) != 1 {
		t.Fatalf("expected 1 edit, got %d", len(edits))
	}

	if edits[0].NewText != "function ($x) : int { return $x + 1; }" {
		t.Fatalf("unexpected replacement text: %q", edits[0].NewText)
	}
}

func TestPHPArrowFunctionToTraditionalCodeActionsNoMatch(t *testing.T) {
	uri := "file:///test.php"
	selected := "function ($x) { return $x + 1; }"
	rng := lsp.Range{
		Start: lsp.Position{Line: 0, Character: 0},
		End:   lsp.Position{Line: 0, Character: len(selected)},
	}

	actions := phpArrowFunctionToTraditionalCodeActions(uri, selected, rng)
	if len(actions) != 0 {
		t.Fatalf("expected 0 actions, got %d", len(actions))
	}
}
