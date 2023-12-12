package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type RailwaysDB struct {
	db *sql.DB
}

var reader = bufio.NewReader(os.Stdin)

// MySQL [jail] supported isolation levels
var levelsMap = map[int]sql.IsolationLevel{
	1: sql.LevelReadUncommitted,
	2: sql.LevelReadCommitted,
	3: sql.LevelRepeatableRead,
	4: sql.LevelSerializable,
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	connectionString := os.Getenv("CONNECTION_STRING")

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		fmt.Println("Error opening database: ", err)
		return
	}
	defer db.Close()

	rdb := RailwaysDB{db: db}

	fmt.Println("Welcome to the Railways Database!")
	fmt.Println("for help, type 'h'")

	for {
		fmt.Print(">> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		switch input {
		case "h":
			help()
		case "e", "q":
			return
		case "cs":
			station, err := rdb.createStation()
			if err != nil {
				fmt.Println("Error creating station: ", err)
			} else {
				fmt.Println("Station created:")
				fmt.Println("Id City Name")
				fmt.Println(station.Id, station.City, station.Name)
			}
		case "gss":
			stations, err := rdb.getAllStations()
			if err != nil {
				fmt.Println("Error getting stations: ", err)
			} else {
				fmt.Println("Stations:")
				fmt.Println("Id City Name")
				for _, station := range *stations {
					fmt.Println(station.Id, station.City, station.Name)
				}
			}
		case "gs":
			id := readInt("Enter station id: ")
			station, err := rdb.getStationById(int64(id))
			if err != nil {
				fmt.Println("Error getting station: ", err)
			} else {
				fmt.Print("Station:")
				fmt.Println("Id City Name")
				fmt.Println(station.Id, station.City, station.Name)
			}
		case "us":
			station, err := rdb.updateStation()
			if err != nil {
				fmt.Println("Error updating station: ", err)
			} else {
				fmt.Println("Station updated:")
				fmt.Println("Id City Name")
				fmt.Println(station.Id, station.City, station.Name)
			}
		case "ds":
			id := readInt("Enter station id: ")
			err := rdb.deleteStation(int64(id))
			if err != nil {
				fmt.Println("Error deleting station: ", err)
			} else {
				fmt.Println("Station deleted")
			}
		case "ct":
			train, err := rdb.createTrain()
			if err != nil {
				fmt.Println("Error creating train: ", err)
			} else {
				fmt.Println("Train created:")
				fmt.Println("Id Code Capacity")
				fmt.Println(train.Id, train.Code, train.Capacity)
			}
		case "gts":
			trains, err := rdb.getAllTrains()
			if err != nil {
				fmt.Println("Error getting trains: ", err)
			} else {
				fmt.Println("Trains:")
				fmt.Println("Id Code Capacity")
				for _, train := range *trains {
					fmt.Println(train.Id, train.Code, train.Capacity)
				}
			}
		case "gt":
			id := readInt("Enter train id: ")
			train, err := rdb.getTrainById(int64(id))
			if err != nil {
				fmt.Println("Error getting train: ", err)
			} else {
				fmt.Print("Train:")
				fmt.Println("Id Code Capacity")
				fmt.Println(train.Id, train.Code, train.Capacity)
			}
		case "ut":
			train, error := rdb.updateTrain()
			if error != nil {
				fmt.Println("Error updating train: ", error)
			} else {
				fmt.Println("Train updated:")
				fmt.Println("Id Code Capacity")
				fmt.Println(train.Id, train.Code, train.Capacity)
			}
		case "dt":
			id := readInt("Enter train id: ")
			err := rdb.deleteTrain(int64(id))
			if err != nil {
				fmt.Println("Error deleting train: ", err)
			} else {
				fmt.Println("Train deleted")
			}
		case "cp":
			id := readInt("Enter amount of seats: ")
			before, err := rdb.countRows(context.Background(), "seats")
			if err != nil {
				fmt.Println("Error counting rows: ", err)
			}
			fmt.Println("Amount of rows before call procedure: ", before)
			err = rdb.callProcedure(int64(id))
			if err != nil {
				fmt.Println("Error calling procedure: ", err)
			}
			fmt.Println("Procedure called")
			after, err := rdb.countRows(context.Background(), "seats")
			if err != nil {
				fmt.Println("Error counting rows: ", err)
			}
			fmt.Println("Amount of rows after call procedure: ", after)
		case "snrr":
			level := readInt("Enter isolation level (1: Read Uncommitted, 2: Read Committed, 3: Repeatable Read, 4: Serializable): ")
			if err != nil {
				fmt.Println("Invalid level: ", level)
				continue
			}
			err = rdb.simulateNonRepeatableRead(levelsMap[level])
			if err != nil {
				fmt.Println("Error simulating non-repeatable read: ", err)
			}
		default:
			fmt.Println("Invalid command")
		}
	}
}

