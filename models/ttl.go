package models

import (
	"log"
	"os"
	"time"

	"github.com/cmxjs/file-server/dao"
	"gorm.io/gorm"
)

type TTL struct {
	Path     string
	Ttl      int64
	UpdateAt int64
}

var DB *gorm.DB

func InitModel() error {
	if err := dao.SqliteDB.AutoMigrate(&TTL{}); err != nil {
		return err
	}
	DB = dao.SqliteDB
	return nil
}

func UpdateTTL(path string, t int64) error {
	var ttl TTL
	if err := DB.First(&ttl, "path = ?", path).Error; err == nil {
		DB.Model(&ttl).Where("path = ?", path).Update("Ttl", t).Update("UpdateAt", time.Now().Unix())
		return nil
	}

	return DB.Create(&TTL{Path: path, Ttl: t, UpdateAt: time.Now().Unix()}).Error
}

func GetAllTTL() ([]TTL, error) {
	var ttls []TTL
	if err := DB.Find(&ttls).Error; err != nil {
		return nil, err
	}
	return ttls, nil
}

func DelTTL(path string) error {
	var ttl TTL
	return DB.Delete(&ttl, "path = ?", path).Error
}

func DeleteExpiredFile() {
	var ttls []TTL
	if err := DB.Find(&ttls, "Ttl != ?", -1).Error; err != nil {
		log.Println(err)
		return
	}

	for _, ttl := range ttls {
		if (ttl.UpdateAt + ttl.Ttl) < time.Now().Unix() {
			DelTTL(ttl.Path)
			os.Remove(ttl.Path)
			log.Println("Delete expired ttl. ", ttl)
		}
	}
}
