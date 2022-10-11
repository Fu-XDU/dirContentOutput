package flags

import "gopkg.in/urfave/cli.v1"

var (
	Folder       string
	OutputFile   string
	IgnoreHidden bool

	FolderFlag = cli.StringFlag{
		Name:        "folder,f",
		Usage:       "Folder path",
		Destination: &Folder,
	}

	OutputFileFlag = cli.StringFlag{
		Name:        "output,o",
		Usage:       "Output file path",
		Value:       "Output.txt",
		Destination: &OutputFile,
	}

	IgnoreHiddenFile = cli.BoolFlag{
		Name:        "ignoreHidden,i",
		Usage:       "Ignore hidden files",
		Destination: &IgnoreHidden,
	}
)
