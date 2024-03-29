package model

import (
    "fmt"
    "os"
    "time"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    "gorm.io/gorm/schema"

    "github.com/hd2yao/ginblog/utils"
)

var db *gorm.DB
var err error

func InitDB() {
    dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        utils.DBUser,
        utils.DBPassword,
        utils.DBHost,
        utils.DBPort,
        utils.DBName,
    )
    db, err = gorm.Open(mysql.Open(dns), &gorm.Config{
        // gorm 日志模式: silent
        Logger: logger.Default.LogMode(logger.Silent),
        // 外键约束
        DisableForeignKeyConstraintWhenMigrating: true,
        // 禁用默认事务（提高运行速度）
        SkipDefaultTransaction: true,
        NamingStrategy: schema.NamingStrategy{
            // 使用单数表名，启用该选项，此时， `User` 的表名应该是 `user`
            SingularTable: true,
        },
    })
    if err != nil {
        fmt.Println("连接数据库失败，请检查参数：", err)
        os.Exit(1)
    }

    // 迁移数据表，在没有数据表结构变更时候，建议注释不执行
    _ = db.AutoMigrate(&User{}, &Article{}, &Category{}, Profile{}, Comment{})

    sqlDB, _ := db.DB()
    // SetMaxIdleConns 设置连接池中的最大闲置连接数
    sqlDB.SetMaxIdleConns(10)

    // SetMaxOpenConns 设置数据库的最大连接数
    sqlDB.SetMaxOpenConns(100)

    // SetConnMaxLifetime 设置连接的最大可复用时间
    sqlDB.SetConnMaxLifetime(10 * time.Second)

}
