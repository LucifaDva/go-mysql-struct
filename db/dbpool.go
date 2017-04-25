package db
import (
	"encoding/json"
	"errors"
	"time"
	"io/ioutil"
)

type DBPool struct {
	Servers 	map[string]*DBServer
	//dbMap 		util.RingMap
}


func (dbpool *DBPool) loadConfigFile(configFile string) (err error) {
	bytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		return
	}
	if err = json.Unmarshal(bytes, &dbpool.Servers); err != nil {
		return
	}
	return
}

func (dbpool *DBPool) Init(configFile string, serverNames []string, maxIdle int, maxDbConn int, reConnect bool) (err error) {
	//Load config file
	if err = dbpool.loadConfigFile(configFile); err != nil {
		return
	}
	if serverNames == nil {

	} else if len(serverNames) == 1 && serverNames[0] == "all" {
		// Init all database server connection
		for serverName, _ := range dbpool.Servers {
			_, err = dbpool.ConnectServer(serverName, maxIdle, maxDbConn, reConnect)
		}
	} else {
		for _, serverName:= range serverNames {
			_, err = dbpool.ConnectServer(serverName, maxIdle, maxDbConn, reConnect)
		}
	}
	//pool.dbMap.Init(5000, true)
	return
}

func (dbpool *DBPool) Clean() {
	for _, server := range dbpool.Servers {
		if server.ConnectionPool != nil {
			server.ConnectionPool.Close()
		}
	}
}

func (dbpool *DBPool) Close(serverName string) {
	server := dbpool.Servers[serverName]
	if server != nil && server.ConnectionPool != nil {
		server.ConnectionPool.Close()
	}
}
//???
func (dbpool *DBPool) Get(serverName string) (server DBServer, err error) {
	return dbpool.ConnectServer(serverName, server.MaxIdle, server.MaxDbConn, false)
}

func (dbpool *DBPool) ConnectServer(serverName string, maxIdle int, maxDbConn int, reConnect bool) (server DBServer, err error) {
	s, ok := dbpool.Servers[serverName]
	if !ok {
		err = errors.New("No configure for this db")
		return
	}
	
	s.mu.Lock()
	if time.Since(s.LastConnectTime) > time.Second {
		if s.ConnectionPool != nil {
			_, err = s.ConnectionPool.Exec("DO 1")
			if err != nil {
				reConnect = true
			} else {
				s.LastConnectTime = time.Now()
			}
		}
	}
	if reConnect || s.ConnectionPool == nil {
		err = s.conn(maxIdle, maxDbConn)
	}
	if err == nil {
		server = *s
	}
	s.mu.Unlock()
	return
 }











