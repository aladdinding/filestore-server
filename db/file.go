package db

import (
	"database/sql"
	"filestore-server/db/mysql"
	"fmt"
)

func OnFileUploadFinished(filehash string, filename string, filesize int64, fileaddr string) bool {
	stmt, err := mysql.DBConn().Prepare(
		"insert ignore into tbl_file (`file_sha1`,`file_name`,`file_size`,`file_addr`,`status`) values (?,?,?,?,1)")
	if err != nil {
		fmt.Println("Failed to prepare statement, err: " + err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			fmt.Printf("File with hash: %s has been uploaded before", filehash)
		}
		return true
	}
	return false
}

type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

func GetFileMeta(filehash string) (*TableFile, error) {
	stmt, err := mysql.DBConn().Prepare(
		"select file_sha1,file_addr,file_name,file_size from tbl_file where file_sha1=? and status=1 limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer stmt.Close()
	tfile := TableFile{}
	err = stmt.QueryRow(filehash).Scan(&tfile.FileHash, &tfile.FileName, &tfile.FileSize)
	if err != nil {
		if err == sql.ErrNoRows {
			// 查不到对应记录，返回参数及错误均为nil
			return nil, nil
		} else {
			fmt.Println(err.Error())
			return nil, err
		}
	}
	return &tfile, nil
}

// GetFileMetaList:从mysql批量获取文件元信息
func GetFileMetaList(limit int) ([]TableFile, error) {
	stmt, err := mysql.DBConn().Prepare(
		"select file_sha1,file_addr,file_name,file_size from tbl_file where status=1 limit ?")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(limit)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	cloumns, _ := rows.Columns()
	values := make([]sql.RawBytes, len(cloumns))
	var tfiles []TableFile
	for i := 0; i < len(values) && rows.Next(); i++ {
		tfile := TableFile{}
		err = rows.Scan(&tfile.FileHash, &tfile.FileAddr, &tfile.FileName, &tfile.FileSize)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		tfiles = append(tfiles, tfile)
	}
	fmt.Println(len(tfiles))
	return tfiles, nil
}

//UpdateFileLocation:更新文件的存储地址（如果文件被转移了）
func UpdateFileLocation(filehash string, fileaddr string) bool {
	stmt, err := mysql.DBConn().Prepare(
		"update tbl_file set `file_addr`=? where `file_sha1`=? limit 1")
	if err != nil {
		fmt.Println("预编译sql失败，err: " + err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(fileaddr, filehash)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			fmt.Println("更新文件location失败，filehash: %s", filehash)
		}
		return true
	}
	return false
}
