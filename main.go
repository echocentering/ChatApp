package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"os/signal"
)

const (
	host = "2601:647:5e00:4744:dc20:4a91:19c7:b0f1"
	port = 5432
	user = "postgres"
	password = "mysecretpassword"
	dbname = "testdb"
)

func main() {
	fmt.Println("HELLO CHAT APP")
	fmt.Println("Good morning, my lovely smiley friend")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	fmt.Println("Postgres INFO", psqlInfo)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// by calling db.Ping(), we force our code to actually open up a connection
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully Connected!")

	stmt := `CREATE TABLE IF NOT EXISTS moodlog (
		name VARCHAR(50) PRIMARY KEY,
		mood VARCHAR(50) NOT NULL,
		training_plan VARCHAR(200) NOT NULL
	);`
	result, err := db.Exec(stmt)
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully added a table", result)

	srv := &http.Server{
		Addr: ":8000",
		//Handler: http.FileServer(http.Dir("./static")),
	}
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./static"))))
	http.Handle("/meditation/", http.StripPrefix("/meditation/", http.FileServer(http.Dir("./meditation"))))
	go func() {
		if err = srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServer(): %v", err)
		}
	}()
	fmt.Println("main: serving on http server at port 8000")

	// Results in an error, because the entry already exists.
	//insertStmt := `INSERT INTO moodlog(name, mood, training_plan)
	//VALUES('emma', 'nervous', 'trying to work on this project');`
	//result, err = db.Exec(insertStmt)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("Successfully entered one row")

	// UNCOMMENT THIS LINE AND IT WILL WORK.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	if err = srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("Shutdown(): %v", err)
	}
}
