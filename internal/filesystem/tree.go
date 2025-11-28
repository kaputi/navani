package filesystem

import (
	"fmt"
	"strings"

	"github.com/kaputi/navani/internal/config"
	"github.com/kaputi/navani/internal/utils"
	"github.com/kaputi/navani/internal/utils/logger"
)

type TreeNode struct {
	name     string
	path     string
	ft       string
	isDir    bool
	open     bool
	children []*TreeNode
}

func NewFileTreeNode(name, path string, isDir bool) *TreeNode {
	ft, _ := utils.FTbyFileName(name)
	return &TreeNode{
		name:     name,
		path:     path,
		ft:       ft,
		isDir:    isDir,
		open:     false,
		children: []*TreeNode{},
	}
}

func (n *TreeNode) Name() string {
	return n.name
}

func (n *TreeNode) IsDir() bool {
	return n.isDir
}

func (n *TreeNode) Path() string {
	return n.path
}

func (n *TreeNode) FileType() string {
	return n.ft
}

func (n *TreeNode) IsOpen() bool {
	return n.open
}

func (n *TreeNode) Open() {
	if n.isDir {
		n.open = true
	}
}

func (n *TreeNode) Close() {
	if n.isDir {
		n.open = false
	}
}

func (n *TreeNode) Toggle() {
	if n.isDir {
		n.open = !n.open
	}
}

func (n *TreeNode) AddChild(child *TreeNode) {
	if n.isDir {
		n.children = append(n.children, child)
	} else {
		logger.Debug("Attempted to add a child to a non-directory node")
	}
}

func (n *TreeNode) Children() []*TreeNode {
	return n.children
}

type FileTree struct {
	Root         *TreeNode
	openNodeList []*TreeNode
	nodeList     []*TreeNode
}

func NewFileTree(rootPath string) *FileTree {
	ft := &FileTree{
		Root:         NewFileTreeNode("navani", rootPath, true),
		openNodeList: []*TreeNode{},
		nodeList:     []*TreeNode{},
	}

	ft.UpdateNodeLists()

	return ft
}

func walkNodes(currNode *TreeNode, nodeList []*TreeNode, all bool) []*TreeNode {
	nodeList = append(nodeList, currNode)
	// if is a dir and is open or i'm checking all nodes
	if currNode.isDir && currNode.open || all {
		for _, child := range currNode.children {
			nodeList = walkNodes(child, nodeList, all)
		}
	}
	return nodeList
}

func (ft *FileTree) UpdateNodeLists() {
	nodeList := walkNodes(ft.Root, []*TreeNode{}, true)
	ft.nodeList = nodeList
}

func (ft *FileTree) UpdateOpenNodeList() {
	nodeList := walkNodes(ft.Root, []*TreeNode{}, false)
	ft.openNodeList = nodeList
}

func (ft *FileTree) OpenNodeList() []*TreeNode {
	return append([]*TreeNode{}, ft.openNodeList...)
}

func (ft *FileTree) NodeList() []*TreeNode {
	return append([]*TreeNode{}, ft.nodeList...)
}

func recursiveStrings(currNode *TreeNode, strList []string, level int, last bool) []string {
	str := ""

	t := config.Theme()

	for range level {
		str += t.Tree.TreeIndentChar + strings.Repeat(" ", t.Tree.TreeIndentSize-1)
	}

	if last {
		str = utils.ReplaceLastOcurrence(str, t.Tree.TreeIndentChar, t.Tree.TreeLastIndentChar)
	}

	if currNode.isDir {
		if !last && len(str) > 0 {
			str = t.Tree.TreeDirIndentChar + str[len(t.Tree.TreeIndentChar):]
		}
		if currNode.open {
			str += t.Tree.TreeOpenChar
		} else {
			str += t.Tree.TreeCloseChar
		}

		icon := utils.GetFtIcon("directory")
		if len(currNode.children) == 0 {
			icon = utils.GetFtIcon("emptyDirectory")
		} else if currNode.open {
			icon = utils.GetFtIcon("openDirectory")
		}

		str += fmt.Sprintf("%s %s", icon, currNode.name)
		strList = append(strList, str)

		if currNode.open && len(currNode.children) > 0 {
			for i, child := range currNode.children {
				last := i == len(currNode.children)-1
				strList = recursiveStrings(child, strList, level+1, last)
			}
		}
	} else {
		icon := utils.GetFtIcon(currNode.ft)
		str += fmt.Sprintf("%s %s", icon, currNode.name)
		strList = append(strList, str)
	}

	return strList
}

func (ft *FileTree) Strings() []string {
	strList := recursiveStrings(ft.Root, []string{}, 0, false)
	return strList
}
