package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"

	"github.com/allez-chauffe/marcel/pkg/config"
	"github.com/allez-chauffe/marcel/pkg/db/internal/db"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func unmarshallRow(row scanable, id interface{}, result interface{}) error {
	var data string

	resultVal := reflect.ValueOf(result)

	if resultVal.Kind() != reflect.Ptr ||
		(resultVal.Elem().Kind() != reflect.Ptr &&
			resultVal.Elem().Kind() != reflect.Interface) {
		panic("result should be a pointer of pointer of the targeted entity (**Client by example)")
	}

	resultPtr := resultVal.Elem()

	if err := row.Scan(&id, &data); err != nil {
		if err == sql.ErrNoRows {
			resultPtr.Set(reflect.Zero(resultPtr.Type()))
			return nil
		}
		return err
	}

	var resultMap map[string]interface{}
	if err := json.Unmarshal([]byte(data), &resultMap); err != nil {
		return err
	}

	entity := resultPtr.Interface().(db.Entity)

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: mapstructure.StringToTimeHookFunc(time.RFC3339),
		Result:     entity,
	})
	if err != nil {
		return err
	}
	if err := decoder.Decode(resultMap); err != nil {
		return err
	}
	entity.SetID(id)

	return nil
}

func prepare(item db.Entity) (interface{}, []byte) {
	data := structs.Map(item)
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(fmt.Errorf("Failed to prepare value fo postgres jsonb: %w", err))
	}
	return item.GetID(), jsonData
}

type scanable interface {
	Scan(...interface{}) error
}

func getConnectionString(database string) string {
	pgConf := config.Default().API().DB().Postgres()
	if database == "" {
		database = pgConf.DBName()
	}
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		pgConf.Host(), pgConf.Port(), pgConf.Username(), pgConf.Password(), database,
	)
}

func createDatabase() (*sql.DB, error) {
	dbName := config.Default().API().DB().Postgres().DBName()
	log.Infof("Creating database %s", dbName)

	tempDB, err := sql.Open("postgres", getConnectionString("postgres"))
	if err != nil {
		return nil, err
	}
	defer tempDB.Close()

	if _, err = tempDB.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, dbName)); err != nil {
		return nil, err
	}

	pg, err := sql.Open("postgres", getConnectionString(""))
	if err != nil {
		return nil, err
	}

	return pg, pg.Ping()
}
