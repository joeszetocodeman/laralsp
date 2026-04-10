package analysis

import (
	"educationalsp/lsp"
	"regexp"
	"strings"
)

var phpTraditionalFunctionPattern = regexp.MustCompile(`(?s)^\s*(static\s+)?function\s*(&\s*)?\(([^)]*)\)\s*(:\s*[^{]+)?\s*(?:use\s*\(([^)]*)\)\s*)?\{\s*return\s+(.+?)\s*;\s*\}\s*$`)
var phpArrowFunctionPattern = regexp.MustCompile(`(?s)^\s*(static\s+)?fn\s*(&\s*)?\(([^)]*)\)\s*(:\s*[^=]+?)?\s*=>\s*(.+?)\s*$`)

func phpTraditionalFunctionToArrowCodeActions(uri string, selectedText string, selectedRange lsp.Range) []lsp.CodeAction {
	matches := phpTraditionalFunctionPattern.FindStringSubmatch(selectedText)
	if len(matches) == 0 {
		return nil
	}

	staticKeyword := strings.TrimSpace(matches[1])
	returnByReference := strings.TrimSpace(matches[2])
	parameters := strings.TrimSpace(matches[3])
	returnType := strings.TrimSpace(matches[4])
	expression := strings.TrimSpace(matches[6])

	newText := ""
	if staticKeyword != "" {
		newText += staticKeyword + " "
	}
	newText += "fn"
	if returnByReference != "" {
		newText += returnByReference
	}
	newText += "(" + parameters + ")"
	if returnType != "" {
		newText += " " + returnType
	}
	newText += " => " + expression

	changes := map[string][]lsp.TextEdit{
		uri: {
			{
				Range:   selectedRange,
				NewText: newText,
			},
		},
	}

	return []lsp.CodeAction{
		{
			Title: "Convert function to arrow function",
			Edit:  &lsp.WorkspaceEdit{Changes: changes},
		},
	}
}

func phpArrowFunctionToTraditionalCodeActions(uri string, selectedText string, selectedRange lsp.Range) []lsp.CodeAction {
	matches := phpArrowFunctionPattern.FindStringSubmatch(selectedText)
	if len(matches) == 0 {
		return nil
	}

	staticKeyword := strings.TrimSpace(matches[1])
	returnByReference := strings.TrimSpace(matches[2])
	parameters := strings.TrimSpace(matches[3])
	returnType := strings.TrimSpace(matches[4])
	expression := strings.TrimSpace(matches[5])

	newText := ""
	if staticKeyword != "" {
		newText += staticKeyword + " "
	}
	newText += "function "
	if returnByReference != "" {
		newText += returnByReference
	}
	newText += "(" + parameters + ")"
	if returnType != "" {
		newText += " " + returnType
	}
	newText += " { return " + expression + "; }"

	changes := map[string][]lsp.TextEdit{
		uri: {
			{
				Range:   selectedRange,
				NewText: newText,
			},
		},
	}

	return []lsp.CodeAction{
		{
			Title: "Convert arrow function to function",
			Edit:  &lsp.WorkspaceEdit{Changes: changes},
		},
	}
}
