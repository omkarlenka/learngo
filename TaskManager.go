package main

import (
	"strconv"
	"fmt"
	"flag"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
)

type task struct{
	id int
	name string
}

func logFatalError(err error){
	if err != nil{
		log.Fatal(err)
	}
}

func main(){
	db,err := sql.Open("mysql","user:pass@tcp(127.0.0.1:3306)/taskdb")
	logFatalError(err)

	//DB Query
	// var dbname = "taskdb"
	var createDbQuery = "CREATE DATABASE IF NOT EXISTS taskdb"
	var createTableQuery = `CREATE TABLE IF NOT EXISTS tasks(
		id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
		name varchar(255) NOT NULL
	  );`
	var insertQuery = "INSERT INTO tasks (name) VALUES "
	var deleteQuery = "DELETE FROM tasks WHERE id="
	var selectQuery = "SELECT id,name from tasks"
	var truncateQuery = "TRUNCATE TABLE tasks"
	var existsQuery = "SELECT EXISTS (SELECT 1 FROM tasks)"

	_,err = db.Exec(createDbQuery)
	logFatalError(err)

	_,err = db.Exec(createTableQuery)
	logFatalError(err)

	// Parse input Flags
	addPtr := flag.String("add","","Add a task");
	deletePtr := flag.Int("done",0,"Enter Task Numer to be deleted");

	flag.Parse()

	// Perform db operations based on input flags
	if(len(*addPtr) > 0){
		// Add task to db
		// fmt.Println(insertQuery + "(\"" + *addPtr + "\")")
		_,err = db.Exec(insertQuery + "(\"" + *addPtr + "\")")
		logFatalError(err)
	} else if(*deletePtr != 0){
		// Delete task from db
		_,err = db.Exec(deleteQuery + strconv.Itoa(*deletePtr))
		logFatalError(err)

		//Check if table is empty, reset id
		rows, err := db.Query(existsQuery)
		logFatalError(err)

		var exists bool
		for rows.Next(){		
			err = rows.Scan(&exists)
			fmt.Println(exists)
		}

		if !exists{
			_,err = db.Exec(truncateQuery)
		}

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
			fmt.Println(strconv.Itoa(task.id) + "-" + task.name)
		}

	}
}