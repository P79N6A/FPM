// client-put project doc.go

/*
client-put document
*/
package main

/*
上传rpm包的shell命令行
*/
/*
所有的方式都是post上传。
在buildrpm时，将网页版的信息利用shell命令直接上传，添加那几个参数就可以了。这样的话是实现一个包一个包
的打包，并不能实现批量打包（可不可以通过脚本来实现）。
在test时，获取测试包的地址，获取地址时需要上传包的名称
在up时，将包上传
*/
