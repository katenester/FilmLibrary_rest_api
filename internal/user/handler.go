package user

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/julienschmidt/httprouter"
	"github.com/katenester/FilmLibrary_rest_api/internal/handlers"
	"github.com/katenester/FilmLibrary_rest_api/internal/repository/postgres"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

const (
	actor        = "/actor"
	actorURL     = "/actor/:uuid"
	actors       = "/actors"
	actorsSearch = "/search/actors"

	movie        = "/movie"
	movieURL     = "/movie/:uuid"
	movies       = "/movies"
	moviesSearch = "/search/movies"
)

type handler struct {
}

func NewHandler() handlers.Handler {
	return &handler{}
}

// Register - регистрация обработчиков handler
func (h *handler) Register(router *httprouter.Router) {
	// регистрируем пути с middleware, для проверки уровня доступа
	router.GET(actors, RoleBasedAuthMiddleware([]string{"admin", "user"}, h.GetActorList))
	router.GET(movies, RoleBasedAuthMiddleware([]string{"admin", "user"}, h.GetMovieList))
	router.POST(actor, RoleBasedAuthMiddleware([]string{"admin"}, h.CreateActor))
	router.POST(movie, RoleBasedAuthMiddleware([]string{"admin"}, h.CreateMovie))
	router.PUT(actorURL, RoleBasedAuthMiddleware([]string{"admin"}, h.UpdateActor))
	router.PUT(movieURL, RoleBasedAuthMiddleware([]string{"admin"}, h.UpdateMovie))
	router.DELETE(actorURL, RoleBasedAuthMiddleware([]string{"admin"}, h.DeleteActor))
	router.DELETE(movieURL, RoleBasedAuthMiddleware([]string{"admin"}, h.DeleteMovie))
	router.GET(actorsSearch, RoleBasedAuthMiddleware([]string{"admin", "user"}, h.SearchMovies))
	router.GET(moviesSearch, RoleBasedAuthMiddleware([]string{"admin", "user"}, h.SearchActors))
}

