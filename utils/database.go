package utils

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	_ "github.com/lib/pq" //pgadmin
)

//	db.QueryRow("INSERT INTO hash(hash,source, first_seen) VALUES($1,$2,$3) ;", hash, "dht", time.Now().Local().Format("2006-01-02"))

//Config from file for the database
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
}

//GetConfig Get the config from the configFile
func GetConfig(configFile string, debug bool, verbose bool) (Config, error) {
	var config Config

	if debug {
		log.Println("Getting database config from  ", configFile)
	}

	_, err := os.Stat(configFile)
	if err != nil {
		return config, err
	}

	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		log.Fatal(err)
	}

	return config, nil

}

//ConnectDb : conect database and return *sql.DB
func ConnectDb(cfgFile string, debug bool, verbose bool) (*sql.DB, error) {

	if debug {
		log.Println("Conecting DB")
	}
	config, err := GetConfig(cfgFile, debug, verbose)
	if err != nil {
		return nil, err
	}

	errdb := CreateDb(cfgFile, config.Dbname, debug, verbose)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	if errdb == nil {
		err = InitDB(db, debug, verbose)
		if err != nil {
			log.Fatal(err)
		}
	}
	return db, err
}

//CreateDb : Create database hash
func CreateDb(cfgFile string, database string, debug bool, verbose bool) error {

	if debug {
		log.Println("Creating DB  ", cfgFile)
	}

	config, err := GetConfig(cfgFile, debug, verbose)
	if err != nil {
		return err
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}

	//Create DB
	_, err = db.Exec(`create database ` + database)
	if err != nil {
		return err
	}

	return nil
}

//InitDB : init Database
func InitDB(db *sql.DB, debug bool, verbose bool) error {

	if debug {
		log.Println("Init DB")
	}

	//read sql file
	gopath := os.Getenv("GOPATH")

	b, err := ioutil.ReadFile(gopath + "/src/github.com/RaulCalvoLaorden/bntoolkit/sql.sql") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	databaseSQL := string(b) // convert content to a 'string'

	//Init db
	_, err = db.Exec(databaseSQL)
	if debug {
		fmt.Println(databaseSQL)
	}
	if err != nil {
		return err
	}
	return nil
}

//ExecuteDb execeute a SQL query. Dangerous
func ExecuteDb(db *sql.DB, debug bool, verbose bool, sqlString string) error {
	//Create DB
	if debug {
		log.Println("Executing in SQL:  ", sqlString)
	}
	_, err := db.Exec(sqlString)
	if err != nil {
		return err
	}
	return nil
}

//InsertHash to the database
func InsertHash(db *sql.DB, debug bool, verbose bool, hash string, source string) error {

	sqlStatement := `
		INSERT INTO hash(hash,source, first_seen) 
		VALUES($1,$2,$3)
		RETURNING hash`

	if debug {
		log.Println("Insert hash  ", hash)
	}

	id := ""
	err := db.QueryRow(sqlStatement, hash, source, time.Now().Local()).Scan(&id)
	if err != nil {
		return err
	}
	if debug {
		log.Println("New record ID is:", id)
	}
	return nil
}

//InsertProject to the database
func InsertProject(db *sql.DB, debug bool, verbose bool, projectName string) error {

	sqlStatement := "INSERT INTO project (\"projectName\", date) VALUES('" + projectName + "',current_timestamp) RETURNING \"projectName\""
	id := ""

	if debug {
		log.Printf("Insert project %v ", projectName)
	}

	err := db.QueryRow(sqlStatement).Scan(&id)
	if err != nil {
		return err
	}
	if debug {
		log.Println("New record ID is:", id)
	}
	return nil
}

//DeleteProject from the database
func DeleteProject(db *sql.DB, debug bool, verbose bool, name string) error {

	sqlStatement := `
		DELETE FROM project where name = $1`

	if debug {
		log.Println("Delete project  ", name)
	}
	db.QueryRow(sqlStatement, name)

	return nil
}

