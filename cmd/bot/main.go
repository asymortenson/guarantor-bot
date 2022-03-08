package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"guarantorplace.com/internal/config"
	"guarantorplace.com/internal/data"
	"guarantorplace.com/internal/jsonlog"

	tele "gopkg.in/telebot.v3"
)



type application struct {
	config *config.Config
	bot *tele.Bot
	logger *jsonlog.Logger
	models data.Models
}


func main() {


	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)


	pref := tele.Settings{
		Token:  cfg.TelegramToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(pref)

	logger.PrintInfo("bot successfully initialized", nil)

	if err != nil {
		logger.PrintFatal(err, nil)
	}

	db, err := openDB(cfg)

	if err != nil {
		logger.PrintFatal(err, nil)
	}

	defer db.Close()

	log.Println("database connection pool established")

	app := &application{
		config: cfg,
		logger: logger,
		bot: bot,
		models: data.NewModels(db),
	}
	

	err = app.handleUpdates()

	if err != nil {
		logger.PrintError(err, nil)
	}
	app.bot.Start()

}


func openDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.Db.Dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.Db.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Db.MaxIdleConns)

	duration, err := time.ParseDuration(cfg.Db.MaxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
