package service

import (
	"gorm.io/gorm"
	"sharing-vision-id/internal/models"
	"sharing-vision-id/pkg"
	"time"
)

type UserService struct {
	DB *gorm.DB
}

func (s *UserService) CreatePost(input models.Post) (interface{}, error) {

	result, err := pkg.WithTransaction(s.DB, func(tz *gorm.DB) (interface{}, error) {
		impl := models.PostModelImpl{
			DBMain: tz,
		}

		if created := impl.Create(models.Post{
			Title:       input.Title,
			Status:      input.Status,
			Category:    input.Category,
			Content:     input.Content,
			CreatedDate: time.Now(),
			UpdatedDate: time.Now(),
		}); created != nil {
			return nil, created
		}
		return input, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(models.Post), nil
}

func (s *UserService) GetData(id int, page, pageSize int) (interface{}, error) {

	result, err := pkg.WithTransaction(s.DB, func(tz *gorm.DB) (interface{}, error) {
		impl := models.PostModelImpl{
			DBMain: tz,
		}
		if id != 0 {

			resultData, errData := impl.GetById(uint(id))
			if errData != nil {
				return nil, errData
			}
			return resultData, nil
		}

		var count, totalPage int64

		switch {
		case pageSize > 100:
			pageSize = pageSize
		case pageSize <= 0:
			pageSize = 10
		}

		getData, errGet := impl.GetAll(page, pageSize)
		tz.Table(models.POSTS).Count(&count)

		if errGet != nil {
			return nil, errGet
		}

		if count < int64(pageSize) {
			totalPage = 1
		} else {
			totalPage = count / int64(pageSize)
			if (count % int64(pageSize)) != 0 {
				totalPage = totalPage + 1
			}
		}

		return getData, nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *UserService) Update(input models.Post) (interface{}, error) {

	result, err := pkg.WithTransaction(s.DB, func(tz *gorm.DB) (interface{}, error) {
		impl := models.PostModelImpl{
			DBMain: tz,
		}

		if created := impl.UpdateById(models.Post{
			Title:       input.Title,
			Status:      input.Status,
			Category:    input.Category,
			Content:     input.Content,
			CreatedDate: time.Now(),
			UpdatedDate: time.Now(),
			ID:          input.ID,
		}); created != nil {
			return nil, created
		}
		return input, nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *UserService) Delete(id uint) (interface{}, error) {

	result, err := pkg.WithTransaction(s.DB, func(tz *gorm.DB) (interface{}, error) {
		impl := models.PostModelImpl{
			DBMain: tz,
		}

		if created := impl.DeleteById(id); created != nil {
			return nil, created
		}
		return id, nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
