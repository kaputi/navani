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

// TODO: make directories sticky with config
func (n *TreeNode) recursiveStrings(strList []string, level int, last bool) []string {
	logger.Debug(fmt.Sprintf("Generating string for node: %s, level: %d, last: %t", n.name, level, last))

	str := ""

	for range int(level) {
		str += config.TreeIndentChar + strings.Repeat(" ", config.TreeIndentSize-1)
	}

	if last {
		str = utils.ReplaceLastOcurrence(str, config.TreeIndentChar, config.TreeLastIndentChar)
	}

	if n.isDir {
		if !last && len(str) > 0 {
			str = config.TreeDirIndentChar + str[len(config.TreeIndentChar):]
		}
		if n.open {
			str += config.TreeOpenChar
		} else {
			str += config.TreeCloseChar
		}
		icon := utils.GetFtIcon("directory")
		str += fmt.Sprintf("%s %s/", icon, n.name)
		strList = append(strList, str)

		if n.open && len(n.children) > 0 {
			for i, child := range n.children {
				last := i == len(n.children)-1
				strList = child.recursiveStrings(strList, level+1, last)
			}
		}
	} else {
		icon := utils.GetFtIcon(n.ft)
		str += fmt.Sprintf("%s %s", icon, n.name)
		strList = append(strList, str)
	}

	return strList
}

func (n *TreeNode) Strings() []string {
	strList := n.recursiveStrings([]string{}, 0, false)
	return strList
}