// GetActorList получает список актёров.
func (h *handler) GetActorList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Println("Вход в GetActorList")
	//Проставить заголовки w.WriteHeader(200)
	// Устанавливаем соединение с бд
	log.Println("Подключение к бд")
	db := postgres.SetupDB()
	// Закрываем взаимодействие с бд в конце
	defer func() {
		log.Println("Закрытие подключения")
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	log.Println("Отправка запроса на получение списка актёров")
	rows, err := db.Query(`
    SELECT
        name AS actor_name,
        title AS movie_title
    FROM
        actors
    JOIN
        movieactors  ON actors.actor_id = movieactors.actor_id
    JOIN
        movies ON movies.movie_id = movieactors.movie_id;
`)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		rows.Close()
	}()
	a := make([]Actors, 0)
	for rows.Next() {
		var b Actors
		if err := rows.Scan(&b); err != nil {
			a = append(a, b)
		} else {
			log.Fatal(err)
		}
	}
	// Преобразуем срез в формат JSON
	jsonResponse, err := json.Marshal(a)
	if err != nil {
		log.Fatal(err)
	}
	// Записываем JSON-ответ в ResponseWriter
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// GetMovieList получает список фильмов. По умолчанию - сортировка по рейтингу
func (h *handler) GetMovieList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Println("Вход в GetMovieList")
	// Устанавливаем соединение с бд
	log.Println("Подключение к бд")
	db := postgres.SetupDB()
	// Закрываем взаимодействие с бд в конце
	defer func() {
		log.Println("Закрытие подключения")
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	log.Println("Отправка запроса на получение списка фильмов")
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
	w.WriteHeader(http.StatusOK)
	// Записываем JSON-ответ в ResponseWriter
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

// CreateActor добавляет информацию об актере.
func (h *handler) CreateActor(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Println("Вход в CreateActor")
	// Прочитать JSON данные из тела запроса
	var a Actors
	err := json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		http.Error(w, "Ошибка декодирования JSON данных", http.StatusBadRequest)
		log.Println("Ошибка декодирования JSON данных", err)
		return
	}
	// Добавляем нового актёра
	// Устанавливаем соединение с бд
	log.Println("Подключение к бд")
	db := postgres.SetupDB()
	// Закрываем взаимодействие с бд в конце
	defer func() {
		log.Println("Закрытие бд")
		err := db.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	_, err = db.Query("INSERT INTO actors (name, gender, date_of_birth) VALUES ($1, $2, $3);", a.ActorsName, a.ActorsGender, a.ActorsDateOfBirth)
	if err != nil {
		http.Error(w, "Database entry error", http.StatusBadRequest)
		log.Println("Database entry error", err)
		return
	}
	log.Printf("Принятый актер: %+v", a)
	// Возвращаем ответ
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Актер успешно создан"))
}

// CreateMovie добавляет информацию о фильме.
func (h *handler) CreateMovie(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Println("Вход в CreateMovie")
	// Прочитать JSON данные из тела запроса
	var b Movie
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		http.Error(w, "Ошибка декодирования JSON данных", http.StatusBadRequest)
		log.Println("Ошибка декодирования JSON данных", err)
		return
	}
	log.Println("Подключение к бд")
	db := postgres.SetupDB()
	defer func() {
		log.Println("Закрытие бд")
		err := db.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	_, err = db.Query("INSERT INTO movies (title, description, release_date,rating) VALUES ($1, $2,$3, $4);", b.MovieName, b.MovieDescription, b.MovieReleaseDate, b.MovieRating)
	if err != nil {
		http.Error(w, "Database entry error", http.StatusBadRequest)
		log.Println("Database entry error", err)
		return
	}
	// Получение ID только что добавленного фильма
	var movieID int
	row := db.QueryRow("SELECT lastval()")
	err = row.Scan(&movieID)
	if err != nil {
		http.Error(w, "Ошибка получения ID последнего вставленного фильма", http.StatusInternalServerError)
		log.Println("Ошибка получения ID последнего вставленного фильма", err)
		return
	}

	// Вставка данных о связи фильма и актёров в таблицу "movieactors"
	for _, actorID := range b.MovieActor {
		_, err = db.Query("INSERT INTO movieactors (movie_id, actor_id) VALUES ($1, $2);", movieID, actorID.ActorsID)
		if err != nil {
			http.Error(w, "Ошибка вставки данных в базу данных", http.StatusBadRequest)
			log.Println("Ошибка вставки данных в базу данных", err)
			return
		}
	}

	// Логирование
	log.Printf("Принятый фильм: %+v", b)
	// Возвращаем ответ
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Фильм успешно создан"))
}

// UpdateActor изменяет информацию об актере
func (h *handler) UpdateActor(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Println("Вход в UpdateActor")
	// Прочитать JSON данные из тела запроса
	var a Actors
	err := json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		http.Error(w, "Ошибка декодирования JSON данных", http.StatusBadRequest)
		log.Fatal(err)
		return
	}
	log.Println("Подключение к бд")
	db := postgres.SetupDB()
	defer func() {
		log.Println("Закрытие бд")
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	// Изменение в actors
	_, err = db.Exec("UPDATE actors SET name=$1, gender=$2, date_of_birth = $3 WHERE actor_id = $4", a.ActorsName, a.ActorsGender, a.ActorsDateOfBirth, a.ActorsID)
	if err != nil {
		http.Error(w, "Database entry error", http.StatusBadRequest)
		log.Fatal(err)
		return
	}
	log.Printf("Измененный актер: %+v", a)
	// Возвращаем ответ
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Актер успешно изменен"))
}

// UpdateMovie изменяет информацию о фильме
func (h *handler) UpdateMovie(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Println("Вход в UpdateMovie")
	// Прочитать JSON данные из тела запроса
	var a Movie
	err := json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		http.Error(w, "Ошибка декодирования JSON данных", http.StatusBadRequest)
		log.Fatal(err)
		return
	}
	log.Println("Подключение к бд")
	db := postgres.SetupDB()
	defer func() {
		log.Println("Закрытие бд")
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	// Изменение в movie
	_, err = db.Exec("UPDATE movies SET title=$1, description=$2, release_date = $3,rating=$4 WHERE movie_id = $5", a.MovieName, a.MovieDescription, a.MovieReleaseDate, a.MovieRating, a.MovieID)
	if err != nil {
		http.Error(w, "Database entry error", http.StatusBadRequest)
		log.Fatal(err)
		return
	}
	log.Printf("Измененный фильм: %+v", a)
	// Возвращаем ответ
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Фильм успешно изменен"))
}

// DeleteActor удаляет информацию об актёре
func (h *handler) DeleteActor(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Println("Вход в DeleteActor")
	vars := mux.Vars(r)
	id := vars["id"]
	log.Println("Подключение к бд")
	db := postgres.SetupDB()
	defer func() {
		log.Println("Закрытие бд")
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	_, err := db.Exec("delete from actors where actor_id = ?", id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Database entry error", http.StatusBadRequest)
		return
	}
	_, err = db.Exec("delete from movieactors where actor_id = ?", id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Database entry error", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Актёр успешно удалён"))
}

// DeleteMovie удаляет информацию о фильме
func (h *handler) DeleteMovie(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Println("Вход в DeleteMovie")
	vars := mux.Vars(r)
	id := vars["id"]
	log.Println("Подключение к бд")
	db := postgres.SetupDB()
	defer func() {
		log.Println("Закрытие бд")
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	_, err := db.Exec("delete from movies where movie_id = ?", id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Database entry error", http.StatusBadRequest)
		return
	}
	_, err = db.Exec("delete from movieactors where movie_id = ?", id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Database entry error", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Фильм успешно удалён"))
}

// SearchMovies выполняет поиск фильмов по фрагменту названия
func (h *handler) SearchMovies(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Println("Вход в SearchMovies")
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Требуется параметр запроса q", http.StatusBadRequest)
		return
	}
	db := postgres.SetupDB()
	defer func() {
		log.Println("Закрытие бд")
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	// Поиск фильмов по фрагменту названия
	rows, err := db.Query("SELECT * FROM movies WHERE title LIKE '%' || $1 || '%'", query)
	if err != nil {
		http.Error(w, "Не удалось загрузить фильмы", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	defer rows.Close()
	// Формирование списка фильмов
	var movies []Movie
	for rows.Next() {
		var m Movie
		if err := rows.Scan(&m.MovieID, &m.MovieName, &m.MovieDescription, &m.MovieReleaseDate, &m.MovieRating); err != nil {
			log.Println(err)
			continue
		}
		movies = append(movies, m)
	}

	// Отправка JSON-ответа
	jsonResponse, err := json.Marshal(movies)
	if err != nil {
		http.Error(w, "Не удалось выстроить JSON", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// SearchActors выполняет поиск актёров по фрагменту имени
func (h *handler) SearchActors(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Println("Вход в SearchActors")
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Требуется параметр запроса q", http.StatusBadRequest)
		return
	}
	// Устанавливаем соединение с бд
	db := postgres.SetupDB()
	defer func() {
		log.Println("Закрытие бд")
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	// Поиск актёров по фрагменту имени
	rows, err := db.Query("SELECT * FROM actors WHERE name LIKE '%' || $1 || '%'", query)
	if err != nil {
		http.Error(w, "Не удалось загрузить актёров", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	defer rows.Close()
	// Формирование списка актёров
	var actors []Actors
	for rows.Next() {
		var a Actors
		if err := rows.Scan(&a.ActorsID, &a.ActorsName, &a.ActorsGender, &a.ActorsDateOfBirth); err != nil {
			log.Println(err)
			continue
		}
		actors = append(actors, a)
	}
	// Отправка JSON-ответа
	jsonResponse, err := json.Marshal(actors)
	if err != nil {
		http.Error(w, "Не удалось выстроить JSON", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
