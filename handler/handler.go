package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shukla2112/counter-app/config"
	"github.com/shukla2112/counter-app/databases"
	"github.com/shukla2112/counter-app/redis"
	"github.com/shukla2112/counter-app/utils"
)

// // NewObjFn : Endpoint handler function
// type NewObjFn func() interface{}

// // Func : Endpoint handler function
// type Func func(input interface{}, appC *config.AppConfig) (resp map[string]interface{}, err error)

// // EndpointHandler :
// func EndpointHandler(newFn NewObjFn, handlerFn Func, appC *config.AppConfig) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		input := newFn()
// 		if err := c.BindJSON(input); err != nil {
// 			utils.ReturnData(http.StatusBadRequest, err.Error(), 0, c)
// 			return
// 		}

// 		resp, err := handlerFn(input, appC)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "message": err.Error()})
// 			return
// 		}

// 		c.JSON(http.StatusOK, resp)
// 	}
// }

// GetCounter : Get the counter for the given key
func GetCounter(appC *config.AppConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		counterKey := c.Param("key")

		val, err := redis.RGet(counterKey, appC.RedisPool)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "message": err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{"status": 1, "value": val, "key": counterKey})
	}
}

// SetCounterInput :
type SetCounterInput struct {
	Operation  string `json:"operation"`
	Multiplier *int64 `json:"multiplier"`
	Key        string `json:"key"`
}

// SetCounter : set the counter for the given key
func SetCounter(appC *config.AppConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		var setCounterInput SetCounterInput
		// var err error
		var retval int64
		counterKey := c.Param("key")

		if err := c.BindJSON(&setCounterInput); err != nil {
			utils.ReturnData(http.StatusInternalServerError, err.Error(), 0, c)
			return
		}
		setCounterInput.Key = c.Param("key")

		fmt.Println(setCounterInput)

		if setCounterInput.Operation == "incr" {
			if setCounterInput.Multiplier != nil {
				if *setCounterInput.Multiplier > 0 {
					val, err := redis.RIncrby(setCounterInput.Key, *setCounterInput.Multiplier, appC.RedisPool)
					// fmt.Println(val)
					if err != nil {
						// fmt.Printf("HERE : %s\n", err.Error())
						utils.ReturnData(http.StatusInternalServerError, err.Error(), 0, c)
						return
					}
					retval = val
				} else {
					errMsg := fmt.Sprintf("Hey, multiplier is %d : multiplier should be > 0", *setCounterInput.Multiplier)
					log.Printf(errMsg)
					err := fmt.Errorf(errMsg)
					utils.ReturnData(http.StatusInternalServerError, err.Error(), 0, c)
					return
				}
			} else {
				val, err := redis.RIncr(setCounterInput.Key, appC.RedisPool)
				if err != nil {
					utils.ReturnData(http.StatusInternalServerError, err.Error(), 0, c)
					return
				}
				retval = val
			}
		} else if setCounterInput.Operation == "decr" {
			if setCounterInput.Multiplier != nil {
				if *setCounterInput.Multiplier > 1 {
					val, err := redis.RDecrby(setCounterInput.Key, *setCounterInput.Multiplier, appC.RedisPool)
					if err != nil {
						utils.ReturnData(http.StatusInternalServerError, err.Error(), 0, c)
						return
					}
					retval = val
				} else {
					errMsg := fmt.Sprintf("Hey, multiplier is %d : multiplier should be > 0", *setCounterInput.Multiplier)
					log.Printf(errMsg)
					err := fmt.Errorf(errMsg)
					utils.ReturnData(http.StatusInternalServerError, err.Error(), 0, c)
					return
				}
			} else {
				val, err := redis.RDecr(setCounterInput.Key, appC.RedisPool)
				if err != nil {
					utils.ReturnData(http.StatusInternalServerError, err.Error(), 0, c)
					return
				}
				retval = val
			}
		} else {
			errMsg := fmt.Sprintf("Hey, %s : Not a valid Operation, please provide 'incr' or 'decr'", setCounterInput.Operation)
			log.Printf(errMsg)
			err := fmt.Errorf(errMsg)
			utils.ReturnData(http.StatusInternalServerError, err.Error(), 0, c)
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": 1, "value": retval, "key": counterKey})
	}
}

// InitCounterInput :
type InitCounterInput struct {
	Key               string                      `json:"key"`
	ConnectionDetails databases.ConnectionDetails `json:"connection-details"`
	Query             string                      `json:"query"`
}

// InitCounter :
func InitCounter(appC *config.AppConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		var initCounterInput InitCounterInput

		if err := c.BindJSON(&initCounterInput); err != nil {
			utils.ReturnData(http.StatusInternalServerError, err.Error(), 0, c)
			return
		}
		initCounterInput.Key = c.Param("key")

		fmt.Printf("Input : %v\n", initCounterInput)

		// Establish a connection
		con, err := databases.EstConnection(initCounterInput.ConnectionDetails)
		defer con.Close()
		// fmt.Println(con)
		if err != nil {
			utils.ReturnData(http.StatusInternalServerError, err.Error(), 0, c)
			return
		}

		// Run the query - against datasource
		counterValue, err := databases.RunQuery(con, initCounterInput.Query)
		if err != nil {
			utils.ReturnData(http.StatusInternalServerError, err.Error(), 0, c)
			return
		}
		// fmt.Println(strconv.Itoa(*counterValue))

		// Set the redis key
		setVal, err := redis.RSet(initCounterInput.Key, strconv.Itoa(*counterValue), appC.RedisPool)
		if err != nil {
			utils.ReturnData(http.StatusInternalServerError, err.Error(), 0, c)
			return
		}
		fmt.Printf("Redis set : Key/Value : %s/%s\n", initCounterInput.Key, setVal)

		c.JSON(http.StatusOK, gin.H{"status": 1, "key": initCounterInput.Key, "value": counterValue})
	}
}
