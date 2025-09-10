package shutdown

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"
)

func GracefulShutdown(kafkaContextCancelFunc context.CancelFunc, srv *http.Server, SQLDB *sql.DB) {
	fmt.Println("Сигнал завершения получен, останавливаем сервис...")
	//контекст с таймаутом для api
	apiContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := srv.Shutdown(apiContext)
	if err != nil {
		log.Println("Ошибка при завершении работы API: ", err)
	}
	kafkaContextCancelFunc()
	err = SQLDB.Close()
	if err != nil {
		log.Println("Ошибка при завершении работы БД: ", err)
	}
}
