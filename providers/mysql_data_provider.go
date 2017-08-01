package providers

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

type MySqlDataProvider struct {
	db *sql.DB
}

func (p *MySqlDataProvider) Initialize(ips []string, env_id int, otherargs string, needInitData bool) error {
	if otherargs == "" {
		return errors.New("The parameter: \"mysql-connection-sql\" MUST be set.")
	}
	err := p.initDB(otherargs)
	if err != nil {
		return fmt.Errorf("Failed initializing to given MYSQL database, error: %s", err.Error())
	}
	if needInitData {
		if ips == nil {
			return errors.New("The parameter: \"ips\" MUST be set.")
		}
		err = p.SaveIPs(ips, env_id)
		if err != nil {
			return fmt.Errorf("Failed initializing data to given MYSQL database, error: %s", err.Error())
		}
	}
	return nil
}

func (p *MySqlDataProvider) initDB(connStr string) error {
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return err
	}
	// file := "mysql_database.sql"
	// if _, err = os.Stat(file); err != nil {
	// 	return fmt.Errorf("MYSQL database restoring file DOES NOT existed. file: %s, error: %s", file, err.Error())
	// }
	// fileData, err := ioutil.ReadFile(file)
	// if err != nil {
	// 	return fmt.Errorf("Failed reading MYSQL database restoring file. file: %s, error: %s", file, err.Error())
	// }
	// log.Info("  Restoring MYSQL database...")
	// _, err = db.Query(string(fileData), nil)
	// if err != nil {
	// 	return fmt.Errorf("Failed restoring MYSQL database restoring file. file: %s, error: %s", file, err.Error())
	// }
	p.db = db
	return nil
}

func (p *MySqlDataProvider) SaveIPs(ips []string, env_id int) error {
	batchSize := 1000
	batchCount := 0
	if len(ips)%batchSize == 0 {
		batchCount = len(ips) / batchSize
	} else {
		batchCount = int(len(ips)/batchSize) + 1
	}
	log.Infof("  Forcasted IP pool size = %d", len(ips))
	log.Infof("  Defined database init batch size = %d", batchSize)
	log.Infof("  Calculated database init batch count = %d", batchCount)
	for i := 0; i < batchCount; i++ {
		log.Infof("    Starting INIT database, batch index: %d", i)
		sqlStr := "INSERT INTO ip_pool(ip, source_ip, env_id) VALUES "
		vals := []interface{}{}
		var batchData []string = nil
		if rst := len(ips) - (i * batchCount); rst >= batchSize {
			batchData = ips[i*batchSize : (i*batchSize)+batchSize]
		} else {
			batchData = ips[i*batchSize : (i*batchSize)+rst]
		}
		log.Infof("    Current batch data size = %d", len(batchData))
		for _, ip := range batchData {
			if ip != "" {
				sqlStr += "(?, ?, ?),"
				vals = append(vals, ip, "", env_id)
			}
		}
		//log.Info(sqlStr)
		//trim the last
		sqlStr = sqlStr[0 : len(sqlStr)-1]
		//log.Info(sqlStr)
		//prepare the statement
		stmt, err := p.db.Prepare(sqlStr)
		if err != nil {
			return err
		}
		defer stmt.Close()
		//format all vals at once
		_, err = stmt.Exec(vals...)
		if err != nil {
			log.Infof("%#v", vals)
			return err
		}
		log.Infof("    Batch %d successfully committed.", i)
	}
	log.Info("  All data batches had been committed successfully!")
	return nil
}

func (p *MySqlDataProvider) GetAvailableIP(nodeIp string, envId int, owner, desc string) (string, error) {
	sqlStr := "CALL USP_CheckinIP(?, ?, ?, ?)"
	stmt, err := p.db.Prepare(sqlStr)
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	rst, err := stmt.Query(nodeIp, envId, owner, desc)
	if err != nil {
		return "", err
	}
	defer rst.Close()
	if !rst.Next() {
		return "", errors.New("No more available IP subnet can be assigns.")
	}
	if err = rst.Err(); err != nil {
		return "", err
	}
	var prviateSubnet string
	err = rst.Scan(&prviateSubnet)
	if err != nil {
		return "", errors.New("No more available IP subnet can be assigns by given conditions.")
	}
	return prviateSubnet, nil
}
