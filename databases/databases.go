package databases

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //mysql dialect
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// import (
// 	"fmt"
// 	"log"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/shukla2112/counter-app/config"
// 	"github.com/shukla2112/counter-app/redis"
// 	"github.com/shukla2112/counter-app/utils"
// )

// ConnectionDetails :
type ConnectionDetails struct {
	Datasource string `json:"datasource"`
	Host       string `json:"host"`
	Port       int64  `json:"port"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Dbname     string `json:"dbname"`
}

// EstConnection :
func EstConnection(connectionDetails ConnectionDetails) (db *gorm.DB, err error) {
	if connectionDetails.Datasource == "mysql" {
		connectionParams := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
			connectionDetails.Username,
			connectionDetails.Password,
			connectionDetails.Host,
			connectionDetails.Port,
			connectionDetails.Dbname,
		)
		db, err = gorm.Open("mysql", connectionParams)
		if err != nil {
			return db, fmt.Errorf("CONN_MYSQL_ERROR : couldn't establish the connection to mysql - error : %s", err.Error())
		}
		fmt.Println("CONN_MYSQL : You have connected to the database successfully")
		// defer db.Close()
	}

	if connectionDetails.Datasource == "postgres" {
		connectionParams := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
			connectionDetails.Host,
			connectionDetails.Port,
			connectionDetails.Username,
			connectionDetails.Dbname,
			connectionDetails.Password,
		)
		db, err = gorm.Open("postgres", connectionParams)
		if err != nil {
			return db, fmt.Errorf("CONN_PG_ERROR : couldn't establish the connection to postgres - error : %s", err.Error())
		}
		fmt.Println("CONN_PG : You have connected to the database successfully")
		// defer db.Close()
	}
	return db, err
}

// RunQueryResult :
type RunQueryResult struct {
	Count *int
}

// RunQuery :
func RunQuery(db *gorm.DB, query string) (count *int, err error) {
	var res RunQueryResult
	err = db.Debug().Raw(query).Scan(&res).Error
	if err != nil {
		return count, fmt.Errorf("QUERY_EXEC_ERROR : failed to run the query : error : %s", err.Error())
	}
	if res.Count == nil {
		return count, fmt.Errorf("QUERY_EXEC_ERROR : Query doesn't return count, please check the query before calling or mark the resulting column as count")
	}
	return res.Count, err
}