//InsertMonitor inserts a hash in teh database for monitor in daemon mode.
func InsertMonitor(db *sql.DB, debug bool, verbose bool, hash string, username string, projectName string) error {

	sqlStatement := `
		INSERT INTO monitor(hash, username, "projectName")
		VALUES($1,$2,$3)
		RETURNING hash`
	id := ""

	if debug {
		log.Println("Insert monitor ", hash, username, projectName)
	}

	err := db.QueryRow(sqlStatement, hash, username, projectName).Scan(&id)
	if err != nil {
		return err
	}
	if debug {
		log.Println("New record ID is:", id)
	}
	return nil
}

//DeleteMonitor from the database
func DeleteMonitor(db *sql.DB, debug bool, verbose bool, hash string) error {

	sqlStatement := `
		DELETE FROM monitor where hash = $1`

	if debug {
		log.Println("Delete monitor  ", hash)
	}

	db.QueryRow(sqlStatement, hash)

	return nil
}

//InsertAlert to the database
func InsertAlert(db *sql.DB, debug bool, verbose bool, ip string, list string, username string, projectName string) error {
	sqlStatement := `
		INSERT INTO alert(ip, list, username, "projectName")
		VALUES($1,$2,$3,$4)
		RETURNING ip`
	id := ""

	if debug {
		log.Println("Insert alert ", ip, username, projectName)
	}

	err := db.QueryRow(sqlStatement, ip, list, username, projectName).Scan(&id)
	if err != nil {
		return err
	}
	if debug {
		log.Println("New record ID is:", id)
	}
	return nil
}

//DeleteAlert from the database
func DeleteAlert(db *sql.DB, debug bool, verbose bool, ip string) error {
	sqlStatement := "delete from alert where ip <<= $1"

	if debug {
		log.Println("delete alert", ip)
	}

	db.QueryRow(sqlStatement, ip)

	return nil
}

//InsertIP inserts a IP or range in the database db with the proyectName
func InsertIP(db *sql.DB, debug bool, verbose bool, ip string, projectName string) error {

	sqlStatement := `
		INSERT INTO ip(ip, "projectName")
		VALUES($1,$2)
		RETURNING ip`
	id := ""

	if debug {
		log.Println("Insert ip  ", ip, projectName)
	}

	err := db.QueryRow(sqlStatement, ip, projectName).Scan(&id)
	if err != nil {
		return err
	}
	if debug {
		log.Println("New record ID is:", id)
	}
	return nil
}

//InsertDownload to the database
func InsertDownload(db *sql.DB, debug bool, verbose bool, ip string, port int, hash string, projectName string) error {

	sqlStatement := `
		INSERT INTO download(ip, port, hash, date, "projectName")
		VALUES($1,$2,$3,current_timestamp,$4)
		RETURNING ip`
	id := ""

	if debug {
		log.Println("Insert download", ip, port, hash, projectName)
	}

	err := db.QueryRow(sqlStatement, ip, port, hash).Scan(&id)
	if err != nil {
		return err
	}
	if debug {
		log.Println("New record ID is:", id)
	}
	return nil
}

//InsertHashList to the database
func InsertHashList(db *sql.DB, debug bool, verbose bool, hashes []string, source string) error {
	sql := `INSERT INTO hash(hash ,source, first_seen) VALUES `
	max := len(hashes)

	for i := 0; i < max; i++ {
		value := hashes[i]
		if !strings.Contains(sql, value) {
			if i == 0 {
				sql += "\n ('" + value + "'," + " '" + source + "' , " + "current_timestamp)"
			} else {
				sql += ",\n('" + value + "'," + " '" + source + "' , " + "current_timestamp)"
			}
		}
		fmt.Println(value)
	}
	sql += " ON CONFLICT DO NOTHING" //UPDATE set last_seen = '" + timeAux + "'::timestamp; "

	if debug {
		log.Println("Insert hashList  ", sql)
	}
	ExecuteDb(db, debug, verbose, sql)

	return nil
}

//InsertFile  to the database
func InsertFile(db *sql.DB, debug bool, verbose bool, filePath string) error {

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var hashList []string

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		hashList = append(hashList, scanner.Text())
		if count == 1000 {
			InsertHashList(db, debug, verbose, hashList, filePath)
			var hashList2 []string
			hashList = hashList2
			count = 0
		}
		count++
	}
	InsertHashList(db, debug, verbose, hashList, filePath)

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

