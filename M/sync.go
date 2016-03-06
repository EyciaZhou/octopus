package MOctopus

import (
	"database/sql"
	"errors"
	"time"
	"fmt"
)

func cronPullVersionInfos() {
	time.Sleep(10 * time.Second)

	for ;; {
		err := pullVersionInfos()
		if err != nil {
			//TODO: LOG ERROR.... not decided which log module to use yet
			fmt.Printf("error: " + err.Error())
		}

		time.Sleep(30*time.Second)
	}
}

func init() {
	go cronPullVersionInfos()
}

// id judgment yes no next versionname type

type node_include_two struct {
	Id string
	Judgment sql.NullString
	Yes sql.NullString
	No sql.NullString
	Next sql.NullString
	VersionName sql.NullString
	Type string
}

func pushVersionInfos() (_err error) {
	defer func() {
		err := recover()
		if err != nil {
			_err = err.(error)
		}
	}()

	versionsMutex.RLock()
	defer versionsMutex.RUnlock()

	tx, err := db.Begin()
	defer tx.Rollback()
	if err != nil {
		return err
	}
//id, judgment, yes, no, next, versionName, type, lstModify
	stmtFork, err := tx.Prepare(`
		INSERT INTO
				versions (id, judgment, yes, no, type)
			VALUES
				(?,?,?,?,"fork")
			ON DUPLICATE KEY UPDATE
				id=VALUES(id),
				judgment=VALUES(judgment),
				yes=VALUES(yes),
				no=VALUES(no),
				type="fork",
				next=NULL,
				versionName=NULL
	`)
	if err != nil {
		return err
	}

	stmtVersion, err := tx.Prepare(`
		INSERT INTO
				versions (id, next, versionName, type)
			VALUES
				(?,?,?,"version")
			ON DUPLICATE KEY UPDATE
				id=VALUES(id),
				judgment=NULL,
				yes=NULL,
				no=NULL,
				type="version",
				next=VALUES(next),
				versionName=VALUES(versionName)`)
	if err != nil {
		return err
	}

	for _, v := range versions {
		switch v.Type() {
		case "fork":
			f := v.(*fork)
			_, err := stmtFork.Exec(f.id, f.judgment, f.yes, f.no)
			if err != nil {
				return err
			}
		case "version":
			ver := v.(*version)
			_, err := stmtVersion.Exec(ver.id, ver.next, ver.VersionName)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func pullVersionInfos() (_err error) {
	defer func() {
		err := recover()
		if err != nil {
			_err = err.(error)
		}
	}()

	row := db.QueryRow(`
		SELECT
				max(lstModify)
			FROM versions
	`)
	var lm int64
	row.Scan(&lm)

	if (lm == lstModify) {
		return
	}

	rows, err := db.Query(`
		SELECT
				id, judgment, yes, no, next, versionName, type
			FROM versions
	`)
	defer rows.Close()

	if err != nil {
		return err
	}

	nodes_tmp := []*node_include_two{}

	for rows.Next() {
		n := &node_include_two{}
		_err = rows.Scan(&n.Id, &n.Judgment, &n.Yes, &n.No, &n.VersionName, &n.Type)

		if _err != nil {
			return _err
		}

		nodes_tmp = append(nodes_tmp, n)
	}

	result := map[string]node{}

	for _, n := range nodes_tmp {
		switch n.Type {
		case "fork":
			if n.Judgment.Valid && n.Yes.Valid && n.No.Valid {
				result[n.Id] = &fork {
					n.Id,
					n.Judgment.String,
					n.Yes.String,
					n.No.String,
				}
			} else {
				return errors.New("invalid fork, id:" + n.Id)
			}
		case "version":
			if n.Next.Valid && n.VersionName.Valid {
				result[n.Id] = &version {
					n.Id,
					n.Next.String,
					n.VersionName.String,
				}
			} else {
				return errors.New("invalid version, id:" + n.Id)
			}
		default:
			return errors.New("no such node type: " + n.Type)
		}
	}

	tmpVersionWithName := map[string]*version{}

	for _, v := range result {
		if v.Type() == "version" {
			ver := v.(*version)
			tmpVersionWithName[ver.VersionName] = ver
		}
	}

	versionsMutex.Lock()
	defer versionsMutex.Unlock()
	versions = result
	versionsWithName = tmpVersionWithName
	return nil
}