func help() {
	fmt.Println("h: help")
	fmt.Println("e: exit")
	fmt.Println("q: quit (same as exit)")
	fmt.Println("cs: create station")
	fmt.Println("gss: get all stations")
	fmt.Println("gs: get station")
	fmt.Println("us: update station")
	fmt.Println("ds: delete station")
	fmt.Println("ct: create train")
	fmt.Println("gts: get all trains")
	fmt.Println("gt: get train")
	fmt.Println("ut: update train")
	fmt.Println("dt: delete train")
	fmt.Println("cp: call procedure")
	fmt.Println("snrr: simulate non-repeatable read")
}

func (rdb *RailwaysDB) simulateNonRepeatableRead(level sql.IsolationLevel) error {
	ctx := context.Background()
	go func() {
		var train Train
		tx, err := rdb.db.BeginTx(ctx, &sql.TxOptions{Isolation: level})
		if err != nil {
			fmt.Println("\nTransaction 1: BeginTx Error:", err)
			return
		}
		defer tx.Rollback()
		// First read
		err = tx.QueryRow("SELECT * FROM trains WHERE train_id = 100").Scan(&train.Id, &train.Code, &train.Capacity)
		if err != nil {
			fmt.Println("\nTransaction 1: QueryRow Error:", err)
			return
		}
		fmt.Printf("\nFirst read:\nid: %d code: %s, capacity: %d\n", train.Id, train.Code, train.Capacity)
		// Wait for a while before the second read
		time.Sleep(2 * time.Second)

		// Second read
		err = tx.QueryRow("SELECT * FROM trains WHERE train_id = 100").Scan(&train.Id, &train.Code, &train.Capacity)
		if err != nil {
			fmt.Println("\nTransaction 1: QueryRow Error:", err)
			return
		}
		fmt.Printf("\nSecond read:\nid: %d code: %s, capacity: %d\n", train.Id, train.Code, train.Capacity)
		if err := tx.Commit(); err != nil {
			fmt.Println("\nTransaction 1: Commit Error:", err)
			return
		}
		fmt.Println("\nTransaction 1: Committed")
	}()

	go func() {
		// Transaction 2: Update the data
		time.Sleep(1 * time.Second) // Wait a bit to ensure the first read happens first

		tx, err := rdb.db.BeginTx(ctx, &sql.TxOptions{Isolation: level})
		if err != nil {
			fmt.Println("\nTransaction 2: BeginTx Error:", err)
			return
		}
		defer tx.Rollback()

		_, err = tx.Exec("UPDATE trains SET capacity = capacity + 1 WHERE train_id = 100")
		if err != nil {
			fmt.Println("\nTransaction 2: Exec Error:", err)
		}

		if err := tx.Commit(); err != nil {
			fmt.Println("\nTransaction 2: Commit Error:", err)
			return
		}
		fmt.Println("\nTransaction 2: Committed")
	}()

	time.Sleep(5 * time.Second) // Wait for the transactions to finish

	return nil
}

func (rdb *RailwaysDB) callProcedure(amount int64) error {
	_, err := rdb.db.Query("CALL GenerateSeats(?)", amount)
	return err
}

// countRows is a utility function to count rows in a table.
func (rdb *RailwaysDB) countRows(ctx context.Context, tableName string) (int, error) {
	var count int
	err := rdb.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM "+tableName).Scan(&count)
	return count, err
}

func readInt(msg string) int {
	fmt.Print(msg)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	id, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Invalid input")
		return -1
	}
	return id
}

func readString(msg string) string {
	fmt.Print(msg)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	return input
}