//GetHash from ID
func GetHash(db *sql.DB, debug bool, verbose bool, i int) (string, error) {

	sqlStatement := `SELECT hash FROM possibles WHERE id=$1;`

	if debug {
		log.Println("Get hash  ", i)
	}

	var hash string
	// Replace 3 with an ID from your database or another random
	// value to test the no rows use case.
	row := db.QueryRow(sqlStatement, i)
	switch err := row.Scan(&hash); err {
	case sql.ErrNoRows:
		return "", err
	}
	return hash, nil
}

//GetPossibles from DB possibles table
func GetPossibles(db *sql.DB, debug bool, verbose bool) (int, error) {

	var count int

	if debug {
		log.Println("Getting possibles")
	}

	row := db.QueryRow("SELECT COUNT(*) FROM possibles")
	err := row.Scan(&count)
	if err != nil {
		return -1, err
	}

	return count, nil
}

//DeletePossibles table
func DeletePossibles(db *sql.DB, debug bool, verbose bool) {

	if debug {
		log.Println("Deelete possibles")
	}

	db.QueryRow("DELETE FROM public.possibles")
}

//DeletePossiblesFalse where possible=false
func DeletePossiblesFalse(db *sql.DB, debug bool, verbose bool) {

	if debug {
		log.Println("Delete possibles false")
	}

	db.QueryRow("DELETE FROM public.possibles where Possible=false")
}

//SetTruePossible a hash from possibles table
func SetTruePossible(db *sql.DB, debug bool, verbose bool, hash string) error {
	var id int
	//db.QueryRow("UPDATE public.possibles SET valid=True WHERE id=$1;")
	sqlStatement := `UPDATE public.possibles SET Possible=True WHERE hash=$1 RETURNING id;`
	err := db.QueryRow(sqlStatement, hash).Scan(&id)

	if debug {
		log.Println("Set possible true  ", hash)
	}
	if err != nil {
		return err
	}
	if debug {
		log.Println("New record ID is:", id)
	}

	return nil
}

//SetTrueDownload a hash from possibles table
func SetTrueDownload(db *sql.DB, debug bool, verbose bool, hash string) error {
	var id int
	//db.QueryRow("UPDATE public.possibles SET valid=True WHERE id=$1;")
	sqlStatement := `UPDATE public.possibles SET download=True WHERE hash=$1 RETURNING id;`

	if debug {
		log.Println("Set download true  ", hash)
	}

	err := db.QueryRow(sqlStatement, hash).Scan(&id)
	if err != nil {
		return err
	}
	if debug {
		log.Println("New record ID is:", id)
	}

	return nil
}

//SetLen to a hash in possibles table
func SetLen(db *sql.DB, debug bool, verbose bool, num int, hash string) error {

	var id int
	//db.QueryRow("UPDATE public.possibles SET valid=True WHERE id=$1;")
	sqlStatement := `UPDATE public.possibles SET num=$1 WHERE hash=$2 RETURNING id;`

	if debug {
		log.Println("Set len  ", hash)
	}

	err := db.QueryRow(sqlStatement, num, hash).Scan(&id)
	if err != nil {
		return err
	}
	if debug {
		log.Println("New record ID is:", id)
	}

	return nil
}

//SetTrueValid in possibles table
func SetTrueValid(db *sql.DB, debug bool, verbose bool, hash string) error {

	var id int
	//db.QueryRow("UPDATE public.possibles SET valid=True WHERE id=$1;")
	sqlStatement := `UPDATE public.possibles SET valid=True WHERE hash=$1 RETURNING id;`

	if debug {
		log.Println("Set true valid  ", hash)
	}

	err := db.QueryRow(sqlStatement, hash).Scan(&id)
	if err != nil {
		return err
	}
	if debug {
		log.Println("New record ID is:", id)
	}

	return nil
}

//SetNamePossibles in possibles table
func SetNamePossibles(db *sql.DB, debug bool, verbose bool, name string, hash string) error {

	var id int
	//db.QueryRow("UPDATE public.possibles SET valid=True WHERE id=$1;")
	sqlStatement := `UPDATE public.possibles SET name=$1 WHERE hash=$2 RETURNING id;`

	if debug {
		log.Println("Set name possibles", name, hash)
	}

	err := db.QueryRow(sqlStatement, name, hash).Scan(&id)
	if err != nil {
		return err
	}
	if debug {
		log.Println("New record ID is:", id)
	}

	return nil
}

