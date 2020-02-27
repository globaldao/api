package main

import (
  "github.com/gorilla/mux"
  "database/sql"
  _"github.com/go-sql-driver/mysql"
  "net/http"
  "encoding/json"
  "fmt"
  //"io/ioutil"
  "os"
)

type Pointrate struct {
  NAME string `json:"name"`
  VALUE string `json:"value"`
  LASTUPDATE string `json:"lastupdate"`
}

var db *sql.DB
var err error

func main() {

  //db, err = sql.Open("mysql", "juantellez:kd84kfi5@tcp(aliza.gootechnology.com:3306)/globaldaoadmin")
  db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", os.Getenv("DAOUSER"), os.Getenv("DAOPASSWORD"), os.Getenv("DAOHOST"), os.Getenv("DAODB")))

  if err != nil {
    panic(err.Error())
  }

  defer db.Close()

  router := mux.NewRouter()
  router.HandleFunc("/v1/public/pointrates", getPointrates).Methods("GET")
  http.ListenAndServe(":8000", router)

}

func getPointrates(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")

  var pointrates []Pointrate
  result, err := db.Query("SELECT name, value, lastupdate from pointrates")

  if err != nil {
    panic(err.Error())
  }

  defer result.Close()

  for result.Next() {
    var pointrate Pointrate

    err := result.Scan(&pointrate.NAME, &pointrate.VALUE, &pointrate.LASTUPDATE)

    if err != nil {
      panic(err.Error())
    }

    pointrates = append(pointrates, pointrate)
  }

  json.NewEncoder(w).Encode(pointrates)
}
