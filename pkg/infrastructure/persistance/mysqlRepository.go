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
	Store       string
	CountryCode string
}

type Result struct {
	Element string
	Total   int
}

type BearerTable struct {
	gorm.Model
	Bearer string
}

type RefreshTokenTable struct {
	gorm.Model
	Token string
}

type CountryTable struct {
	gorm.Model
	Country string
}

func NewStoreTable(storeName string, countryCode string) *StoreTable {
	return &StoreTable{Store: storeName, CountryCode: countryCode}
}

func intialMigration(user string, pwd string, ip string, database string) {
	db, err := gorm.Open(mysql.Open(user+":"+pwd+"@tcp("+ip+")/"+database+"?charset=utf8&parseTime=True&loc=Local"), &gorm.Config{})

	if err != nil {
		zap.L().Panic("Failed to connect to database: ", zap.Error(err))
	}

	db.AutoMigrate(&StoreTable{})
	db.AutoMigrate(&BearerTable{})
	db.AutoMigrate(&RefreshTokenTable{})
	db.AutoMigrate(&CountryTable{})
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

func (db *MysqlRepository) GetRefreshToken() string {
	database := openConnection(db.user, db.pwd, db.ip, db.database)

	var refreshToken RefreshTokenTable
	database.Find(&refreshToken)

	return refreshToken.Token
}

func (db *MysqlRepository) GetCountries() []domain.Country {
	database := openConnection(db.user, db.pwd, db.ip, db.database)

	var countries []CountryTable
	database.Find(&countries)

	return CountryTableToCountryObject(countries)
}

func (db *MysqlRepository) UpdateRefreshToken(newToken string) {
	currentToken := db.GetRefreshToken()
	database := openConnection(db.user, db.pwd, db.ip, db.database)

	var refreshTokenTable RefreshTokenTable

	if currentToken == "" {
		refreshTokenTable.Token = newToken
		database.Create(&refreshTokenTable)
	} else {
		refreshTokenTable.Token = newToken
		database.Model(&RefreshTokenTable{}).Where("token = ?", currentToken).Update("token", newToken)
	}
}

func (db *MysqlRepository) GetStores() []domain.Store {
	database := openConnection(db.user, db.pwd, db.ip, db.database)

	var stores []StoreTable

	today := time.Now().Format("2006-01-02")
	database.Where("created_at > ?", today).Find(&stores)
	return StoreTablesToStoreObjects(stores)
}

func (db *MysqlRepository) GetStoresByTimesAppeared(frequency string, countryCode string) []domain.StoreCounter {
	database := openConnection(db.user, db.pwd, db.ip, db.database)

	var stores []StoreTable
	var result []Result

	database.Model(&stores).Select("store as element, count(store) as total").Where("created_at > ? AND country_code = ?", frequency, countryCode).Group("store").Order("total").Find(&result)

	return StoreTableCountResultsToStoreCounterObjects(result)
}

func (db *MysqlRepository) GetStoresByDayOfWeek(frequency string, countryCode string) []domain.StoreCounter {
	database := openConnection(db.user, db.pwd, db.ip, db.database)

	var stores []StoreTable
	var result []Result
	database.Model(&stores).Select("DAYNAME(CREATED_AT) as element, COUNT(CREATED_AT) as total").Where("created_at > ? AND country_code = ?", frequency, countryCode).Group("DAYNAME(CREATED_AT)").Order("total").Find(&result)

	return StoreTableCountResultsToStoreCounterObjects(result)
}

func (db *MysqlRepository) GetStoresByHourOfDay(frequency string, countryCode string) []domain.StoreCounter {
	database := openConnection(db.user, db.pwd, db.ip, db.database)

	var stores []StoreTable
	var result []Result
	database.Model(&stores).Select("HOUR(CREATED_AT) as element, COUNT(CREATED_AT) as total").Where("created_at > ? AND country_code = ?", frequency, countryCode).Group("HOUR(CREATED_AT)").Order("total").Find(&result)

	return StoreTableCountResultsToStoreCounterObjects(result)
}

func (db *MysqlRepository) GetCountryCodes() []string {
	database := openConnection(db.user, db.pwd, db.ip, db.database)

	var stores []StoreTable
	var countries []string
	database.Model(&stores).Distinct().Pluck("country_code", &countries)

	return countries
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
