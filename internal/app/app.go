package app

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-media-stream/internal/config"
	"go-media-stream/internal/domain"
	"go-media-stream/internal/handlers"
	"go-media-stream/internal/log"
	"go-media-stream/internal/repository"
	"go-media-stream/internal/services"
	"go-media-stream/internal/utils"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"

	_ "github.com/go-sql-driver/mysql"
)

func Run() error {
	conf := config.NewConfig()
	logger := log.New()

	db, err := sql.Open("mysql", conf.DBConnect)
	if err != nil {
		return err
	}
	defer db.Close()
	err = Bootstrap(context.Background(), db)
	if err != nil {
		return err
	}

	userRepository := repository.NewUserRepository(db)
	videoRepository := repository.NewVideoRepository(db)
	audioRepository := repository.NewAudioRepository(db)
	storeRepository := repository.NewStoreRepository()
	queueRepository := repository.NewQueueRepository(db)

	authService := services.NewAuthServices(userRepository)
	uploadService := services.NewUploadServices(videoRepository, storeRepository, queueRepository)

	chiRouter := chi.NewRouter()
	fs := http.FileServer(http.Dir("./static"))
	chiRouter.Handle("/static/*", http.StripPrefix("/static/", fs))

	handlers.NewAuthHandler(logger, authService).Register(chiRouter)
	handlers.NewVideoHandler(videoRepository, audioRepository).Register(chiRouter)
	handlers.NewHomeHandler(videoRepository, audioRepository).Register(chiRouter)
	handlers.NewUploadHandler(uploadService).Register(chiRouter)
	handlers.NewStreamHandler(logger, videoRepository, audioRepository).Register(chiRouter)

	server := &http.Server{Addr: conf.ServerAddress, Handler: chiRouter}
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logrus.Error(err)
		}
	}()

	ffmpeg := utils.NewFFMPEG(logger, db, 1)
	ffmpeg.CreateWorkers(conf.DecoderWorkerCount)

	chQueueDone := make(chan struct{})
	defer close(chQueueDone)
	timeQueueTracker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-timeQueueTracker.C:
				// num := runtime.NumGoroutine()
				// fmt.Printf("Количество запущенных горутин: %d\n", num)
				rows, err := db.QueryContext(context.Background(), `
					SELECT id, folder, type FROM queue WHERE is_done=0 AND is_work=0
				`)
				if err != nil {
					logrus.Error(err)
					return
				}
				for rows.Next() {
					queue := domain.Queue{}
					if err := rows.Scan(&queue.ID, &queue.Folder, &queue.Type); err != nil {
						logrus.Error(err)
						return
					}
					_, err = db.Exec("UPDATE queue SET is_work = 1 WHERE id = ?;", queue.ID)
					if err != nil {
						logrus.Error(err)
						break
					}
					ffmpeg.Add(utils.JobFFMPEG{
						QueueId: queue.ID,
					})
				}
			case <-chQueueDone:
				logrus.Info("Queue остановлен")
				return
			}
		}
	}()

	logrus.Infof("Сервер запушен %s", conf.ServerAddress)

	terminateSignals := make(chan os.Signal, 1)
	signal.Notify(terminateSignals, syscall.SIGINT)
	<-terminateSignals
	logrus.Info("Ждем воркер ffmpeg")
	chQueueDone <- struct{}{}
	timeQueueTracker.Stop()
	ffmpeg.Close()
	ffmpeg.Wait()
	logrus.Info("Сервер остановлен нормально")
	return nil
}
