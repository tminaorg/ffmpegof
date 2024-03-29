package processor

import (
	"errors"
	"fmt"

	"database/sql"
	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

func sqlInsertProcess(dbType string) (string, error) {
	switch dbType {
	case "sqlite":
		return `INSERT INTO processes (host_id, process_id, cmd) VALUES (?, ?, ?) `, nil
	case "postgres":
		return `INSERT INTO processes (host_id, process_id, cmd) VALUES ($1, $2, $3) `, nil
	default:
		return "", fmt.Errorf("incorrect database type")
	}
}

func (store *datastore) InsertProcess(process Process) error {
	sqlInsertProcess, err := sqlInsertProcess(store.dbType)
	if err != nil {
		return err
	}

	tx, err := store.Begin()
	if err != nil {
		return err
	}

	if _, err = tx.Exec(sqlInsertProcess, process.HostId, process.ProcessId, process.Cmd); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			panic(rollbackErr)
		}
		return err
	}

	return tx.Commit()
}

func sqlDeleteProcesses(dbType string) (string, error) {
	switch dbType {
	case "sqlite":
		return `DELETE FROM processes`, nil
	case "postgres":
		return `DELETE FROM processes`, nil
	default:
		return "", fmt.Errorf("incorrect database type")
	}
}

func (store *datastore) DeleteProcesses() error {
	sqlDeleteProcess, err := sqlDeleteProcesses(store.dbType)
	if err != nil {
		return err
	}

	_, err = store.Exec(sqlDeleteProcess)
	if err != nil {
		return fmt.Errorf("delete processes: %w", err)
	}

	return nil
}

func sqlDeleteProcessesWhere(dbType string) (string, error) {
	switch dbType {
	case "sqlite":
		return `DELETE FROM processes WHERE %s=?`, nil
	case "postgres":
		return `DELETE FROM processes WHERE %s=$1`, nil
	default:
		return "", fmt.Errorf("incorrect database type")
	}
}

func (store *datastore) DeleteProcessesWhere(field string, process Process) error {
	sqlDeleteProcessesWhere, err := sqlDeleteProcessesWhere(store.dbType)
	if err != nil {
		return err
	}
	sqlDeleteProcessesWhere = fmt.Sprintf(sqlDeleteProcessesWhere, field)
	
	if field == "id" {
		_, err = store.Exec(sqlDeleteProcessesWhere, process.Id)
	} else if field == "host_id" {
		_, err = store.Exec(sqlDeleteProcessesWhere, process.HostId)
	} else if field == "process_id" {
		_, err = store.Exec(sqlDeleteProcessesWhere, process.ProcessId)
	} else {
		return fmt.Errorf("delete processes where: wrong field name")
	}
	
	if err != nil {
		return fmt.Errorf("delete processes where: %w", err)
	}

	return nil
}

func sqlSelectCountProcesses(dbType string) (string, error) {
	switch dbType {
	case "sqlite":
		return `SELECT COUNT(id) FROM processes`, nil
	case "postgres":
		return `SELECT COUNT(id) FROM processes`, nil
	default:
		return "", fmt.Errorf("incorrect database type")
	}
}

func (store *datastore) SelectCountProcesses() (int, error) {
	sqlSelectCountProcesses, err := sqlSelectCountProcesses(store.dbType)
	if err != nil {
		return 0, err
	}

	row := store.QueryRow(sqlSelectCountProcesses)

	remaining := 0
	err = row.Scan(&remaining)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return remaining, nil
	case err != nil:
		return remaining, fmt.Errorf("select count processes: %w", err)
	}

	return remaining, nil
}

func sqlSelectCountProcessesWhere(dbType string) (string, error) {
	switch dbType {
	case "sqlite":
		return `SELECT COUNT(id) FROM processes WHERE host_id=?`, nil
	case "postgres":
		return `SELECT COUNT(id) FROM processes WHERE host_id=$1`, nil
	default:
		return "", fmt.Errorf("incorrect database type")
	}
}

