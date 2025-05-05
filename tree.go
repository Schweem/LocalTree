package main

import (
	"strings"

	"github.com/rivo/tview"
)

// Helper to serialize the tree as Markdown
func serializeTreeMarkdown(node *tview.TreeNode, depth int) string {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}
	md := indent + "- " + node.GetText() + "\n"
	for _, child := range node.GetChildren() {
		md += serializeTreeMarkdown(child, depth+1)
	}
	return md
}

// Render the tree as text with tree-drawing characters
func renderTreeAsText(node *tview.TreeNode, prefix string, isRoot bool) string {
	var result string
	if isRoot {
		result += node.GetText() + "\n"
		children := node.GetChildren()
		n := len(children)
		for i, child := range children {
			var branch string
			if i == n-1 {
				branch = "└── "
			} else {
				branch = "├── "
			}
			result += branch + child.GetText() + "\n"
			result += renderTreeAsText(child, "    ", false)
		}
		return result
	}
	children := node.GetChildren()
	n := len(children)
	for i, child := range children {
		var branch, newPrefix string
		if i == n-1 {
			branch = prefix + "└── "
			newPrefix = prefix + "    "
		} else {
			branch = prefix + "├── "
			newPrefix = prefix + "│   "
		}
		result += branch + child.GetText() + "\n"
		result += renderTreeAsText(child, newPrefix, false)
	}
	return result
}

// Parse a tree from text lines (tree-drawing format)
func parseTreeFromText(lines []string) (*tview.TreeNode, map[*tview.TreeNode]*tview.TreeNode) {
	root := tview.NewTreeNode("Root").SetColor(tview.Styles.PrimaryTextColor)
	parentMap := make(map[*tview.TreeNode]*tview.TreeNode)
	type stackEntry struct {
		node  *tview.TreeNode
		depth int
	}
	stack := []stackEntry{{root, -1}}
	for _, line := range lines {
		// Remove any BOM, non-breaking spaces, or invisible characters
		line = strings.Map(func(r rune) rune {
			if r == '\uFEFF' || r == '\u200B' || r == '\u00A0' { // BOM, zero-width space, non-breaking space
				return -1
			}
			return r
		}, line)
		trimmed := strings.TrimLeft(line, " ")
		if trimmed == "" {
			continue
		}
		depth := (len(line) - len(trimmed)) / 4
		// Remove branch characters (handle possible corrupted variants)
		if strings.HasPrefix(trimmed, "├── ") || strings.HasPrefix(trimmed, "└── ") ||
			strings.HasPrefix(trimmed, "|-- ") || strings.HasPrefix(trimmed, "+-- ") {
			trimmed = trimmed[4:]
		}
		// Remove any leading replacement characters
		trimmed = strings.TrimLeft(trimmed, "�")
		node := tview.NewTreeNode(trimmed)
		for len(stack) > 0 && stack[len(stack)-1].depth >= depth {
			stack = stack[:len(stack)-1]
		}
		parent := stack[len(stack)-1].node
		parent.AddChild(node)
		parentMap[node] = parent
		stack = append(stack, stackEntry{node, depth})
	}
	return root, parentMap
}
