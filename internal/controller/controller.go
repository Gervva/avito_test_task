package controller

import (
	"net/http"

	auth "github.com/Gervva/avito_test_task/pkg/authorisation"
)

func AdminMW(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("token")

		exist, role := auth.GetUser(token)
		if !exist {
			http.Error(w, ErrMsgUnauthorized, http.StatusUnauthorized)
			return
		}
		if role != auth.UserRoleAdmin {
			http.Error(w, ErrMsgForbidden, http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func UserMW(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("token")

		exist, role := auth.GetUser(token)
		if !exist {
			http.Error(w, ErrMsgUnauthorized, http.StatusUnauthorized)
			return
		}

		if role != auth.UserRoleRegular && role != auth.UserRoleAdmin {
			http.Error(w, ErrMsgForbidden, http.StatusForbidden)
			return
		}

		if role == auth.UserRoleAdmin {
			query := r.URL.Query()
			query.Add("is_admin", "true")
			r.URL.RawQuery = query.Encode()
		}

		next.ServeHTTP(w, r)
	})
}