func (store *datastore) SelectCountProcessesWhere(host Host) (int, error) {
	sqlSelectCountProcessesWhere, err := sqlSelectCountProcessesWhere(store.dbType)
	if err != nil {
		return 0, err
	}

	row := store.QueryRow(sqlSelectCountProcessesWhere, host.Id)

	remaining := 0
	err = row.Scan(&remaining)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return remaining, nil
	case err != nil:
		return remaining, fmt.Errorf("select count processes where: %w", err)
	}

	return remaining, nil
}

func sqlSelectProcesses(dbType string) (string, error) {
	switch dbType {
	case "sqlite":
		return `SELECT %s FROM processes`, nil
	case "postgres":
		return `SELECT %s FROM processes`, nil
	default:
		return "", fmt.Errorf("incorrect database type")
	}
}

func (store *datastore) SelectProcesses() (processes []Process, err error) {
	sqlSelectProcesses, err := sqlSelectProcesses(store.dbType)
	if err != nil {
		return processes, err
	}
	sqlSelectProcesses = fmt.Sprintf(sqlSelectProcesses, "*")

	rows, err := store.Query(sqlSelectProcesses)
	if err != nil {
		return processes, err
	}

	defer rows.Close()
	for rows.Next() {
		process := Process{}
		err = rows.Scan(&process.Id, &process.HostId, &process.ProcessId, &process.Cmd)
		if err != nil {
			return processes, err
		}

		processes = append(processes, process)
	}

	return processes, rows.Err()
}

func (store *datastore) SelectProcessesId() (processes []Process, err error) {
	sqlSelectProcesses, err := sqlSelectProcesses(store.dbType)
	if err != nil {
		return processes, err
	}
	sqlSelectProcesses = fmt.Sprintf(sqlSelectProcesses, "id")

	rows, err := store.Query(sqlSelectProcesses)
	if err != nil {
		return processes, err
	}

	defer rows.Close()
	for rows.Next() {
		process := Process{}
		err = rows.Scan(&process.Id)
		if err != nil {
			return processes, err
		}

		processes = append(processes, process)
	}

	return processes, rows.Err()
}

func sqlSelectProcessesWhere(dbType string) (string, error) {
	switch dbType {
	case "sqlite":
		return `SELECT %s FROM processes WHERE host_id=? ORDER BY id DESC`, nil
	case "postgres":
		return `SELECT %s FROM processes WHERE host_id=$1 ORDER BY id DESC`, nil
	default:
		return "", fmt.Errorf("incorrect database type")
	}
}

func (store *datastore) SelectProcessesWhere(host Host) (processes []Process, err error) {
	sqlSelectProcessesWhere, err := sqlSelectProcessesWhere(store.dbType)
	if err != nil {
		return processes, err
	}
	sqlSelectProcessesWhere = fmt.Sprintf(sqlSelectProcessesWhere, "*")

	rows, err := store.Query(sqlSelectProcessesWhere, host.Id)
	if err != nil {
		return processes, err
	}

	defer rows.Close()
	for rows.Next() {
		process := Process{}
		err = rows.Scan(&process.Id, &process.HostId, &process.ProcessId, &process.Cmd)
		if err != nil {
			return processes, err
		}

		processes = append(processes, process)
	}

	return processes, rows.Err()
}

func (store *datastore) SelectProcessesIdWhere(host Host) (processes []Process, err error) {
	sqlSelectProcessesWhere, err := sqlSelectProcessesWhere(store.dbType)
	if err != nil {
		return processes, err
	}
	sqlSelectProcessesWhere = fmt.Sprintf(sqlSelectProcessesWhere, "id")

	rows, err := store.Query(sqlSelectProcessesWhere, host.Id)
	if err != nil {
		return processes, err
	}

	defer rows.Close()
	for rows.Next() {
		process := Process{}
		err = rows.Scan(&process.Id)
		if err != nil {
			return processes, err
		}

		processes = append(processes, process)
	}

	return processes, rows.Err()
}
