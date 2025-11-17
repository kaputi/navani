package filesystem

import (
	"fmt"
	"strings"

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
	n.open = true
}

func (n *TreeNode) Close() {
	n.open = false
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

// TODO: get characters from config
var (
	openChar   = "▼ "
	closeChar  = "▶ "
	indentChar = " "
	indentSize = 2
)

func (n *TreeNode) recursiveStrings(strList []string, level int) []string {
	if n.isDir {
		dirStr := strings.Repeat(indentChar, level*indentSize)
		if n.open {
			dirStr += openChar
		} else {
			dirStr += closeChar
		}
		icon := utils.GetFtIcon("directory")
		dirStr += fmt.Sprintf("%s %s/", icon, n.name)
		strList = append(strList, dirStr)

		if n.open && len(n.children) > 0 {
			for _, child := range n.children {
				strList = child.recursiveStrings(strList, level+1)
			}
		}
	} else {
		fileStr := strings.Repeat(indentChar, level*indentSize)
		// TODO: get file icon from config based on file type
		icon := utils.GetFtIcon(n.ft)
		fileStr += fmt.Sprintf("%s %s", icon, n.name)
		strList = append(strList, fileStr)
	}

	return strList
}

func (n *TreeNode) Strings() []string {
	strList := n.recursiveStrings([]string{}, 0)
	return strList
}
