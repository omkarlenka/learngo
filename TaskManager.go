package main

import (
	"strconv"
	"fmt"
	"flag"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"os"
	"strings"
	// "errors"
)

type task struct{
	id int
	name string
}

var db *sql.DB
var err error
var idToPriority map[int]int
var priorityToId map[int]int


var createDbQuery = "CREATE DATABASE IF NOT EXISTS taskdb"
var useDbQuery = "USE taskdb"
var createTableQuery = `CREATE TABLE IF NOT EXISTS tasks(
	id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
	name varchar(255) NOT NULL
  );`
var insertQuery = "INSERT INTO tasks (name) VALUES "
var insertAllQuery = "INSERT INTO tasks VALUES "
var deleteQuery = "DELETE FROM tasks WHERE id="
var selectQuery = "SELECT id,name from tasks"
var truncateQuery = "TRUNCATE TABLE tasks"
var existsQuery = "SELECT EXISTS (SELECT 1 FROM tasks)"
var deleteAllQuery = "DELETE from tasks"
var selectIDQuery = "SELECT id from tasks"
var selectNameQuery = "SELECT name from tasks WHERE id="

func init(){
	db,err = sql.Open("mysql","root:omkar@123@tcp(127.0.0.1:3306)/")
	logFatalError(err)

	_,err = db.Exec(createDbQuery)
	logFatalError(err)

	_,err = db.Exec(useDbQuery)
	logFatalError(err)

	_,err = db.Exec(createTableQuery)
	logFatalError(err)

	rows, err := db.Query(selectIDQuery)
	logFatalError(err)

	ids := make([]*int,0)
	for rows.Next(){
		var id int
		err := rows.Scan(&id)
		logFatalError(err)

		ids = append(ids, &id)
	}

	idToPriority = make(map[int]int)
	priorityToId = make(map[int]int)
	p:=1
	for _,id := range ids{
		idToPriority[*id] = p
		priorityToId[p] = *id
		p++
	}
}

func swapPriority(id1 int, id2 int){
	var name1 string
	var name2 string

	rows1, err := db.Query(selectNameQuery + strconv.Itoa(id1))
	logFatalError(err)
	for rows1.Next(){
		err := rows1.Scan(&name1)
		logFatalError(err)
	}

	rows2, err := db.Query(selectNameQuery + strconv.Itoa(id2))
	logFatalError(err)
	for rows2.Next(){
		err := rows2.Scan(&name2)
		logFatalError(err)
	}

	_,err = db.Exec(deleteQuery + strconv.Itoa(id1))
	logFatalError(err)

	_,err = db.Exec(deleteQuery + strconv.Itoa(id2))
	logFatalError(err)

	_,err = db.Exec(insertAllQuery + "(" + strconv.Itoa(id1) + ",\"" + name2 + "\")")
	logFatalError(err)

	_,err = db.Exec(insertAllQuery + "(" + strconv.Itoa(id2) + ",\"" + name1 + "\")")
	logFatalError(err)
}

func logFatalError(err error){
	if err != nil{
		log.Fatal(err)
		os.Exit(1)
	}
}

func main(){
	// Parse input Flags
	addPtr := flag.String("add","","Add a task");
	deletePtr := flag.Int("done",0,"Enter Task Numer to be deleted");
	deleteAllPtr := flag.Bool("doneall",false,"Mark all tasks as done")
	swapPtr := flag.String("swap", "", "Enter the comma separated task numbers to be swapped")

	flag.Parse()

	// Perform db operations based on input flags
	if(len(*addPtr) > 0){
		// Add task to db
		_,err = db.Exec(insertQuery + "(\"" + *addPtr + "\")")
		logFatalError(err)
	} else if(*deletePtr != 0){
		// Delete task from db
		_,err = db.Exec(deleteQuery + strconv.Itoa(priorityToId[*deletePtr]))
		logFatalError(err)

		//Check if table is empty, reset id
		rows, err := db.Query(existsQuery)
		logFatalError(err)

		var exists bool
		for rows.Next(){		
			err = rows.Scan(&exists)
		}

		if !exists{
			_,err = db.Exec(truncateQuery)
			logFatalError(err)
		}
	}else if(*deleteAllPtr == true){
		_,err = db.Exec(deleteAllQuery)
		logFatalError(err)

		_,err = db.Exec(truncateQuery)
		logFatalError(err)
	} else if (len(*swapPtr) != 0){
		
		s := strings.Split(*swapPtr, ",")
		if(len(s)>2){
			os.Exit(1)
		}

		id1,err := strconv.Atoi(s[0])
		if(err != nil){
			os.Exit(1)
		}

		id2,err := strconv.Atoi(s[1])
		if(err != nil){
			os.Exit(1)
		}

		swapPriority(priorityToId[id1], priorityToId[id2])
	}else{
		rows, err := db.Query(selectQuery)
		logFatalError(err)

		tasks := make([]*task,0)
		for rows.Next(){
			task := new(task)
			err := rows.Scan(&task.id,&task.name)
			logFatalError(err)
			tasks = append(tasks, task)
		}
		
		for _,task := range tasks{
			fmt.Println(strconv.Itoa(idToPriority[task.id]) + "-" + task.name)
		}

	}
}