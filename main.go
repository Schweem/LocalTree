package main

import (
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	parentMap := make(map[*tview.TreeNode]*tview.TreeNode)

	rootNode := tview.NewTreeNode("Root").SetColor(tview.Styles.PrimaryTextColor)
	child1 := tview.NewTreeNode("Child 1")
	child2 := tview.NewTreeNode("Child 2")
	rootNode.AddChild(child1)
	rootNode.AddChild(child2)
	parentMap[child1] = rootNode
	parentMap[child2] = rootNode

	tree := tview.NewTreeView().SetRoot(rootNode).SetCurrentNode(rootNode)

	// Set background color to terminal default
	tree.SetBackgroundColor(tcell.ColorDefault)

	footer := tview.NewTextView().SetText("Enter: Add Child  |  Space: Expand/Collapse  |  r: Rename  |  d: Delete  |  Ctrl+S: Save  |  q: Quit").SetTextAlign(tview.AlignCenter)
	footer.SetBackgroundColor(tcell.ColorDefault)
	footer.SetDynamicColors(true)

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tree, 0, 1, true).
		AddItem(footer, 1, 1, false)

	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		// Always add a new child to the selected node
		newChild := tview.NewTreeNode("New Child")
		node.AddChild(newChild)
		parentMap[newChild] = node
	})

	tree.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		node := tree.GetCurrentNode()
		switch event.Rune() {
		case 'r': // Rename node
			input := tview.NewInputField().SetLabel("Rename: ").SetText(node.GetText())
			input.SetBackgroundColor(tcell.ColorDefault)
			input.SetDoneFunc(func(key tcell.Key) {
				if key == tcell.KeyEnter {
					node.SetText(input.GetText())
				}
				app.SetRoot(layout, true).SetFocus(tree)
			})
			modal := tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(input, 3, 1, true)
			app.SetRoot(modal, true).SetFocus(input)
			return nil
		case 'd': // Delete node (not root)
			if node != rootNode {
				parent := parentMap[node]
				if parent != nil {
					parent.RemoveChild(node)
					delete(parentMap, node)
				}
			}
			return nil
		case ' ': // Space toggles expand/collapse
			node.SetExpanded(!node.IsExpanded())
			return nil
		case 'q': // Quit
			app.Stop()
			return nil
		}
		if event.Key() == tcell.KeyCtrlS {
			// Save tree as a Markdown code block with tree-drawing characters
			treeText := renderTreeAsText(rootNode, "", true)
			md := "```text\n" + treeText + "```\n"
			file, err := os.Create("tree.md")
			if err == nil {
				file.WriteString(md)
				file.Close()
				footer.SetText("[green]Tree saved to tree.md[-]  |  Enter: Add Child  |  Space: Expand/Collapse  |  r: Rename  |  d: Delete  |  Ctrl+S: Save  |  q: Quit")
				go func() {
					app.QueueUpdateDraw(func() {
						footer.SetText("Enter: Add Child  |  Space: Expand/Collapse  |  r: Rename  |  d: Delete  |  Ctrl+S: Save  |  q: Quit")
					})
				}()
			}
			return nil
		}
		return event
	})

	if err := app.SetRoot(layout, true).Run(); err != nil {
		panic(err)
	}
}

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
