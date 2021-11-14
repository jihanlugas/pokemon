package constant


type DocType uint8

const (
	FOTO_USER DocType = iota
)

var FolderForFile map[DocType]string

func init() {
	FolderForFile = make(map[DocType]string)
	FolderForFile[FOTO_USER] = "user"
}

