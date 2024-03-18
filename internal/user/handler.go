package user

import (
	"encoding/json"
	"github.com/gorilla/mux"
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
	//Проставить заголовки w.WriteHeader(200)
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
        actors a
    JOIN
        movieactors ma ON a.actor_id = ma.actor_id
    JOIN
        movies m ON ma.movie_id = m.movie_id;
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
	jsonResponse, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}
	// Записываем JSON-ответ в ResponseWriter
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

// GetMovieList получает список фильмов. По умолчанию - сортировка по рейтингу
func (h *handler) GetMovieList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Устанавливаем соединение с бд
	db := postgres.SetupDB()
	// Закрываем взаимодействие с бд в конце
	defer func() {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}()
	rows, err := db.Query(`SELECT * FROM movies ORDER BY rating DESC;`)
	if err != nil {
		panic(err)
	}
	defer func() {
		rows.Close()
	}()
	a := make([]Movie, 0)
	for rows.Next() {
		var b Movie
		if err := rows.Scan(&b); err != nil {
			log.Fatal(err)
		} else {
			a = append(a, b)
		}
	}
	// Преобразуем срез в формат JSON
	jsonResponse, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}
	// Записываем JSON-ответ в ResponseWriter
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

// CreateActor добавляет информацию об актере.
func (h *handler) CreateActor(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Прочитать JSON данные из тела запроса
	var a Actors
	err := json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		http.Error(w, "Failed to decode JSON data", http.StatusBadRequest)
		return
	}

	// Добавляем нового актёра
	// Устанавливаем соединение с бд
	db := postgres.SetupDB()
	// Закрываем взаимодействие с бд в конце
	defer func() {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}()
	_, err := db.Query(`INSERT INTO actors (name, gender, date_of_birth) VALUES (a.ActorsName,a.ActorsGender,a.ActorsDateOfBirth);`)
	if err != nil {
		http.Error(w, "Database entry error", http.StatusBadRequest)
		return
	}
	// Логирование
	log.Printf("Received actor: %+v", a)

	// Возвращаем ответ
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Actor created successfully"))
}

// CreateMovie добавляет информацию о фильме.
func (h *handler) CreateMovie(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Прочитать JSON данные из тела запроса
	var b Movie
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		http.Error(w, "Failed to decode JSON data", http.StatusBadRequest)
		return
	}
	// Добавляем нового актёра
	// Устанавливаем соединение с бд
	db := postgres.SetupDB()
	// Закрываем взаимодействие с бд в конце
	defer func() {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}()
	_, err := db.Query(`INSERT INTO movies (title, description, release_date,rating) VALUES (b.MovieName,b.MovieDescription,b.MovieReleaseDate,b.MovieRating);`)
	if err != nil {
		http.Error(w, "Database entry error", http.StatusBadRequest)
		return
	}
	// Логирование
	log.Printf("Received actor: %+v", b)
	// Возвращаем ответ
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Actor created successfully"))
}

// UpdateActor изменяет информацию об актере
func (h *handler) UpdateActor(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	var a Actors
	a.ActorsID = r.FormValue("actorid")
	a.ActorsName = r.FormValue("actorname")
	a.ActorsGender = r.FormValue("actorgender")
	a.ActorsDateOfBirth = r.FormValue("actordateofbirth")

	db := postgres.SetupDB()
	// Закрываем взаимодействие с бд в конце
	defer func() {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}()
	_, err = db.Exec("update actors set name=?, gender=?, date_of_birth = ? where actor_id = ?", a.ActorsName, a.ActorsGender, a.ActorsDateOfBirth, a.ActorsID)
	if err != nil {
		log.Println(err)
	}
	http.Redirect(w, r, "/", 301)
}

// UpdateMovie изменяет информацию о фильме
func (h *handler) UpdateMovie(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	var b Movie
	b.MovieID = r.FormValue("movieid")
	b.MovieName = r.FormValue("moviename")
	b.MovieDescription = r.FormValue("moviedescription")
	b.MovieReleaseDate = r.FormValue("release_date")
	b.MovieRating = r.FormValue("movierating")
	db := postgres.SetupDB()
	// Закрываем взаимодействие с бд в конце
	defer func() {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}()
	_, err = db.Exec("update movies set title=?, description=?, release_date = ?,rating=? where movie_id = ?", b.MovieName, b.MovieDescription, b.MovieReleaseDate, b.MovieRating, b.MovieID)
	if err != nil {
		log.Println(err)
	}
	http.Redirect(w, r, "/", 301)
}

// DeleteActor удаляет информацию об актёре
func (h *handler) DeleteActor(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	vars := mux.Vars(r)
	id := vars["id"]
	// Устанавливаем соединение с бд
	db := postgres.SetupDB()
	// Закрываем взаимодействие с бд в конце
	defer func() {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}()
	_, err := db.Exec("delete from actors where actor_id = ?", id)
	if err != nil {
		log.Println(err)
	}
	http.Redirect(w, r, "/", 301)
}

// DeleteMovie удаляет информацию о фильме
func (h *handler) DeleteMovie(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	vars := mux.Vars(r)
	id := vars["id"]
	// Устанавливаем соединение с бд
	db := postgres.SetupDB()
	// Закрываем взаимодействие с бд в конце
	defer func() {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}()
	_, err := db.Exec("delete from movies where movie_id = ?", id)
	if err != nil {
		log.Println(err)
	}
	http.Redirect(w, r, "/", 301)
}
