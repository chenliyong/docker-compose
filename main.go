package main

import (
	"flag"
	"github.com/jmoiron/sqlx"
	"time"
	_ "github.com/go-sql-driver/mysql" // mysql driver
	"log"
)

var (
	URL = flag.String("database_url", "", "Description")
	// 数据库连接
	db *sqlx.DB

	schemas = []string{
		// 客户标签中间表
		`CREATE TABLE IF NOT EXISTS customer_tag (
			tag_id INT(11) NOT NULL,
			customer_id INT(11) NOT NULL,
			automatic TINYINT(1) DEFAULT 0 COMMENT '手动关联or自动关联 0 手动 1自动',
			created DATETIME DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (tag_id, customer_id)
		);`,
	}
)

// Init 初始化
func Init() (*sqlx.DB, error) {
	var err error
	db, err = Connect("mysql", *URL, 10, 10)
	if err != nil {
		return nil, err
	}
	db.MapperFunc(mapper)    // db.NamedExec 使用的mapper方法
	sqlx.NameMapper = mapper // sqlx.Named 使用的mapper方法
	for _, s := range schemas {
		_, err = db.Exec(s)
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}

func mapper(name string) string {
	var s []byte
	for i, r := range []byte(name) {
		if r >= 'A' && r <= 'Z' {
			r += 'a' - 'A'
			if i != 0 {
				s = append(s, '_')
			}
		}
		s = append(s, r)
	}
	return string(s)
}

// Connect 初始化数据库连接池
func Connect(driver, uri string, maxOpen, maxIdel int) (db *sqlx.DB, err error) {
	db, err = sqlx.Connect(driver, uri+"?charset=utf8mb4&parseTime=true")
	if err != nil {
		return
	}
	// 配置连接池
	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxIdel)
	db.SetConnMaxLifetime(2 * time.Hour)
	return
}

func main() {
	flag.Parse()
	log.Println(*URL)
	Init()
	tag_id := 0
	customer_id := 0
	for {
		tag_id++
		customer_id++
		time.Sleep(1*time.Second)
		log.Println(tag_id, customer_id)
		if _, err := db.Exec(`INSERT INTO customer_tag (tag_id, customer_id) VALUES (?,?)`, tag_id, customer_id); err != nil {
			log.Println(err)
		}
	}
}