//InsertPossible in possibles table
func InsertPossible(db *sql.DB, debug bool, verbose bool, num int, hash string, valid bool, projectName string) error {

	possible := true
	download := false
	sqlStatement := `
		INSERT INTO public.possibles(id, hash, download, valid, possible, num, "projectName") 
		VALUES($1,$2,$3,$4,$5,$6,$7)
		RETURNING hash`
	id := ""
	db.QueryRow(sqlStatement, num, hash, download, valid, possible, 0, projectName).Scan(&id)

	if debug {
		log.Println(sqlStatement)
		log.Println("New record ID is:", id)
	}

	if id == "" {
		return errors.New("DB Insert error")
	}

	return nil
}

//CheckExist : comprobar si el Possible esta en la otra lista
func CheckExist(db *sql.DB, debug bool, verbose bool, hash string) (bool, error) {
	var count int
	row := db.QueryRow("SELECT count(*)	FROM public.hash where hash = '" + hash + "';")
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	if count != 0 {
		return true, nil
	}
	return false, nil

}

var waitGroup sync.WaitGroup

//DownloadPossibles Get and download (torrentLib) from possibles where possible true
func DownloadPossibles(db *sql.DB, timeout int, debug bool, verbose bool, projectName string) error {
	rows, err := db.Query("SELECT id, hash, download, valid, possible, num FROM possibles where Possible=True")
	if err != nil {
		return err
	}
	for rows.Next() {
		var id int
		var hash string
		var download bool
		var valid bool
		var possible bool
		var num int
		err := rows.Scan(&id, &hash, &download, &valid, &possible, &num)
		if err != nil {
			return err
		}
		waitGroup.Add(1)
		go torrentLib(db, timeout, debug, verbose, hash)
	}
	waitGroup.Wait()
	//time.Sleep(time.Minute)
	return nil
}

//DownloadValid Get and download (torrentLib) from possibles where valid = true
func DownloadValid(db *sql.DB, timeout int, debug bool, verbose bool, projectName string) error {
	rows, err := db.Query("SELECT id, hash, download, valid, possible, num FROM possibles where Valid=True")
	if err != nil {
		return err
	}
	for rows.Next() {
		var id int
		var hash string
		var download bool
		var valid bool
		var possible bool
		var num int
		err := rows.Scan(&id, &hash, &download, &valid, &possible, &num)
		if err != nil {
			return err
		}
		waitGroup.Add(1)
		go torrentLib(db, timeout, debug, verbose, hash)
	}
	waitGroup.Wait()
	//time.Sleep(time.Minute)
	return nil
}

//SelectPossiblesWhere columnTrue=True from the database. Possibles values: download, valid and possible
func SelectPossiblesWhere(columnTrue string, db *sql.DB, debug bool, verbose bool, projectName string) error {
	rows, err := db.Query("SELECT id, hash, download, valid, possible, num, \"projectName\" FROM possibles where " + columnTrue + "=True")
	if err != nil {
		return err
	}
	fmt.Println(" id \t|\t\t\t hash \t\t\t\t|\t download \t|\t valid \t\t|\t Possible \t|\t num \t\t|\t projectName")
	for rows.Next() {
		var id int
		var hash string
		var download bool
		var valid bool
		var possible bool
		var num int
		var projectName string
		err := rows.Scan(&id, &hash, &download, &valid, &possible, &num, &projectName)
		if err != nil {
			return err
		}
		fmt.Printf("%3v \t|\t %8v \t|\t %6v \t|\t %6v \t|\t %6v \t|\t %6v \t|\t %8v \n", id, hash, download, valid, possible, num, projectName)

	}
	//time.Sleep(time.Minute)
	return nil
}

//GetHashes from database
func GetHashes(db *sql.DB, debug bool, verbose bool) ([][]byte, error) {
	var infohashes [][]byte
	max, err := GetPossibles(db, debug, verbose)
	if err != nil {
		return nil, err
	}

	for i := 0; i < max; i++ {
		hash, err := GetHash(db, debug, verbose, i)
		if err != nil {
			return nil, err
		}

		infohashes = append(infohashes, []byte(hash))
	}

	return infohashes, nil
}

