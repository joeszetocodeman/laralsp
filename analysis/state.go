package analysis

import (
	"educationalsp/lsp"
	"log"
	"strings"
)

type State struct {
	// Map of file names to contents
	Documents map[string]string
}

func NewState() State {
	return State{Documents: map[string]string{}}
}

func getDiagnosticsForFile(text string) []lsp.Diagnostic {
	diagnostics := []lsp.Diagnostic{}
	for row, line := range strings.Split(text, "\n") {
		if strings.Contains(line, "VS Code") {
			idx := strings.Index(line, "VS Code")
			diagnostics = append(diagnostics, lsp.Diagnostic{
				Range:    LineRange(row, idx, idx+len("VS Code")),
				Severity: 1,
				Source:   "Common Sense",
				Message:  "Please make sure we use good language in this video",
			})
		}

		if strings.Contains(line, "Neovim") {
			idx := strings.Index(line, "Neovim")
			diagnostics = append(diagnostics, lsp.Diagnostic{
				Range:    LineRange(row, idx, idx+len("Neovim")),
				Severity: 2,
				Source:   "Common Sense",
				Message:  "Great choice :)",
			})

		}
	}

	return diagnostics
}

func (s *State) OpenDocument(uri, text string) []lsp.Diagnostic {
	s.Documents[uri] = text

	return getDiagnosticsForFile(text)
}

func (s *State) UpdateDocument(uri, text string) []lsp.Diagnostic {
	s.Documents[uri] = text

	return getDiagnosticsForFile(text)
}

func (s *State) Hover(id int, uri string, position lsp.Position) lsp.HoverResponse {
	// In real life, this would look up the type in our type analysis code...

	// document := s.Documents[uri]

	return lsp.HoverResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: lsp.HoverResult{
			Contents: "",
		},
	}
}

func (s *State) Definition(id int, uri string, position lsp.Position) lsp.DefinitionResponse {
	// In real life, this would look up the definition

	return lsp.DefinitionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: lsp.Location{
			URI: uri,
			Range: lsp.Range{
				Start: lsp.Position{
					Line:      position.Line - 1,
					Character: 0,
				},
				End: lsp.Position{
					Line:      position.Line - 1,
					Character: 0,
				},
			},
		},
	}
}
func (s *State) TextDocumentCodeAction(logger *log.Logger, id int, params lsp.TextDocumentCodeActionParams) lsp.TextDocumentCodeActionResponse {
	uri := params.TextDocument.URI
	text := s.Documents[uri]
	hoverRange := params.Range
	hoverText := textInRange(text, hoverRange)

	actions := []lsp.CodeAction{}
	actions = append(actions, phpTraditionalFunctionToArrowCodeActions(uri, hoverText, hoverRange)...)
	actions = append(actions, phpArrowFunctionToTraditionalCodeActions(uri, hoverText, hoverRange)...)

	response := lsp.TextDocumentCodeActionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: actions,
	}

	return response
}

func (s *State) TextDocumentCompletion(id int, uri string) lsp.CompletionResponse {

	// Ask your static analysis tools to figure out good completions
	items := []lsp.CompletionItem{
		{
			Label:         "Neovim (BTW)",
			Detail:        "Very cool editor",
			Documentation: "Fun to watch in videos. Don't forget to like & subscribe to streamers using it :)",
		},
	}

	response := lsp.CompletionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: items,
	}

	return response
}

func LineRange(line, start, end int) lsp.Range {
	return lsp.Range{
		Start: lsp.Position{
			Line:      line,
			Character: start,
		},
		End: lsp.Position{
			Line:      line,
			Character: end,
		},
	}
}

func textInRange(text string, r lsp.Range) string {
	lines := strings.Split(text, "\n")
	if len(lines) == 0 {
		return ""
	}

	if r.Start.Line < 0 || r.End.Line < 0 || r.Start.Line >= len(lines) || r.End.Line >= len(lines) {
		return ""
	}
	if r.Start.Line > r.End.Line {
		return ""
	}

	if r.Start.Line == r.End.Line {
		line := lines[r.Start.Line]
		if r.Start.Character < 0 || r.End.Character < r.Start.Character || r.End.Character > len(line) {
			return ""
		}
		return line[r.Start.Character:r.End.Character]
	}

	var b strings.Builder
	for line := r.Start.Line; line <= r.End.Line; line++ {
		start := 0
		end := len(lines[line])
		if line == r.Start.Line {
			start = r.Start.Character
		}
		if line == r.End.Line {
			end = r.End.Character
		}
		if start < 0 || end < start || end > len(lines[line]) {
			return ""
		}

		if line > r.Start.Line {
			b.WriteString("\n")
		}
		b.WriteString(lines[line][start:end])
	}

	return b.String()
}
