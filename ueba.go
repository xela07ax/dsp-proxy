package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var allRecords []Employee

func getItems(w http.ResponseWriter, r *http.Request) {
	log.Print("got /get-items request")
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	param := r.URL.RawQuery
	idsStr := strings.Split(param, ",")
	itemsRecords := make([]Employee, 0)
	for _, id := range idsStr {
		if id == "0" { // если 0, то выведем всех
			dat, err := json.Marshal(allRecords)
			if err != nil {
				panic(err)
			}
			w.Write(dat)
			return
		}
		itemId, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for _, each := range allRecords {
			if each.Id == itemId {
				itemsRecords = append(itemsRecords, each)
			}
		}
	}
	dat, err := json.Marshal(itemsRecords)
	if err != nil {
		panic(err)
	}
	w.Write(dat)
	return
}

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

type Employee struct {
	Id                                                   int
	Uid, Domain, Cn, Department, Title, Who, Logon_count string
}

func main() {
	records := readCsvFile("./ueba.csv")

	for _, each := range records {
		var oneRecord Employee
		oneRecord.Id, _ = strconv.Atoi(each[1])
		oneRecord.Uid = each[2]
		oneRecord.Domain = each[3]
		oneRecord.Cn = each[4]
		oneRecord.Department = each[5]
		oneRecord.Title = each[6]
		oneRecord.Who = each[7]
		oneRecord.Logon_count = each[8]
		allRecords = append(allRecords, oneRecord)
	}
	fmt.Println("Hello server at 3333\nhttp://localhost:3333/get-items?0\n")
	fmt.Println("http://localhost:3333/get-items?852,7656,133\n")
	http.HandleFunc("/get-items", getItems)

	err := http.ListenAndServe(":3333", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed ok\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}

	dat, err := json.Marshal(allRecords)
	if err != nil {
		panic(err)
	}
	fmt.Printf(">%s\n", dat)
	fmt.Println(records)
}
