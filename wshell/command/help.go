package command

import (
	"fmt"
	"os"
	"runtime"

	"client-put/model"
)

var optionDocs = map[string]string{
	"-v": "Show version",
	"-h": "Show help",
}

var cmds = []string{
	//"account",
	"init",
	"add",
	"delete",
	"list",
	"build",
	"put",
	"pack",
	"test",
	"upload",
}

var cmdDocs = map[string][]string{
	//"account": []string{"wshell account [<Name>]", "Get/Set account info."},
	"init":   []string{"wshell init <fileName-Version-Iteration>", "init directory"},
	"add":    []string{"wshell add <source>", "Add the source folder to the corresponding directory of the destination host"},
	"delete": []string{"wshell delete <fileName>", "Delete file. If there are no parameters,delete all files"},
	"list":   []string{"wshell list <fileName> <Level>", "fileName and level are Selectable"},
	"build":  []string{"wshell build <fileName-Version-Iteration>", "Package and compress a local file or directroy. If there are no parameters, compress files at default locations"},
	"put":    []string{"wshell put <fileName>", "Upload a local file."},
	"pack":   []string{"wshell pack packager -n name1 name2 name3... -d des1 des2 des3...", "Batch packing"},
	"test":   []string{"wshell test name1 name2 name3...", "Batch test"},
	"upload": []string{"wshell upload name1 name2 name3...", "Batch upload"},
}

func Help(cmd string, params ...string) {
	//upgrade.VersionCheck()
	if len(params) == 0 {
		fmt.Println(CmdList())
	} else {
		CmdHelps(params...)
	}
}

func CmdList() string {
	helpAll := fmt.Sprintf("wShell %s\r\n\r\n", model.Version)
	helpAll += "Options:\r\n"
	for k, v := range optionDocs {
		helpAll += fmt.Sprintf("\t%-20s%-20s\r\n", k, v)
	}
	helpAll += "\r\n"
	helpAll += "Commands:\r\n"
	for _, cmd := range cmds {
		if help, ok := cmdDocs[cmd]; ok {
			cmdDesc := help[1]
			helpAll += fmt.Sprintf("\t%-20s%-20s\r\n", cmd, cmdDesc)
		}
	}
	return helpAll
}

func CmdHelps(cmds ...string) {
	defer os.Exit(1)
	if len(cmds) == 0 {
		fmt.Println(CmdList())
	} else {
		for _, cmd := range cmds {
			CmdHelp(cmd)
		}
	}
}

func CmdHelp(cmd string) {
	docStr := fmt.Sprintf("Unknow cmd `%s`", cmd)
	if cmdDoc, ok := cmdDocs[cmd]; ok {
		docStr = fmt.Sprintf("Usage: %s\r\n  %s\r\n", cmdDoc[0], cmdDoc[1])
	}
	fmt.Println(docStr)
}

func Version() {
	fmt.Printf("wshell/%s (%s; %s; %s)\n", model.Version, runtime.GOOS, runtime.GOARCH, runtime.Version())
}
