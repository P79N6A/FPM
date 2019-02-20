// client-put project main.go
package main

import (
	//"client-put/dbInterface"
	"flag"
	"fmt"
	"os"
	//"os/user"

	"client-put/command"
	"client-put/model"
)

var supportedCmds = map[string]model.CliFunc{
	"account": command.Account,
	"put":     command.Put,
	"build":   command.Build,
	"add":     command.Add,
	"delete":  command.Delete,
	"init":    command.Init,
	"list":    command.List,
	"pack":    command.Create,
	"test":    command.TEST,
	"upload":  command.Upload,
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Use help or help [cmd1 [cmd2 [cmd3 ...]]] to see supported commands.")
		os.Exit(1)
	}
	var helpMode bool
	var versionMode bool
	flag.BoolVar(&helpMode, "h", false, "show help") //Bool()方式返回一个相应的指针，BoolVar()绑定到一个变量，Var()绑定自定义类型。
	flag.BoolVar(&versionMode, "v", false, "show version")
	flag.Parse()

	if helpMode {
		command.Help("help")
		return
	}

	if versionMode {
		command.Version()
		return
	}

	//设置命令、参数
	args := flag.Args()
	cmd := args[0]
	params := args[1:]

	if cliFunc, ok := supportedCmds[cmd]; ok {
		cliFunc(cmd, params...)
	} else {
		fmt.Printf("Error: unknown cmd `%s`\r\n", cmd)
		os.Exit(1)
	}
}
