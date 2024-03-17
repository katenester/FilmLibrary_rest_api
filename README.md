Конечные точки

Авторизация:
POST /login: --аутентификация пользователей. Возвращает токен доступа.

Актёры:
POST /actor --Добавление информации об актёре. 
PUT /actor/{actor_id} -- Изменение информации об актёре.
DELETE /actor/{actor_id} -- Удаление информации об актёре.
GET /actors -- Получение списка актёров.

Фильмы:
POST /movie --Добавление информации о фильме.
PUT /movie/{movie_id} --Изменение информации о фильме.
DELETE /movie/{movie_id} --Удаление информации о фильме.
GET /movies --Получение списка фильмов с возможностью сортировки(путь для сортировки по умолчанию )

Поиск:
GET /search/movies?q={query} --Поиск фильма по фрагменту названия.
GET /search/actors?q={query} --Поиск актёра по фрагменту имени.
