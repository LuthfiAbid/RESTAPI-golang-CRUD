package seed

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/LuthfiAbid/golang_interview/api/models"
)

var users = []models.User{
	models.User{
		Username: "LuthfiAbid",
		Password: "password123",
		Nama_Lengkap : "Muhammad Luthfi Abid Cahyadi",
	},
	models.User{
		Username: "AbidLuthfi",
		Password: "password321",
		Nama_Lengkap : "Muhammad Luthfi Abid Cahyadi",
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	// err = db.Debug().Model().AddForeignKey("author_id", "users(id)", "cascade", "cascade").Error
	// if err != nil {
	// 	log.Fatalf("attaching foreign key error: %v", err)
	// }

	for i, _ := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
	}
}