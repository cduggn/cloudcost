package storage

import (
	"database/sql"
	"github.com/cduggn/cloudcost/internal/pkg/logger"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"os"
)

var (
	createTableStmt = `
		CREATE TABLE cloudCostData (
		    till_id INTEGER PRIMARY KEY AUTOINCREMENT, 
		    dimension VARCHAR(64) NULL, 
		    dimension2 VARCHAR(255),
		    tag VARCHAR(255) NOT NULL, 
		    metric_name VARCHAR(255) NOT NULL,
		    amount FLOAT NOT NULL,
		    unit VARCHAR(255) NOT NULL,
		    granularity VARCHAR(255) NOT NULL,
		    start_date DATETIME NOT NULL,
		    end_date DATETIME NOT NULL)
    `
	insertStmt = "INSERT INTO cloudCostData (dimension, dimension2, tag, metric_name, amount, unit, granularity, start_date, end_date) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	conn       *CostDataStorage
)

func (c *CostDataStorage) New(dbName string) error {

	// if database already exists then return
	if conn != nil {
		return nil
	}

	// create the physical file if not already in place
	_, err := c.CreateFile(dbName)
	if err != nil {
		return DBError{msg: err.Error()}
	}

	// create the database name
	err = c.Set(dbName)
	if err != nil {
		logger.Error(err.Error())
		return DBError{msg: err.Error()}
	}

	// create the table
	res := c.createCostDataTable()
	if res == -1 {
		msg := "Could not create table"
		return DBError{msg: msg}
	}

	return nil
}

// return 0 if creation was a success or -1 if file was not created
func (c *CostDataStorage) CreateFile(dbName string) (int, error) {
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		file, err := os.Create(dbName)
		if err != nil {
			return -1, &DBError{msg: "Could not create database"}
		}
		defer file.Close()
		return 0, nil
	}
	return -1, nil
}

// create a database with named provided as arg
func (c *CostDataStorage) Set(s string) error {
	db, err := sql.Open("sqlite3", s)
	if err != nil {
		return err
	}
	c.SQLite = db
	return nil
}

// return -1 and or error if table was not created , return 0 if table was created
func (c *CostDataStorage) createCostDataTable() int {
	_, err := c.SQLite.Exec(createTableStmt)
	if err != nil {
		return -1
	}
	logger.Info("Table created", zap.String("table", "cloudCostData"))
	return 0
}

func (c *CostDataStorage) Insert(data CostDataInsert) int {

	stmt, err := c.SQLite.Prepare(insertStmt)
	if err != nil {
		logger.Error(err.Error())
		return -1
	}
	defer stmt.Close()

	res, err := stmt.Exec(data.Dimension, data.Dimension2, data.Tag, data.MetricName, data.Amount, data.Unit, data.Granularity, data.StartDate, data.EndDate)
	if err != nil {
		logger.Error(err.Error())
		return -1
	}
	id, err := res.LastInsertId()
	if err != nil {
		logger.Error(err.Error())
		return -1
	}
	logger.Info("Row added", zap.Int64("rowId", id))
	return 0
}
