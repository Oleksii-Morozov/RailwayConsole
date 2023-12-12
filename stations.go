package main

import (
	"fmt"
)

func (rdb *RailwaysDB) createStation() (*Station, error) {
	var station StationCreate
	station.City = readString("city: ")
	station.Name = readString("name:")

	stmt, err := rdb.db.Prepare("INSERT INTO stations(city, name) VALUES(?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(station.City, station.Name)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return rdb.getStationById(int64(id))
}

func (rdb *RailwaysDB) getAllStations() (*[]Station, error) {
	var stations []Station
	rows, err := rdb.db.Query(`SELECT * FROM stations 
														ORDER BY station_id ASC
														LIMIT 100`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var s Station

		if err := rows.Scan(&s.Id, &s.City, &s.Name); err != nil {
			return nil, err
		}

		stations = append(stations, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &stations, nil
}

func (rdb *RailwaysDB) getStationById(id int64) (*Station, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid id: %d", id)
	}

	s := &Station{}
	query := `SELECT * FROM stations WHERE station_id = ?`
	err := rdb.db.QueryRow(query, id).Scan(&s.Id, &s.City, &s.Name)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (rdb *RailwaysDB) updateStation() (*Station, error) {
	var station Station

	station.Id = int64(readInt("id: "))
	station.City = readString("city: ")
	station.Name = readString("name: ")

	query := `UPDATE stations SET city = ?, name = ? WHERE station_id = ?`
	res, err := rdb.db.Exec(query, station.City, station.Name, station.Id)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		fmt.Println("No rows affected")
	}

	return rdb.getStationById(station.Id)
}

func (rdb *RailwaysDB) deleteStation(id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid id: %d", id)
	}

	query := `DELETE FROM stations WHERE station_id = ?`
	res, err := rdb.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		fmt.Println("No rows affected")
	}

	fmt.Println("Deleted station with id", id)
	return nil
}
