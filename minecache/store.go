package minecache

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/myadamtest/logkit"
	"reflect"
	"time"
)

type Store struct {
	db *sqlx.DB
	rc *redis.Client
}

//"stu:1234qwer@tcp(10.0.0.241:3307)/test?charset=utf8"
func NewStore(mysqlAddr, redisAddr string) *Store {
	db, err := sqlx.Open("mysql", mysqlAddr)
	if err != nil {
		logkit.Errorf("open mysql failed,err :%s", err)
		return nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,      // "127.0.0.1:6379"
		Password: "dengdeng123w", // no password set
		DB:       0,              // use default DB
	})

	store := &Store{}
	store.db = db
	store.rc = client
	return store
}

func (s *Store) Add(value interface{}) {
	t := reflect.TypeOf(value)
	if t.Kind() != reflect.Ptr {
		logkit.Errorf("need ptr type value")
		return
	}

	sp := getAddScope(value)
	if sp == nil {
		return
	}

	result, err := s.db.Exec(sp.sql, sp.sqlVars...)
	if err != nil {
		logkit.Errorf("exec sql,err :%s", err)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		logkit.Errorf("get last insert id,err :%s", err)
		return
	}

	if id <= 0 {
		return
	}

	if sp.primaryKey != "" {
		pp := reflect.ValueOf(value) // 取得struct变量的指针
		pp.Elem().FieldByName(sp.primaryKey).SetInt(id)
	}

	b, _ := json.Marshal(value)
	err = s.rc.Set(fmt.Sprintf(sp.redisKey, id), string(b), time.Minute).Err()
	if err != nil {
		logkit.Errorf("redis err :%s", err)
		return
	}
}

func (s *Store) Finds(value interface{}, expression string, where ...interface{}) {
	t := reflect.TypeOf(value)
	h := reflect.New(t.Elem().Elem())

	sp := getFindScope(h, expression, where...)

	fmt.Println(h)
}