//QueryHash from database
func QueryHash(db *sql.DB, debug bool, verbose bool, sqlQuery string, hash string, source string) (string, error) {
	salida := ""

	where := ""
	if sqlQuery != "" {
		where = sqlQuery
	} else if hash != "" {
		where = " where hash='" + hash + "'"
		if source != "" {
			where += " AND source='" + source + "'"
		}
	} else if source != "" {
		where = " where source='" + source + "'"
	}
	//fmt.Println("SELECT hash , source , first_seen FROM hash " + where + ";")
	rows, err := db.Query("SELECT hash , source , first_seen, path, name FROM hash " + where + ";")

	if err != nil {
		return "", err
	}
	salida += "\t\t\t hash \t\t\t\t|\t source \t\t|\t\t first_seen \t\t\t|\t\t path \t\t|\t\t name " + "\n"
	//fmt.Println("hash \t\t\t\t\t\t\t\t|\t source \t|\t first_seen")
	for rows.Next() {
		var hash string
		var source string
		var firstSeen string
		var path sql.NullString
		var name sql.NullString
		err := rows.Scan(&hash, &source, &firstSeen, &path, &name)
		if err != nil {
			return "", err
		}
		salida += "\t" + hash + " \t|\t   " + source + "   \t\t|\t " + firstSeen + "\t\t|\t\t" + path.String + "\t\t|\t\t" + name.String + "\n"
		//fmt.Printf("%3v \t|\t %8v \t|\t %6v \n", hash, source, first_seen)
	}
	return salida, nil
}

//QueryPossibles from the database
func QueryPossibles(db *sql.DB, debug bool, verbose bool, sql string, hash string) (string, error) {
	salida := ""

	where := ""
	if sql != "" {
		where = sql
	} else if hash != "" {
		where = " where hash='" + hash + "'"
	}
	//fmt.Println("SELECT hash , source , first_seen FROM hash " + where + ";")
	rows, err := db.Query("SELECT id, hash, download, valid, possible, num, \"projectName\" FROM possibles " + where + ";")

	if err != nil {
		return "", err
	}
	salida += (" id \t|\t\t\t hash \t\t\t\t|\t download \t|\t valid \t\t|\t Possible \t|\t num \t\t|\t projectName\n")
	for rows.Next() {
		var id int
		var hash string
		var download bool
		var valid bool
		var possible bool
		var num int
		var projectName string
		err := rows.Scan(&id, &hash, &download, &valid, &possible, &num, &projectName)
		if err != nil {
			return "", err
		}
		salida += fmt.Sprintf("%3v \t|\t %8v \t|\t %6v \t|\t %6v \t|\t %6v \t|\t %6v \t|\t %8v \n", id, hash, download, valid, possible, num, projectName)

	}
	//time.Sleep(time.Minute)
	return salida, nil
}

//QueryProjects from  the database projects table
func QueryProjects(db *sql.DB, debug bool, verbose bool, sql string, nombre string) (string, error) {
	salida := ""

	where := ""
	if sql != "" {
		where = sql
	} else if nombre != "" {
		where = " where nombre='" + nombre + "'"
	}
	//fmt.Println("SELECT hash , source , first_seen FROM hash " + where + ";")
	rows, err := db.Query("SELECT \"projectName\", date FROM project " + where + ";")

	if err != nil {
		return "", err
	}
	salida += "projectName \t|\t date" + "\n"
	//fmt.Println("hash \t\t\t\t\t\t\t\t|\t source \t|\t first_seen")
	for rows.Next() {
		var nombre string
		var fecha string
		err := rows.Scan(&nombre, &fecha)
		if err != nil {
			return "", err
		}
		salida += nombre + " \t|\t " + fecha + "\n"
		//fmt.Printf("%3v \t|\t %8v \t|\t %6v \n", hash, source, first_seen)
	}
	return salida, nil
}

//QueryIP SELECT from the database ip table
func QueryIP(db *sql.DB, debug bool, verbose bool, sql string, ip string) (string, error) {
	salida := ""

	where := ""
	if sql != "" {
		where = sql
	} else if ip != "" {
		where = " where ip='" + ip + "'"
	}
	//fmt.Println("SELECT hash , source , first_seen FROM hash " + where + ";")
	rows, err := db.Query("SELECT ip FROM ip " + where + ";")

	if err != nil {
		return "", err
	}
	salida += "ip" + "\n"
	//fmt.Println("hash \t\t\t\t\t\t\t\t|\t source \t|\t first_seen")
	for rows.Next() {
		var ip string
		err := rows.Scan(&ip)
		if err != nil {
			return "", err
		}
		salida += ip + "\n"
		//fmt.Printf("%3v \t|\t %8v \t|\t %6v \n", hash, source, first_seen)
	}
	return salida, nil
}

