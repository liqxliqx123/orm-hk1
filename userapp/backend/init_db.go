package main

import (
	"gitee.com/geektime-geekbang/geektime-go/orm"
	"gitee.com/geektime-geekbang/geektime-go/orm/middleware/opentelemetry"
	"gitee.com/geektime-geekbang/geektime-go/orm/middleware/prometheus"
	"gitee.com/geektime-geekbang/geektime-go/orm/middleware/querylog"
	"go.uber.org/zap"
)

func initDB() *orm.DB {
	db, err := orm.Open("mysql", "root:root@tcp(localhost:13306)/userapp",
		orm.DBWithMiddleware(prometheus.MiddlewareBuilder{
			Name:        "userapp",
			Subsystem:   "orm",
			ConstLabels: map[string]string{"db": "userapp"}}.Build(),
			opentelemetry.MiddlewareBuilder{}.Build(),
			querylog.NewBuilder().LogFunc(func(sql string, args []any) {
				// 一般不建议记录参数，因为参数里面可能有一些加密信息，
				// 当然如果你确定自己是在开发环境下使用，那么你就可以记录参数
				zap.L().Info("query", zap.String("SQL", sql))
			}).Build()))
	if err != nil {
		panic(err)
	}
	return db
}
