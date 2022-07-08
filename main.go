package main
import (
"fmt"
"net/http"
"os"
"time"
"github.com/gin-gonic/gin"
mgo "gopkg.in/mgo.v2"
"gopkg.in/mgo.v2/bson"
)
// database and collection names are statically declared
const database, collection = "go-mongo-practice", "user"
// User structure
type User struct {
ID                bson.ObjectId `bson:"_id"`
Name          string         `bson:"name"`
Address       string         `bson:"address"`
Age              int       `bson:"age"`
CreatedAt    time.Time     `bson:"created_at"`
UpdatedAt   time.Time       `bson:"updated_at"`
}
// DB connection
func connect() *mgo.Session {
session, err := mgo.Dial("localhost")
if err != nil {
fmt.Println("session err:", err)
os.Exit(1)
}
return session
}
// StartService function
func StartService() {
router := gin.Default()
api := router.Group("/api")
{
// Create user record
api.POST("/users", func(c *gin.Context) {
user := User{}
err := c.Bind(&user)
if err != nil {
c.JSON(http.StatusBadRequest,
gin.H{
"status": "failed",
"message": "Invalid request body",
})
return
}
user.ID = bson.NewObjectId()
user.CreatedAt, user.UpdatedAt = time.Now(), time.Now()
session := connect()
defer session.Close()
err = session.DB(database).C(collection).Insert(user)
if err != nil {
c.JSON(http.StatusBadRequest,
gin.H{
"status": "failed",
"message": "Error in the user insertion",
})
return
}
c.JSON(http.StatusOK,
gin.H{
"status": "success",
"user": &user,
})
})
}
router.NoRoute(func(c *gin.Context) {
c.AbortWithStatus(http.StatusNotFound)
})
router.Run(":8000")
}
func main() {
StartService()
}