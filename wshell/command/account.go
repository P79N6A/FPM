package command

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"client-put/model"
	//"furion-stable/fpmfshell/upgrade"
)

var FShellRootPath string

func Account(cmd string, params ...string) {
	//upgrade.VersionCheck()
	if len(params) == 0 {
		account, gErr := GetAccount()
		if gErr != nil {
			fmt.Println(gErr)
			os.Exit(1)
		}
		fmt.Println(account.String())
	} else if len(params) == 1 {
		name := params[0]
		sErr := SetAccount(name)
		if sErr != nil {
			fmt.Println(sErr)
			os.Exit(1)
		}
	} else {
		CmdHelp(cmd)
	}
}

func GetAccount() (account model.Acc, err error) {
	storageDir := filepath.Join(FShellRootPath, ".fshell")
	accountFname := filepath.Join(storageDir, "account.json")
	accountFh, openErr := os.Open(accountFname)
	if openErr != nil {
		err = fmt.Errorf("Open account file error, %s, please use `account` to set Name first", openErr)
		return
	}
	defer accountFh.Close()

	accountBytes, readErr := ioutil.ReadAll(accountFh)
	if readErr != nil {
		err = fmt.Errorf("Read account file error, %s", readErr)
		return
	}

	if umError := json.Unmarshal(accountBytes, &account); umError != nil {
		err = fmt.Errorf("Parse account file error, %s", umError)
		return
	}

	return
}

func SetAccount(name string) (err error) {
	storageDir := filepath.Join(FShellRootPath, ".fshell")
	if _, sErr := os.Stat(storageDir); sErr != nil {
		if mErr := os.MkdirAll(storageDir, 0755); mErr != nil {
			err = fmt.Errorf("Mkdir `%s` error, %s", storageDir, mErr)
			return
		}
	}

	accountFname := filepath.Join(storageDir, "account.json")

	accountFh, openErr := os.Create(accountFname)
	if openErr != nil {
		err = fmt.Errorf("Open account file error, %s", openErr)
		return
	}
	defer accountFh.Close()

	var account model.Acc
	account.Name = name

	jsonStr, mErr := account.ToJson()
	if mErr != nil {
		err = mErr
		return
	}
	_, wErr := accountFh.WriteString(jsonStr)
	if wErr != nil {
		err = fmt.Errorf("Write account info error, %s", wErr)
		return
	}

	return
}
