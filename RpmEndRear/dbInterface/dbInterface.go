package dbInterface

import (
	"FPMtestUpload/dataStruct"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/Go-SQL-Driver/MYSQL"
)

/*
Description:链接数据库的常量设置
*/
const (
	dbName           = "fpm"
	port             = "14809"
	ip               = "10.142.234.19"
	account          = "fpm"
	accountPwd       = "7c7837af25964aba"
	tableName        = "rpm_source_packages"
	source_id        = "source_id"
	source_name      = "source_name"
	source_version   = "source_version"
	source_iteration = "source_iteration"
	pre_installsh    = "pre_installsh"
	post_installsh   = "post_installsh"
	preun_installsh  = "preun_installsh"
	postun_installsh = "postun_installsh"
)

var dbObj *sql.DB

/*
function:连接数据库
output1:连接数据库失败，返回错误信息
output2:连接数据库成功，返回数据库对象
*/
func ConnectMysql() (err error) {
	dbObj, err = sql.Open("mysql", account+":"+accountPwd+"@tcp("+ip+":"+port+")/"+dbName+"?"+"charset=utf8")
	return
}

/*
function:向mysql数据库中插入数据
input:rpm包的基本信息结构体
output:插入失败，返回错误信息
*/
func InsertMysql(rpmInfo dataStruct.FpmBasicInfo) (err error) {
	stmt, err := dbObj.Prepare("INSERT " + tableName + " SET " + source_name +
		"=?," + source_version + "=?," + source_iteration + "=?," + pre_installsh + "=?," +
		post_installsh + "=?," + preun_installsh + "=?," + postun_installsh + "=?")
	fmt.Println("通过http向mysql数据库中插入数据")
	if err != nil {
		fmt.Println("通过http向mysql数据库中插入数据失败")
		return
	}

	res, err := stmt.Exec(rpmInfo.Name, rpmInfo.Version, rpmInfo.Iteration, rpmInfo.PreInstallSh, rpmInfo.PostInstallSh, rpmInfo.PreUninstallSh, rpmInfo.PostUninstallSh)

	if err != nil {
		fmt.Println(err)
		fmt.Println(res)
		return
	}
	return

}

/*
function:从mysql数据库中查询数据,回调函数，处理客户端查询数据的请求。
output:查询失败时返回错误信息

*/
func SelectMysql(response http.ResponseWriter, request *http.Request) {

	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	fileName := request.PostFormValue("n")
	fmt.Println("filename:" + fileName)

	stmt, err := dbObj.Prepare("select source_version,source_iteration from rpm_source_packages where source_name = ?")
	if err != nil {
		fmt.Println(err)
		return
	}
	rows, err := stmt.Query(fileName)

	defer rows.Close()
	var source_version string
	var source_iteration string
	for rows.Next() {

		rows.Scan(&source_version, &source_iteration)

	}

	var result dataStruct.ResponseMysqlInfo
	result.Data.Version = source_version
	result.Data.Iteration = source_iteration
	jsonResult, er := json.Marshal(result)
	if er == nil {
		response.Write(jsonResult)
	}

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(*rows)

	return
}
