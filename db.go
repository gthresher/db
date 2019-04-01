package main

import (
	"database/sql"
)

func insert(ProjectId int, IssueId string, PongCounter string) (sql.Result, error) {
	return db.Exec("INSERT INTO phonebook VALUES (default, $1, $2, $3)",
		ProjectId, IssueId, PongCounter)
}

func remove(id int) (sql.Result, error) {
	return db.Exec("DELETE FROM phonebook WHERE id=$1", id)
}

func update(ProjectId int, IssueId string, PongCounter string) (sql.Result, error) {
	return db.Exec("UPDATE phonebook SET PongsCounter = $3 WHERE projectid=$1 AND issueid=$2",
		ProjectId, IssueId, PongCounter)
}

func readIssuesById(ProjectId int) (Record, error) {
	var record Record
	row := db.QueryRow("SELECT * FROM phonebook WHERE projectid=$1 ORDER BY projectid", ProjectId)
	return record, row.Scan(&record.ProjectId, &record.IssueId, &record.PongsCounter)
}

func readByName(str string) ([]Record, error) {
	var rows *sql.Rows //обьявляем переменную, в которую положили адрес экземпляра структуры строк бд
	var err error

	if str != "" { //если поданное в функцию значение ненулевое, тогда  отправляем запрос в бд с просьбой выдать все строки, где name похоже на поданную переменную, если нет - то выдаем все строки
		rows, err = db.Query("SELECT * FROM phonebook WHERE name LIKE $1 ORDER BY id", "%"+str+"%")
	} else {
		rows, err = db.Query("SELECT * FROM phonebook ORDER BY id")
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recordsArray = make([]Record, 0) //Обьявляем массив Record емкостью ноль
	var record Record
	for rows.Next() { // здесь обрабатываем машинопонятно респонс базы данных, чтобы выплюнуть []Record
		if err = rows.Scan(&record.ProjectId, &record.IssueId, &record.PongsCounter); err != nil {
			return nil, err
		}
		recordsArray = append(recordsArray, record)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return recordsArray, nil
}
