package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Vishnukvsvk/Metrorail/dbutils"
	"github.com/emicklei/go-restful"
	_ "github.com/mattn/go-sqlite3"
)

//DB driver visible to whole program
var db *sql.DB

//Train Resource is model for holding rail info
type TrainResource struct {
	ID              int
	DriverName      string
	OperatingStatus bool
}

//Station Resource is model for holding station info
type StationResource struct {
	ID          int
	Name        string
	OpeningTime time.Time
	ClosingTime time.Time
}

//Schedule resource links train and station
type ScheduleResource struct {
	ID          int
	TrainId     int
	StationId   int
	ArrivalTime time.Time
}

//Register adds path and routes to a new service
func (t *TrainResource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/v1/trains").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)

	//getTrain,createTrain,deleteTrain are function handlers
	ws.Route(ws.GET("/{train-id}").To(t.getTrain))
	ws.Route(ws.POST("").To(t.createTrain))
	ws.Route(ws.DELETE("/{train-id}").To(t.deleteTrain))
	container.Add(ws)
}

// GET http://localhost:8000/v1/trains/1
func (t TrainResource) getTrain(req *restful.Request, resp *restful.Response) {
	id := req.PathParameter("train-id")
	err := db.QueryRow("SELECT ID, DRIVER_NAME, OPERATING_STATUS FROM train WHERE ID=?", id).Scan(&t.ID, &t.DriverName, &t.OperatingStatus)
	if err != nil {
		log.Println(err)
		resp.AddHeader("Content-Type", "text/plain")
		resp.WriteErrorString(http.StatusNotFound, "Train could not be found.")
	} else {
		resp.WriteEntity(t)
	}
}

// POST http://localhost:8000/v1/trains
func (t TrainResource) createTrain(req *restful.Request, resp *restful.Response) {
	log.Println(req.Request.Body)
	decoder := json.NewDecoder(req.Request.Body)
	var b TrainResource
	err := decoder.Decode(&b)
	log.Println(b.DriverName, b.OperatingStatus)
	if err != nil {
		log.Println(err)
	}
	statement, _ := db.Prepare("INSERT INTO TRAIN (DRIVER_NAME,OPERATING_STATUS)VALUES(?,?)")
	result, err := statement.Exec(b.DriverName, b.OperatingStatus)
	if err == nil {
		newID, _ := result.LastInsertId()
		b.ID = int(newID)
		resp.WriteHeaderAndEntity(http.StatusCreated, b)
	} else {
		resp.AddHeader("Content-Type", "text/plain")
		resp.WriteErrorString(http.StatusInternalServerError, err.Error())
	}
}

func (t TrainResource) deleteTrain(req *restful.Request, resp *restful.Response) {
	id := req.PathParameter("train-id")
	statement, _ := db.Prepare("delete from train where id=?")
	_, err := statement.Exec(id)
	if err == nil {
		resp.WriteHeader(http.StatusOK)
	} else {
		resp.AddHeader("Content-Type", "text/plain")
		resp.WriteErrorString(http.StatusInternalServerError,
			err.Error())
	}
}

func main() {
	//Connect to database
	db, err := sql.Open("sqlite3", "./railapi.db")
	if err != nil {
		log.Println("driver creation failed")
	}

	//Create tables
	dbutils.Initialize(db)
	wsContainer := restful.NewContainer()
	wsContainer.Router(restful.CurlyRouter{})
	t := TrainResource{}
	t.Register(wsContainer)
	log.Printf("start listening on localhost:8000")
	server := &http.Server{Addr: ":8000", Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}
