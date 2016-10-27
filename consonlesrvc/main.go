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

/*
var routes map[string]interface{}= map[string]interface{}{
	"authsrvc.Login": new(authsrvc.Login),
	"authsrvc.Signup": new(authsrvc.Signup),
	"authsrvc.Refresh": new(authsrvc.Refresh),
	"authsrvc.Logout": new(authsrvc.Logout),
}
*/

func main() {
	var db *sql.DB = database.GetDB()
	defer db.Close()
	var c = cron.New()
	c.AddJob("*/5 * * * * ?", &cronjob.JobCreateAccount{})
	c.AddJob("*/5 * * * * ?", &cronjob.JobAccountTransfer{})
	c.Start()
	defer c.Stop()

	router := httprouter.New()

	/*
		prefix := "/auth/"
		for _, instance := range routes {
			v := reflect.ValueOf(instance)
			t := v.Type()
			for i := 0; i < t.NumMethod(); i++{
				methodType := v.Method(i).Type()
				consLogger.Debugf("func (%s) %s%s", t, t.Method(i).Name,
					strings.TrimPrefix(methodType.String(), "func"))
				paths := strings.Split(t.String(),".")
				path := prefix + strings.ToLower(paths[len(paths)-1])
				switch instance.(type) {
				case authsrvc.Login:
					f := instance.(authsrvc.Login)
					router.Handle(strings.ToUpper(t.Method(i).Name), path, f.Post)
				case authsrvc.Signup:
					f := instance.(authsrvc.Signup)
					router.Handle(strings.ToUpper(t.Method(i).Name), path, f.Post)
				case authsrvc.Refresh:
					f := instance.(authsrvc.Refresh)
					router.Handle(strings.ToUpper(t.Method(i).Name), path, f.Post)
				case authsrvc.Logout:
					f := instance.(authsrvc.Logout)
					router.Handle(strings.ToUpper(t.Method(i).Name), path, f.Post)
				}
				consLogger.Infof("register router %s with mehtod %s", strings.ToUpper(t.Method(i).Name), path)
			}
		}
	*/

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
