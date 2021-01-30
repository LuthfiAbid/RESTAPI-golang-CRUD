package models

import (
	"errors"
	"html"
	"log"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct{
	ID uint32 `gorm:"primary_key;auto_increment" json:id`
	Username string `gorm:"size:255;not null;unique" json:username`
	Password string `gorm:"size:100;not null" json:password`
	Nama_Lengkap string `gorm:"size:100;not null;" json:nama_lengkap`
	Foto string `gorm:"size:255;null" json:foto`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func Hash(password string) ([]byte, error){
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error{
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u *User) BeforeSave() error{
	hashedPassword, err := Hash(u.Password)
	if err!=nil{return err}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) Prepare(){
	u.ID=0
	u.Username=html.EscapeString(strings.TrimSpace(u.Username))
	u.Nama_Lengkap=html.EscapeString(strings.TrimSpace(u.Nama_Lengkap))
	u.Foto=html.EscapeString(strings.TrimSpace(u.Foto))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *User) Validate(action string) error{
	switch strings.ToLower(action){
		case "update":
			if u.Username==""{return errors.New("Required Username!")}
			if u.Password==""{return errors.New("Required Password!")}
			if u.Foto==""{return errors.New("Required Photo!")}
			return nil
		case "login":
			if u.Username==""{return errors.New("Required Username")}
			if u.Password==""{return errors.New("Required Password!")}
			return nil
		default:
			if u.Username==""{return errors.New("Required Username")}
			if u.Password==""{return errors.New("Required Password!")}
			return nil
	}
}

func (u *User) SaveUser(db *gorm.DB) (*User, error){
	var err error
	err = db.Debug().Create(&u).Error
	if err!=nil{
		return &User{}, err
	}
	return u, nil
}

func (u *User) FindUserByID(db *gorm.DB, uid uint32) (*User, error){
	var err error
	err = db.Debug().Model(User{}).Where("id=?", uid).Take(&u).Error
	if err!=nil{return &User{}, err}
	if gorm.IsRecordNotFoundError(err){
		return &User{}, errors.New("User Not Found!")
	}
	return u, err
}

func (u *User) FindAllUsers(db *gorm.DB) (*[]User, error){
	var err error
	users := []User{}
	err = db.Debug().Model(&User{}).Limit(100).Find(&users).Error
	if err!=nil{return &[]User{} ,err}
	return &users, err
}

func (u *User) UpdateUser(db *gorm.DB, uid uint32) (*User, error){
	err := u.BeforeSave()
	if err!=nil{log.Fatal(err)}
	db=db.Debug().Model(&User{}).Where("id=?", uid).Take(&User{}).UpdateColumn(
		map[string]interface{}{
			"username": u.Username,
			"password": u.Password,
			"update_at": time.Now(),
			"foto": u.Foto,
		},
	)
	if db.Error!=nil{return &User{}, db.Error}
	err = db.Debug().Model(&User{}).Where("id=?", uid).Take(&u).Error
	if err!=nil{return &User{}, err}
	return u, nil
}

func (u *User) DeleteUser(db *gorm.DB, uid uint32)(int64, error){
	db=db.Debug().Model(&User{}).Where("id=?", uid).Take(&User{}).Delete(&User{})
	if db.Error!=nil{return 0, db.Error}
	return db.RowsAffected, nil
}