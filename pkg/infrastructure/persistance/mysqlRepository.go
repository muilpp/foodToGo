package persistance

import (
	"time"

	"github.com/marc/get-food-to-go/pkg/domain"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type StoreTable struct {
	gorm.Model
	Store string
}

type BearerTable struct {
	gorm.Model
	Bearer string
}

func NewStoreTable(storeName string) *StoreTable {
	return &StoreTable{Store: storeName}
}

func intialMigration(user string, pwd string, ip string, database string) {
	db, err := gorm.Open(mysql.Open(user+":"+pwd+"@tcp("+ip+")/"+database+"?charset=utf8&parseTime=True&loc=Local"), &gorm.Config{})

	if err != nil {
		zap.L().Panic("Failed to connect to database: ", zap.Error(err))
	}

	db.AutoMigrate(&StoreTable{})
	db.AutoMigrate(&BearerTable{})
}

type MysqlRepository struct {
	user     string
	pwd      string
	ip       string
	database string
}

func NewMysqlRepository(user string, pwd string, ip string, database string) *MysqlRepository {
	intialMigration(user, pwd, ip, database)
	return &MysqlRepository{user, pwd, ip, database}
}

func (db *MysqlRepository) GetBearer() string {
	database := openConnection(db.user, db.pwd, db.ip, db.database)

	var bearer BearerTable
	database.Find(&bearer)

	return bearer.Bearer
}

func (db *MysqlRepository) UpdateBearer(newBearer string) {
	currentBearer := db.GetBearer()
	database := openConnection(db.user, db.pwd, db.ip, db.database)

	var bearerTable BearerTable

	if currentBearer == "" {
		bearerTable.Bearer = newBearer
		database.Create(&bearerTable)
	} else {
		bearerTable.Bearer = newBearer
		database.Model(&BearerTable{}).Where("bearer = ?", currentBearer).Update("bearer", newBearer)
	}
}

func (db *MysqlRepository) GetStores() []domain.Store {
	database := openConnection(db.user, db.pwd, db.ip, db.database)

	var stores []StoreTable
	database.Find(&stores)

	today := time.Now().Format("2006-01-02")
	database.Where("created_at > ?", today).Find(&stores)
	return StoreTablesToStoreObjects(stores)
}

func (db *MysqlRepository) AddStores(stores []domain.Store) {
	database := openConnection(db.user, db.pwd, db.ip, db.database)
	database.CreateInBatches(StoreObjectsToStoreTables(stores), 10)
}

func openConnection(user string, pwd string, ip string, database string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(user+":"+pwd+"@tcp("+ip+")/"+database+"?charset=utf8&parseTime=True&loc=Local"), &gorm.Config{})

	if err != nil {
		zap.L().Panic("Failed to connect to database: ", zap.Error(err))
	}

	return db
}
