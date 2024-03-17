package user

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/katenester/FilmLibrary_rest_api/internal/handlers"
	"github.com/katenester/FilmLibrary_rest_api/internal/repository/postgres"
	"log"
	"net/http"
)

const (
	actor    = "/actor"
	actorURL = "/actor/:uuid"
	actors   = "/actors"

	movie    = "/movie"
	movieURL = "/movie/:uuid"
	movies   = "/movies"
)

type handler struct {
}

func NewHandler() handlers.Handler {
	return &handler{}
}

// Register - регистрация обработчиков handler
func (h *handler) Register(router *httprouter.Router) {
	// регистрируем пути
	router.GET(actors, h.GetActorList)
	router.GET(movies, h.GetMovieList)
	router.POST(actor, h.CreateActor)
	router.POST(movie, h.CreateMovie)
	router.PUT(actorURL, h.UpdateActor)
	router.PUT(movieURL, h.UpdateMovie)
	router.DELETE(actorURL, h.DeleteActor)
	router.DELETE(movieURL, h.DeleteMovie)
}

// GetActorList получает список актёров.
func (h *handler) GetActorList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Проставляем заголовки
	w.WriteHeader(200)
	// Устанавливаем соединение с бд
	db := postgres.SetupDB()
	// Закрываем взаимодействие с бд в конце
	defer func() {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}()
	rows, err := db.Query(`
    SELECT
        a.actor_id,
        a.name AS actor_name,
        m.movie_id,
        m.title AS movie_title
    FROM
        Actors a
    JOIN
        MovieActors ma ON a.actor_id = ma.actor_id
    JOIN
        Movies m ON ma.movie_id = m.movie_id;
`)
	if err != nil {
		panic(err)
	}
	defer func() {
		rows.Close()
	}()
	a := make([]Actors, 0)
	for rows.Next() {
		var b Actors
		if err := rows.Scan(&b); err != nil {
			log.Fatal(err)
		} else {
			a = append(a, b)
		}
	}
	// Преобразуем срез в формат JSON
	jsonResponse, err := json.Marshal(actors)
	if err != nil {
		panic(err)
	}
	// Записываем JSON-ответ в ResponseWriter
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

// GetMovieList получает список фильмов.
func (h *handler) GetMovieList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Проставляем заголовки
	w.WriteHeader(200)
	w.Write([]byte("This is list of users"))
}

// CreateActor добавляет информацию об актере.
func (h *handler) CreateActor(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(201)
	w.Write([]byte("This is create user"))
}

// CreateMovie добавляет информацию о фильме.
func (h *handler) CreateMovie(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(201)
	w.Write([]byte("This is create user"))
}

// UpdateActor изменяет информацию об актере
func (h *handler) UpdateActor(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(204)
	w.Write([]byte("This is fully update user "))
}

// UpdateMovie изменяет информацию о фильме
func (h *handler) UpdateMovie(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(204)
	w.Write([]byte("This is fully update user "))
}

// DeleteActor удаляет информацию об актёре
func (h *handler) DeleteActor(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(204)
	w.Write([]byte("This is delete user by id"))
}

// DeleteMovie удаляет информацию о фильме
func (h *handler) DeleteMovie(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(204)
	w.Write([]byte("This is delete user by id"))
}
