package main

import (
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

	footer := tview.NewTextView().SetText("Enter: Add Child  |  Space: Expand/Collapse  |  r: Rename  |  d: Delete  |  q: Quit").SetTextAlign(tview.AlignCenter)
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
		return event
	})

	if err := app.SetRoot(layout, true).Run(); err != nil {
		panic(err)
	}
}
