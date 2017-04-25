package db
import (
	"database/sql"
	"sync"
	"time"
	"errors"
	"fmt"
	"reflect"
	"unicode"
	_ "github.com/go-sql-driver/mysql"
)

type DBServer struct {
	Host 			string
	Port 			uint16
	User 			string
	Pass 			string
	Charset 		string
	Db 				string
	MaxIdle 		int
	MaxDbConn 		int
	ConnectionPool 	*sql.DB
	LastConnectTime time.Time
	mu 				sync.Mutex
}

func (dbserver *DBServer) conn(maxIdle, maxDbConn int) (err error) {
	if maxIdle <= 0 {
		maxIdle = 8
	}
	if maxDbConn <= 0 {
		maxDbConn = 8
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&timeout=%s",
		dbserver.User,
		dbserver.Pass,
		dbserver.Host,
		dbserver.Port,
		dbserver.Db,
		dbserver.Charset,
		//"30s",
		"300s")
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		db.Close()
		fmt.Println("Open mysql err:", err)
		return
	} else {
		db.SetMaxIdleConns(maxIdle)
		db.SetMaxOpenConns(maxDbConn)
		dbserver.MaxIdle   = maxIdle
		dbserver.MaxDbConn = maxDbConn
		if dbserver.ConnectionPool != nil {
			dbserver.ConnectionPool.Close()
		}
		dbserver.ConnectionPool = db
		dbserver.LastConnectTime = time.Now()
	}
	return
}

func (dbserver *DBServer) Query(result interface{}, query string, args ...interface{}) (err error) {
	for i := 1; i <= 3; i++ {
		err = dbserver.QueryHelper(result, query, args...)
		if err == nil {
			break
		} else {
			if err == errors.New("driver: bad connection") || err == errors.New("sql: database is closed") {
				e := dbserver.conn(dbserver.MaxIdle, dbserver.MaxDbConn)
				if e != nil {
					err = errors.New(err.Error() + " AND" + e.Error())
				}
			}
			time.Sleep(time.Duration(i) * time.Second)
			//log.Errorf(fmt.Sprintf("%v %v %v %v", err, i, query, args))
		}
	}
	return
}

func (dbserver *DBServer) QueryHelper(result interface{}, query string, args ...interface{}) (err error) {
	db := dbserver.ConnectionPool
	if db == nil {
		return errors.New("No connection pool aviailable")
	}

	var rows *sql.Rows
	//log.Debug(query, args)
	rows, err = db.Query(query, args...)
fmt.Println(rows)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Get columns name
	var columns []string
	columns, err = rows.Columns()
	if err != nil {
		return err
	}
	fields := make([]string, len(columns))
	for i, columnName := range columns {
		fields[i] = dbserver.firstCharToUpper(columnName)
	}
fmt.Println(fields)
	// Convert columns name to result struct exported field name and match each other
	// Not matched columns would be ignored.
	rv := reflect.ValueOf(result) // rv means result value
	if rv.Kind() == reflect.Ptr { // rv must by a slice pointer to send rows data back to caller
		rv = rv.Elem()
	} else {
fmt.Println(rv.Kind())
		return errors.New("Parameter result must be a slice pointer")
	}
	if rv.Kind() == reflect.Slice { // rv must be a slice to fill zero or more rows
		elemType := rv.Type().Elem()
		if elemType.Kind() == reflect.Struct { // pre row, pre struct
			ev := reflect.New(elemType)              // New slice struct element
			nv := reflect.MakeSlice(rv.Type(), 0, 0) // New slice for fill
			for rows.Next() {                        // for each rows
				scanArgs := make([]interface{}, len(fields))
				for i, fieldName := range fields {
					fv := ev.Elem().FieldByName(fieldName)
					if fv.Kind() != reflect.Invalid {
						scanArgs[i] = fv.Addr().Interface()
					} else {
						return errors.New("Invalid struct filed type or struct field `" + fieldName + "` does not exist. Query:" + query)
					}
				}
				err = rows.Scan(scanArgs...)
				if err != nil {
					return err
				}
				nv = reflect.Append(nv, ev.Elem())
			}
			rv.Set(nv) // return rows data back to caller
		}
	} else {
		return errors.New("Parameter result must be a slice pointer")
	}

	return
}

func (dbserver *DBServer) firstCharToUpper(str string) string {
	if   len(str) > 0 {
		runes := []rune(str)
		firstRune := unicode.ToUpper(runes[0])
		str = string(firstRune) + string(runes[1:])
	}
	return str
}











