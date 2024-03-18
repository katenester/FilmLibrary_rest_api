package user

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

// RoleBasedAuthMiddleware проверяет роль пользователя и разрешает или запрещает доступ к ресурсам.
func RoleBasedAuthMiddleware(allowedRoles []string, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		log.Println("Проверка валидации")
		// Роли пользователя передаются в заголовке "X-User-Role".
		userRole := r.Header.Get("X-User-Role")

		// Проверяем, есть ли у пользователя необходимая роль.
		roleAllowed := false
		for _, role := range allowedRoles {
			if role == userRole {
				roleAllowed = true
				break
			}
		}
		// Если у пользователя есть роль, пропускаем запрос к следующему обработчику.
		if roleAllowed {
			log.Println("Проверка пройдена.Переход к обработчику")
			next(w, r, params)
		} else {
			// Если у пользователя нет необходимой роли, возвращаем ошибку доступа.
			http.Error(w, "Доступ запрещен", http.StatusForbidden)
			log.Println("Попытка перехода на запрещенный доступ")
		}
	}
}
