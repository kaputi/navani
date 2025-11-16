package filesystem

import (
	"fmt"
	"strings"

	"github.com/kaputi/navani/internal/utils"
	"github.com/kaputi/navani/internal/utils/logger"
)

type FileTreeNode struct {
	name     string
	path     string
	ft       string
	isDir    bool
	open     bool
	children []*FileTreeNode
}

func NewFileTreeNode(name, path string, isDir bool) *FileTreeNode {
	ft, _ := utils.FTbyFileName(name)
	return &FileTreeNode{
		name:     name,
		path:     path,
		ft:       ft,
		isDir:    isDir,
		open:     false,
		children: []*FileTreeNode{},
	}
}

func (n *FileTreeNode) Name() string {
	return n.name
}

func (n *FileTreeNode) IsDir() bool {
	return n.isDir
}

func (n *FileTreeNode) Path() string {
	return n.path
}

func (n *FileTreeNode) FileType() string {
	return n.ft
}

func (n *FileTreeNode) IsOpen() bool {
	return n.open
}

func (n *FileTreeNode) Open() {
	n.open = true
}

func (n *FileTreeNode) Close() {
	n.open = false
}

func (n *FileTreeNode) AddChild(child *FileTreeNode) {
	if n.isDir {
		n.children = append(n.children, child)
	} else {
		logger.Debug("Attempted to add a child to a non-directory node")
	}
}

func (n *FileTreeNode) Children() []*FileTreeNode {
	return n.children
}

// TODO: get characters from config
var (
	openChar   = "▼ "
	closeChar  = "▶ "
	indentChar = " "
	indentSize = 2
)

func (n *FileTreeNode) Strings(strList []string, level int) []string {

	if n.isDir {
		dirStr := strings.Repeat(indentChar, level*indentSize)
		if n.open {
			dirStr += openChar
		} else {
			dirStr += closeChar
		}
		icon := utils.GetFtIcon("directory")
		children := len(n.children)
		dirStr += fmt.Sprintf("%s %s (%v):", icon, n.name, children)
		strList = append(strList, dirStr)

		if n.open && len(n.children) > 0 {
			for _, child := range n.children {
				strList = child.Strings(strList, level+1)
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
