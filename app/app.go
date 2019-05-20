package app

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type App struct {
	Router   *mux.Router
	Database *sql.DB
}

func (app *App) SetupRouter() {
	basedir := "/shares"
	app.Router.
		Methods("POST").
		Path(basedir + "/users/{userId}/sharing/{postId}").
		HandlerFunc(app.createShare)

	app.Router.
		Methods("DELETE").
		Path(basedir + "/users/{userId}/sharing/{postId}").
		HandlerFunc(app.deleteShare)

	app.Router.
		Methods("GET").
		Path(basedir + "/users/{userId}/sharing").
		HandlerFunc(app.getShares)
}

func (app *App) createShare(w http.ResponseWriter, r *http.Request) {
	flag := true
	vars := mux.Vars(r)
	userId, ok := vars["userId"]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print("No userId in the path")
		flag = false
	}
	postId, ok := vars["postId"]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal("No postId in the path")
		flag = false
	}

	_, err := app.Database.Exec("INSERT INTO `shares` (`userId`, `postId`) VALUES (?, ?)", userId, postId)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err)
		flag = false
	}
	if flag {
		log.Println("You created a share relationship!")
		w.WriteHeader(http.StatusCreated)
	}

}

func (app *App) getShares(w http.ResponseWriter, r *http.Request) {
	log.Println("trying to get postIds that the user shared")
	vars := mux.Vars(r)
	userId, ok := vars["userId"]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print("No userId in the path")
	}
	log.Println(userId)

	rows, err := app.Database.Query("SELECT postId, sharedAt FROM `shares` where userId = ?", userId)
	defer rows.Close()
	if err != nil {
		log.Println(err)

	}
	querys := make([]*Post, 0)
	for rows.Next() {
		post := &Post{}
		errRow := rows.Scan(&post.PostId, &post.SharedAt)
		if errRow != nil {
			log.Println(err)
			continue
		}
		log.Println(post.PostId)
		querys = append(querys, post)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(querys); err != nil {
		panic(err)
	}

}

func (app *App) deleteShare(w http.ResponseWriter, r *http.Request) {
	flag := true
	vars := mux.Vars(r)
	userId, ok := vars["userId"]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print("No userId in the path")
		flag = false
	}
	postId, ok := vars["postId"]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal("No postId in the path")
		flag = false
	}

	_, err := app.Database.Exec("DELETE FROM `shares` WHERE (`userId` = ?) and (`postId` = ?)", userId, postId)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err)
		flag = false
	}
	if flag {
		log.Println("You deleted a share relationship!")
		w.WriteHeader(http.StatusCreated)
	}

}
