package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
)

func PingDB() http.HandlerFunc {
	fn := func(res http.ResponseWriter, req *http.Request) {

		if req.Method != http.MethodGet {
			fmt.Print("Only GET requests are allowed!\n")
			http.Error(res, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
			return
		}

		ctx := req.Context()
		dbAny := ctx.Value("db")
		if db, ok := dbAny.(*sql.DB); ok {
			err := db.PingContext(ctx)
			if err != nil {
				fmt.Printf("Ошибка соединения с базой: %v \n", err)
				res.WriteHeader(http.StatusInternalServerError)
			}
			fmt.Println("Соединение с базой установлено")
			res.WriteHeader(http.StatusOK)
		} else {
			fmt.Println("Ошибка чтения контекста")
			res.WriteHeader(http.StatusInternalServerError)
		}

	}
	return http.HandlerFunc(fn)
}
