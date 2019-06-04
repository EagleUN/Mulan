package app

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/machinebox/graphql"
)

type App struct {
	Router   *mux.Router
	Database *sql.DB
}

func (app *App) SetupRouter() {
	app.Router.
		Methods("GET").
		Path("/").
		HandlerFunc(app.root)

	basedir := "/shares"
	app.Router.
		Methods("POST").
		Path(basedir + "/create/{userId}/{postId}").
		HandlerFunc(app.createShare)

	app.Router.
		Methods("DELETE").
		Path(basedir + "/delete/{userId}/{postId}").
		HandlerFunc(app.deleteShare)

	app.Router.
		Methods("GET").
		Path(basedir + "/get/{userId}").
		HandlerFunc(app.getShares)
}

func (app *App) createShare(w http.ResponseWriter, r *http.Request) {
	log.Println("trying to create a share relationship")
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

		client := graphql.NewClient("http://35.232.95.82:5000/graphql")

		req := graphql.NewRequest(`
		mutation ($follower: String!, $postId: String!) {
			createShareNotification(notification: {
				follower: $follower
				post_id: $postId
			}) {
				follower
				post_id
			}
		}
	`)

		req.Var("follower", userId)
		req.Var("postId", postId)

		ctx := context.Background()

		var res response
		if err := client.Run(ctx, req, &res); err != nil {
			log.Println(err)
		} else {
			log.Println("Notifications triggered correctly")
		}

		w.WriteHeader(http.StatusCreated)
		share := &Share{}
		err := app.Database.QueryRow("SELECT * FROM `shares` WHERE userId = ? and  postId = ?", userId, postId).Scan(&share.UserId, &share.PostId, &share.SharedAt)
		if err != nil {
			log.Println("Database SELECT failed")
		}
		if err := json.NewEncoder(w).Encode(share); err != nil {
			panic(err)
		}
	}

}

func (app *App) root(w http.ResponseWriter, r *http.Request) {

	log.Println("Mulan is ready to save china")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Mulan is ready to save china")
}

func (app *App) getShares(w http.ResponseWriter, r *http.Request) {
	log.Println("trying to get postIds that the user shared")
	vars := mux.Vars(r)
	userId, ok := vars["userId"]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print("No userId in the path")
	}
	rows, err := app.Database.Query("SELECT * FROM `shares` where userId = ?", userId)
	defer rows.Close()
	if err != nil {
		log.Println(err)

	}
	querys := make([]*Share, 0)
	for rows.Next() {
		share := &Share{}
		errRow := rows.Scan(&share.UserId, &share.PostId, &share.SharedAt)
		if errRow != nil {
			log.Println(err)
			continue
		}
		querys = append(querys, share)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(querys); err != nil {
		panic(err)
	}

}

func (app *App) deleteShare(w http.ResponseWriter, r *http.Request) {
	log.Println("trying to delete a share relationship")
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

	share := &Share{}
	errSel := app.Database.QueryRow("SELECT * FROM `shares` WHERE userId = ? and  postId = ?", userId, postId).Scan(&share.UserId, &share.PostId, &share.SharedAt)
	if errSel != nil {
		log.Println("Database SELECT failed")
	}
	log.Println(share.PostId)
	log.Println(share.UserId)

	_, err := app.Database.Exec("DELETE FROM `shares` WHERE (`userId` = ?) and (`postId` = ?)", userId, postId)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err)
		flag = false
	}
	if flag {
		log.Println("You deleted a share relationship!")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(share); err != nil {
			panic(err)
		}
	}

}
