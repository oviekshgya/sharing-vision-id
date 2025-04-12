package pkg

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"reflect"
	"strings"
	"time"
)

func WithTransaction(db *gorm.DB, fn func(tz *gorm.DB) (interface{}, error)) (interface{}, error) {

	tx := db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			panic(r)
		} else if tx.Error != nil {
			_ = tx.Rollback()
		} else {
			cerr := tx.Commit().Error
			if cerr != nil {
				tx.Error = fmt.Errorf("error committing transaction: %v", cerr)
			}
		}
	}()
	res, err := fn(tx)
	if err != nil {
		tx.Error = err
		return nil, err
	}

	return res, nil
}

func UpdateFieldsDynamic(input interface{}) map[string]interface{} {
	v := reflect.ValueOf(input)
	ind := reflect.Indirect(v)
	updateData := make(map[string]interface{})
	for i := 0; i < ind.NumField(); i++ {
		value := v.Field(i)
		fieldType := ind.Type().Field(i)
		jsonTag := fieldType.Tag.Get("gorm")
		stringSplt := strings.Split(jsonTag, ":")
		var stringsplit, source string
		if jsonTag != "" {
			stringsplit = stringSplt[1]
			source = stringSplt[1]
		}
		if strings.Contains(jsonTag, "primaryKey") {
			continue
		}
		if strings.Contains(stringsplit, ";") {
			sources := strings.SplitAfter(stringsplit, ";")
			source = strings.Replace(sources[0], ";", "", -1)
		}
		fmt.Println("source", source)
		if len(stringSplt) >= 2 {
			if !value.IsZero() && stringsplit != "" {
				valuetype := value.Type().String()
				fmt.Println("valuetype", valuetype)
				switch v := value.Interface().(type) {
				case *int:
					updateData[source] = *v
				case *string:
					updateData[source] = *v
				case *float32:
					updateData[source] = *v
				case *float64:
					updateData[source] = *v
				case *time.Time:
					updateData[source] = *v
				default:
					updateData[source] = v
				}
			}
		}
	}
	return updateData
}

func WithTransactionMongo(client *mongo.Client, fn func(sessCtx mongo.SessionContext) (interface{}, error)) (interface{}, error) {
	session, err := client.StartSession()
	if err != nil {
		return nil, fmt.Errorf("error starting session: %v", err)
	}
	defer session.EndSession(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var res interface{}
	err = mongo.WithSession(ctx, session, func(sessCtx mongo.SessionContext) error {
		err := sessCtx.StartTransaction()
		if err != nil {
			return fmt.Errorf("error starting transaction: %v", err)
		}

		res, err = fn(sessCtx)
		if err != nil {
			_ = session.AbortTransaction(sessCtx)
			return err
		}

		err = session.CommitTransaction(sessCtx)
		if err != nil {
			return fmt.Errorf("error committing transaction: %v", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}
