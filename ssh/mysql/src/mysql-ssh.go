package src

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
)

type MysqlConfig struct {
	UserName,PassWd,IpPort,Dbname string
}

//MysqlClient -- ----------------------------
//--> @Description
//--> @Param
//--> @return
//-- ----------------------------
func MysqlClient(client Client,mysqlCfg *MysqlConfig)(db *sql.DB,err error){
	mysql.RegisterDialContext("mysql+tcp", (&ViaSSHDialer{client: client.client}).Dial)
	//mysql.RegisterDial("mysql+tcp", (&ViaSSHDialer{client.client}).Dial)
	//if db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@mysql+tcp(%s)/%s",mysqlCfg.UserName,mysqlCfg.PassWd,mysqlCfg.IpPort,mysqlCfg.Dbname));err != nil{
	//	return
	//}
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@mysql+tcp(%s)/%s",mysqlCfg.UserName,mysqlCfg.PassWd,mysqlCfg.IpPort,mysqlCfg.Dbname))

	return

	//if rows, err := db.Query("SELECT filed1,filed2 FROM table  limit 1"); err == nil {
	//	for rows.Next() {
	//		var arp_host_id int64
	//		var ip string
	//		rows.Scan(&arp_host_id, &ip)
	//		fmt.Printf("ID: %d Name: %s\n", arp_host_id, ip)
	//	}
	//	rows.Close()
	//} else {
	//	fmt.Printf("Failure: %s", err.Error())
	//}
	//
	//db.Close()
}
