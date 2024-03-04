package main

import (
	"go-media-stream/internal/app"

	"github.com/sirupsen/logrus"
)

func main() {
	err := app.Run()
	if err != nil {
		logrus.Error(err)
	}

	// r := chi.NewRouter()

	// fs := http.FileServer(http.Dir("./static"))
	// r.Handle("/static/*", http.StripPrefix("/static/", fs))

	// store := local.NewLocalStore()
	// conf := config.NewConfig()
	// mysql, err := mysql.NewMysql(conf.DBConnect)
	// if err != nil {
	// 	panic(err)
	// }

	// result := make(chan ffmpeg.WorkerResult, conf.FilmDecWrCount)
	// chExit := make(chan struct{}, 1)
	// ff := ffmpeg.NewFFMPEG(conf)
	// ff.CreateWorkers(conf.FilmDecWrCount, result)
	// go func() {
	// 	for {
	// 		select {
	// 		case workerRes := <-result:
	// 			logrus.Info("Worker result", workerRes)
	// 			err := mysql.AddFilm(context.Background(), &workerRes)
	// 			if err != nil {
	// 				logrus.Error(err)
	// 			}
	// 		case <-chExit:
	// 			logrus.Info("chExit")
	// 			return
	// 		}
	// 	}
	// }()

	// index.NewHandler(conf, store, mysql).Register(r)
	// video.NewHandler(conf, store, mysql).Register(r)
	// stream.NewHandler(conf, store, mysql).Register(r)
	// upload.NewHandler(conf, store, mysql, ff).Register(r)

	// server := &http.Server{Addr: conf.ServerAddress, Handler: r}
	// go func() {
	// 	logrus.Infof("Сервер запушен %s", conf.ServerAddress)
	// 	if err := server.ListenAndServe(); err != nil {
	// 		if !errors.Is(err, http.ErrServerClosed) {
	// 			logrus.Error(err)
	// 		}
	// 	}
	// }()

	// terminateSignals := make(chan os.Signal, 1)
	// signal.Notify(terminateSignals, syscall.SIGINT)
	// <-terminateSignals

	// if err := server.Shutdown(context.Background()); err != nil {
	// 	logrus.Error(err)
	// }

	// ff.Close()
	// ff.Wait()
	// chExit <- struct{}{}
	// close(chExit)

	// logrus.Info("Сервер остановлен нормально")
}
