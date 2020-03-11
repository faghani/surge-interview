package database

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

type Rule struct {
	Threshold int
	Coefficient float64
}

// Initialize rules
var Rules = map[int]float64{
	0: 1,
	1000: 1.05,
	3000: 1.10,
	5000: 1.20,
	10000: 1.50,
	25000: 2,
	50000: 3,
}

type Database struct {
	DB *sql.DB
}

func (d *Database) Setup() error {
	if !d.ShouldInsertRules() {
		return nil
	}
	log.Info("Importing initial rules")
	for threshold, coefficient := range Rules {
		st, err := d.DB.Prepare("INSERT INTO rule(threshold, coefficient) VALUES (?, ?)")
		if err != nil {
			return err
		}

		st.Exec(threshold, coefficient)
	}
	log.Info("Initial rules imported")
	return nil
}

func (d *Database) ShouldInsertRules() bool {
	var count int
	row := d.DB.QueryRow("SELECT COUNT(*) FROM rule")
	row.Scan(&count)
	return count == 0
}

func (d *Database) FetchRules(rules *[]Rule) error {
	rows, err := d.DB.Query("SELECT `threshold`, `coefficient` FROM `rule` ORDER BY `threshold` ASC")
	if err != nil {
		return err
	}
	for rows.Next()  {
		var (
			threshold int
			coefficient float64
		)
		err = rows.Scan(&threshold, &coefficient)
		if err != nil {
			return err
		}
		*rules = append(*rules, Rule{threshold, coefficient})
	}
	return nil
}

func (d *Database) SaveRules(rules []Rule) error {
	_, err := d.DB.Exec("TRUNCATE TABLE rule")
	if err != nil {
		return err
	}

	for _, rule := range rules {
		st, err := d.DB.Prepare("INSERT INTO rule(threshold, coefficient) VALUES (?, ?)")
		if err != nil {
			return err
		}

		st.Exec(rule.Threshold, rule.Coefficient)
	}

	return  nil
}

//ConnectToDatabase Connect to mysql database
func New(host, user, pass, name, port string) (*Database, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", user, pass, host, port, name)

	d, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = d.Ping()
	if err != nil {
		return nil, err
	}

	return &Database{DB: d}, nil
}