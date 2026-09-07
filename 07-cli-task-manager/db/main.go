package db

import (
	"strconv"
	"time"

	"github.com/AbePlays/gophercises/07-cli-task-manager/utils"
	bolt "go.etcd.io/bbolt"
)

var taskBucket = []byte("tasks")
var db *bolt.DB

type Task struct {
	Key   string
	Value string
}

func Init(dbPath string) error {
	var err error
	db, err = bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(taskBucket)
		return err
	})
}

func CreateTask(task string) (int, error) {
	var id int
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(taskBucket)
		id64, _ := bucket.NextSequence()
		id = int(id64)
		key := utils.Itob(id)
		return bucket.Put(key, []byte(task))
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}

func AllTasks() ([]Task, error) {
	var tasks []Task
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(taskBucket)
		cursor := bucket.Cursor()

		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			tasks = append(tasks, Task{
				Key:   strconv.Itoa(utils.Btoi(key)),
				Value: string(value),
			})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func DeleteTask(id int) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(taskBucket)
		key := utils.Itob(id)
		return bucket.Delete(key)
	})
}
