package main

import (

	"net/http"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/julienschmidt/httprouter"
	"github.com/op/go-logging"

	"baas/app-wallet/consonlesrvc/auth"
	"baas/app-wallet/consonlesrvc/common"
	"baas/app-wallet/consonlesrvc/database"
	"baas/app-wallet/consonlesrvc/wallet/account"
	"baas/app-wallet/consonlesrvc/wallet/transaction"
	"baas/app-wallet/consonlesrvc/blockchain"
	"github.com/robfig/cron"
	"baas/app-wallet/consonlesrvc/wallet/cronjob"
)

var consLogger *logging.Logger = common.NewLogger("console") //logging.MustGetLogger("Console")


func main() {
	var db *sql.DB = database.GetDB()
	defer db.Close()
	var c = cron.New()
	c.AddJob("*/1 * * * * ?", &cronjob.JobCreateAccount{})
	c.AddJob("*/1 * * * * ?", &cronjob.JobAccountTransfer{})
	c.Start()
	defer c.Stop()

	router := httprouter.New()


	router.Handle("POST", "/auth/login", authsrvc.LoginPost)
	router.Handle("POST", "/auth/signup", authsrvc.SignupPost)
	router.Handle("POST", "/auth/refresh", authsrvc.RefreshPost)
	router.Handle("POST", "/auth/logout", authsrvc.LogoutPost)

	router.Handle("POST", "/wallet/account/create", account.AccountCreatePost)
	router.Handle("POST", "/wallet/account/list", account.AccountListPost)
	router.Handle("POST", "/wallet/account/transfer", account.TransferPost)
	router.Handle("POST", "/wallet/transaction/list", transaction.TransactionListPost)

	router.Handle("POST", "/blockchain/transaction", blockchain.TransactionDetailPost)

	consLogger.Info("start to listen and serve for localhost:8765")
	consLogger.Fatal(http.ListenAndServe(":8765", router))
}
