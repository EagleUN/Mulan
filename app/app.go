package app

import (
	"database/sql"
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

	_, err := app.Database.Exec("INSERT INTO `shares` (`uuid`, `userId`, `postId`, `sharedAt`) VALUES (uuid(),?, ? , CURTIME())", userId, postId)

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
