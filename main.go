package main

import (
	"dirContentOutput/flags"
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	clientIdentifier = "App" // Client identifier to advertise over the network
	clientVersion    = "1.0.0"
	clientUsage      = "App"
)

var (
	app       = cli.NewApp()
	baseFlags = []cli.Flag{
		flags.FolderFlag,
		flags.OutputFileFlag,
		flags.IgnoreHiddenFile,
	}

	f *os.File
)

func init() {
	app.Action = App
	app.Name = clientIdentifier
	app.Version = clientVersion
	app.Usage = clientUsage
	app.Commands = []cli.Command{}
	app.Flags = append(app.Flags, baseFlags...)
}

func App(ctx *cli.Context) error {
	if args := ctx.Args(); len(args) > 0 {
		return fmt.Errorf("invalid command: %q", args[0])
	}
	prepare()
	process()
	return nil
}

func prepare() {
	if flags.Folder == "" {
		log.Fatal("未指定要读取的文件夹")
	}
	log.Println("将要读取文件夹" + flags.Folder)
	if flags.IgnoreHidden {
		log.Println("忽略隐藏文件")
	} else {
		log.Println("不忽略隐藏文件")
	}
	log.Println("输出至文件" + flags.OutputFile)
}

func process() {
	var err error
	filePath := flags.OutputFile
	f, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal("文件打开失败", err)
	}
	//及时关闭file句柄
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			log.Fatal("文件关闭失败", err)
		}
	}(f)
	GetFiles(strings.TrimRight(flags.Folder, "/"))
}

func GetFiles(folder string) {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		log.Fatal("读取文件夹失败，error: ", err)
	}
	for _, file := range files {
		if flags.IgnoreHidden && strings.HasPrefix(file.Name(), ".") {
			continue
		}
		if file.IsDir() {
			GetFiles(folder + "/" + file.Name())
		} else {
			readFile(folder, file.Name())
		}
	}
}

func readFile(folder, filename string) {
	bytes, err := ioutil.ReadFile(folder + "/" + filename)
	if err != nil {
		log.Fatal("读取文件失败", err)
	}
	relateFilePath := strings.Replace(folder, flags.Folder+"/", "", 1) + "/" + filename
	filePath := fmt.Sprintf("/*********************  %s *********************/\n\n", relateFilePath)
	_, err = f.WriteString(filePath)
	if err != nil {
		log.Fatal("写入文件失败", err)
	}
	_, err = f.Write(bytes)
	if err != nil {
		log.Fatal("写入文件失败", err)
	}
	_, err = f.WriteString("\n")
	if err != nil {
		log.Fatal("写入文件失败", err)
	}
}

func main() {
	if err := app.Run(os.Args); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
