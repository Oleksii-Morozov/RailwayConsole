package main

func (rdb *RailwaysDB) createTrain() (*Train, error) {
	var train TrainCreate

	train.Code = readString("code: ")
	train.Capacity = readInt("capacity: ")

	stmt, err := rdb.db.Prepare("INSERT INTO trains(code, capacity) VALUES(?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(train.Code, train.Capacity)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return rdb.getTrainById(int64(id))
}

func (rdb *RailwaysDB) getAllTrains() (*[]Train, error) {
	var trains []Train
	rows, err := rdb.db.Query(`SELECT * FROM trains 
															ORDER BY train_id ASC
															LIMIT 100`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var t Train

		if err := rows.Scan(&t.Id, &t.Code, &t.Capacity); err != nil {
			return nil, err
		}

		trains = append(trains, t)
	}

	return &trains, nil
}

func (rdb *RailwaysDB) getTrainById(id int64) (*Train, error) {
	var t Train
	err := rdb.db.QueryRow(`SELECT * FROM trains 
													WHERE train_id = ?`, id).Scan(&t.Id, &t.Code, &t.Capacity)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (rdb *RailwaysDB) updateTrain() (*Train, error) {
	var train Train

	train.Id = int64(readInt("id"))
	train.Code = readString("code")
	train.Capacity = readInt("capacity")

	stmt, err := rdb.db.Prepare("UPDATE trains SET code = ?, capacity = ? WHERE train_id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(train.Code, train.Capacity, train.Id)
	if err != nil {
		return nil, err
	}

	return rdb.getTrainById(int64(train.Id))
}

func (rdb *RailwaysDB) deleteTrain(id int64) error {
	stmt, err := rdb.db.Prepare("DELETE FROM trains WHERE train_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}
