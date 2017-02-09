package main

import (
	"github.com/gin-gonic/gin"
	"crypto/md5"
	"strconv"
	"encoding/hex"
	"strings"
	"container/list"
	"errors"
)

var MX, MY int = 8, 8
var server *gin.Engine

//Structure for move
type Dir struct {
	y, x int
}

//Object for check
type Field struct{
	x, y int
	horseDir []Dir
}

//convert string position to coordinate
func (field *Field) posToCoordinate(pos string) error {
	var err error
	//convert symbolic notation coordinate to numeric
	field.x = int(strings.ToUpper(pos)[0] - byte('A')+1)
	field.y, err = strconv.Atoi(string(pos[1]))

	//validate input coordinates
	if (field.x > MX && field.y > MY && field.x <= 0 && field.y <= 0 ){
		err = errors.New("Out of range")
	}
	return err
}

//convert move to new string position
func (field *Field) moveToPos(dir Dir) string {
	var  nX, nY int
	nX = field.x + dir.x
	nY = field.y + dir.y

	return string(byte(nX)+64) + strconv.Itoa(nY)

}

//Init object Field, fill all horse's directions
func (field *Field) init()  {
	field.horseDir = []Dir{ {-2, -1}, {-1, -2},
		{-2, 1},{-1, 2},
		{2, -1}, {1, -2},
		{2, 1},{1, 2} }
}

//Validate current move
func (field *Field) moveValid(dir Dir) bool  {
	var dx, dy int
	dx = field.x + dir.x
	dy = field.y + dir.y
	return  (dx <= MX && dy <= MY && dx > 0 && dy > 0 )
}

//Determine all valid moves
func (field *Field) getValidMoves() *list.List  {
	var lMoves *list.List = list.New()
	for _,d := range field.horseDir {
		if field.moveValid(d) {
			var strMove string = field.moveToPos(d)
			lMoves.PushBack(strMove)
		}
	}
	return lMoves
}

//Callback function for GET request /horse
func horse(c *gin.Context) {
	pos := c.Param("xy")
	//check position contain two symbol
	if (len(pos) > 2){
		c.String(400,"Position isn't valid!")
		return
	}
	var field Field = Field{}
	field.init()
	err := field.posToCoordinate(pos)
	if err != nil {
		c.String(400,"Position isn't valid!")
		return
	}
	lMoves := field.getValidMoves()
	var answer string = ""
	for v := lMoves.Front() ; v != nil ; v =v.Next() {

		answer += string(v.Value.(string)) + " "
	}
	c.String(200, answer)
}

//Callback request for POST request for 1-st task
func posting(c *gin.Context)  {

	idstr := c.PostForm("id")
	id64, err := strconv.ParseInt(c.PostForm("id"),10,32)
	var id int = int(id64)
	if err != nil {
		c.String(400,"Not valid JSON")
	}
	text := c.PostForm("text")
	//validate positive number in id and length of text
	if (len(text) > 100 && id < 0) {
		c.String(400,"Not valid JSON")
		return
	} else {
		var X int = id % 2

		param := idstr + text + strconv.Itoa(X)
		var hash [16] byte
		hash = md5.Sum([]byte(param))
		c.String(200, hex.EncodeToString(hash[:]))
	}
}

//Gin engine
func getEngine() *gin.Engine {
	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	router.GET("/horse/:xy", horse)
	router.POST("/md5", posting)
	return router
}

func main() {

    // Listen and server on 0.0.0.0:8000
	server = getEngine()
	server.Run(":8000")

}
