package models

import (
	"gorm.io/gorm"
	"time"
)

const POSTS = "posts"

type PostModelImpl struct {
	DBMain *gorm.DB
}

type Post struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string    `gorm:"type:varchar(200);not null" json:"title"`
	Content     string    `gorm:"type:text;not null" json:"content"`
	Category    string    `gorm:"type:varchar(100);not null" json:"category"`
	CreatedDate time.Time `gorm:"autoCreateTime" json:"created_date"`
	UpdatedDate time.Time `gorm:"autoUpdateTime" json:"updated_date"`
	Status      string    `gorm:"type:varchar(100);not null" json:"status"` // Publish | Draft | Thrash
}

func (post Post) TableName() string {
	return POSTS
}

func (m *PostModelImpl) Create(data Post) error {
	return m.DBMain.Create(&data).Error
}

func (m *PostModelImpl) GetAll() ([]Post, error) {
	var data []Post
	m.DBMain.Find(&data)
	return data, nil
}

func (m *PostModelImpl) GetById(id uint) (Post, error) {
	var data Post
	m.DBMain.Where("id = ?", id).Take(&data)
	return data, nil
}

func (m *PostModelImpl) UpdateById(data Post) error {
	return m.DBMain.Table(POSTS).Where("id = ?", data.ID).Updates(&data).Error
}

func (m *PostModelImpl) DeleteById(id uint) error {
	return m.DBMain.Table(POSTS).Where("id = ?", id).Delete(&Post{}).Error
}
