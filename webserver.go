package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "encoding/json"
	"encoding/csv"
	"os"
    "io"
    "io/ioutil"
	"strconv"
	"bufio"
)

type ResponseMessage struct {
    Field1 string
    Field2 string
}

type RequestMessage struct {
    Field1 string
    Field2 string
}

type RentalMessage struct {
	CarMaker string
	CarModel string
	NDays int
	NUnits int
}


func main() {

router := mux.NewRouter().StrictSlash(true)
router.HandleFunc("/", Index)
router.HandleFunc("/endpoint/{param}", endpointFunc)
router.HandleFunc("/endpoint2/{param}", endpointFunc2JSONInput)
router.HandleFunc("/new_car_rental", newRental)
router.HandleFunc("/list_car_rental", listRental)

log.Fatal(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Service OK")
}

func endpointFunc(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    param := vars["param"]
    res := ResponseMessage{Field1: "Text1", Field2: param}
    json.NewEncoder(w).Encode(res)
}	

func endpointFunc2JSONInput(w http.ResponseWriter, r *http.Request) {
    var requestMessage RequestMessage
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        panic(err)
    }
    if err := r.Body.Close(); err != nil {
        panic(err)
    }
    if err := json.Unmarshal(body, &requestMessage); err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(422) // unprocessable entity
        if err := json.NewEncoder(w).Encode(err); err != nil {
            panic(err)
        }
    } else {
        fmt.Fprintln(w, "Successfully received request with Field1 =", requestMessage.Field1)
        fmt.Println(r.FormValue("queryparam1"))
    }
}

func newRental(w http.ResponseWriter, r *http.Request) {
	var rentalMessage RentalMessage
	var total_price int
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        panic(err)
    }
    if err := r.Body.Close(); err != nil {
        panic(err)
    }
    if err := json.Unmarshal(body, &rentalMessage); err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(422) // unprocessable entity
        if err := json.NewEncoder(w).Encode(err); err != nil {
            panic(err)
        }
    } else {
		total_price = rentalMessage.NDays * rentalMessage.NUnits * 20
		writeToFile(rentalMessage, w)
        fmt.Println(w, "Successfully received request with price: ", total_price)
    }
}

func listRental(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("rentals.csv")
    if err!=nil {
    json.NewEncoder(w).Encode(err)
    return
    }
    reader := csv.NewReader(bufio.NewReader(file))
    for {
        record, err := reader.Read()
        if err == io.EOF {
                break
            }
            fmt.Fprintf(w, "CarMaker: %q ", record[0])
			fmt.Fprintf(w, "CarModel: %q ", record[1])
			fmt.Fprintf(w, "Ndays: %q ", record[2])
			fmt.Fprintf(w, "NUnits: %q \n", record[3])
    }

}

func writeToFile(c RentalMessage, w http.ResponseWriter) {
    file, err := os.OpenFile("rentals.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
    if err!=nil {
        json.NewEncoder(w).Encode(err)
        return
    }
    writer := csv.NewWriter(file)
	var data = []string{c.CarMaker, c.CarModel, strconv.Itoa(c.NDays), strconv.Itoa(c.NUnits)}
    writer.Write(data)
    writer.Flush()
    file.Close()
}