//QueryMonitor SELECT from the database monitor table
func QueryMonitor(db *sql.DB, debug bool, verbose bool, sqlQuery string, hash string, user string) (string, error) {
	salida := ""

	where := ""
	if sqlQuery != "" {
		where = sqlQuery
	} else if hash != "" {
		where = " where hash='" + hash + "'"
		if user != "" {
			where += " AND user='" + user + "'"
		}
	} else if user != "" {
		where = " where user='" + user + "'"
	}
	//fmt.Println("SELECT hash , source , first_seen FROM hash " + where + ";")
	rows, err := db.Query("SELECT hash, username, \"projectName\" FROM monitor " + where + ";")

	if err != nil {
		return "", err
	}
	salida += "hash \t\t\t\t\t\t|\t username \t|\t projectName" + "\n"
	//fmt.Println("hash \t\t\t\t\t\t\t\t|\t source \t|\t first_seen")
	for rows.Next() {
		var hash sql.NullString
		var user sql.NullString
		var projectName sql.NullString

		err := rows.Scan(&hash, &user, &projectName)
		if err != nil {
			return "", err
		}
		salida += hash.String + " \t|\t " + user.String + " \t|\t " + projectName.String + "\n"
		//fmt.Printf("%3v \t|\t %8v \t|\t %6v \n", hash, source, first_seen)
	}
	return salida, nil
}

//QueryCount from the database hash table
func QueryCount(db *sql.DB, debug bool, verbose bool) (string, error) {
	salida := ""

	//fmt.Println("SELECT hash , source , first_seen FROM hash " + where + ";")
	rows, err := db.Query("SELECT count(*) FROM hash " + ";")

	if err != nil {
		return "", err
	}
	salida += "count" + "\n"
	//fmt.Println("hash \t\t\t\t\t\t\t\t|\t source \t|\t first_seen")
	for rows.Next() {
		var count string
		err := rows.Scan(&count)
		if err != nil {
			return "", err
		}
		salida += count + "\n"
		//fmt.Printf("%3v \t|\t %8v \t|\t %6v \n", hash, source, first_seen)
	}
	return salida, nil
}

//QueryAlert from the database
func QueryAlert(db *sql.DB, debug bool, verbose bool, sql string, ip string, user string) (string, error) {
	salida := ""

	where := ""
	if sql != "" {
		where = sql
	} else if ip != "" {
		where = " where ip='" + ip + "'"
		if user != "" {
			where += " AND user='" + user + "'"
		}
	} else if user != "" {
		where = " where user='" + user + "'"
	}
	if verbose {
		fmt.Println("SELECT hash , source , first_seen FROM hash " + where + ";")
	}
	rows, err := db.Query("SELECT ip, list, username FROM alert " + where + ";")

	if err != nil {
		return "", err
	}
	salida += "ip \t\t|\t list \t|\t user" + "\n"
	//fmt.Println("hash \t\t\t\t\t\t\t\t|\t source \t|\t first_seen")
	for rows.Next() {
		var ip string
		var user string
		var list string

		err := rows.Scan(&ip, &user, &list)
		if err != nil {
			return "", err
		}
		salida += ip + " \t|\t " + user + " \t|\t " + list + "\n"
		//fmt.Printf("%3v \t|\t %8v \t|\t %6v \n", hash, source, first_seen)
	}
	return salida, nil
}

//GetMonitor from the database
func GetMonitor(db *sql.DB, debug bool, verbose bool, projectName string) ([]string, error) {
	var infohashes []string
	rows, err := db.Query("SELECT hash FROM monitor where \"projectName\"='" + projectName + "'")
	if err != nil {
		// handle this error better than this
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var hash string
		err = rows.Scan(&hash)
		if err != nil {
			// handle this error
			panic(err)
		}
		if debug {
			fmt.Println(hash)
		}

		infohashes = append(infohashes, hash)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return infohashes, nil
}
