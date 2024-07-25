//go:generate mockgen -source=storage.go -destination=mock/mock_storage.go -package=mock
package storage

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type IStorage interface {
	CreateUser(entity UserTable) error
	IsExistUserAccount(account string) (exist bool, err error)
	GetUser(id int64) (entity UserTable, err error)
	GetUserByAccount(account string) (entity UserTable, err error)
}

var _ IStorage = (*BlockActionDB)(nil)

type BlockActionDB struct {
	conn *gorm.DB
}

type BlockActionDBCfg struct {
	UserName        string
	Password        string
	Address         string
	DBName          string
	MaxOpenConn     int
	MaxIdleConn     int
	ConnMaxLifetime int
}

type UserTable struct {
	ID        int64  `gorm:"<-:create;column:id;type:bigint;primaryKey;autoIncrement:false;"`                 // ID
	Account   string `gorm:"column:account;type:varchar(45);not null;index:uk_merchant_user_account,unique;"` // 帳號
	Secret    string `gorm:"column:secret;type:varchar(255);not null;"`                                       // Secret
	Name      string `gorm:"column:name;type:varchar(100);not null;"`                                         // 名稱
	Desc      string `gorm:"column:description;type:varchar(200);not null;"`                                  // 描述
	State     int    `gorm:"column:state;type:smallint;not null;"`                                            // 狀態
	CreatedAt int64  `gorm:"<-:create;column:created_at;type:bigint;not null;autoUpdateTime:milli;"`          // 建立時間
	Creator   int64  `gorm:"<-:create;column:creator;type:bigint;not null;"`                                  // 建立者
	UpdatedAt int64  `gorm:"column:updated_at;type:bigint;not null;autoUpdateTime:milli;"`                    // 修改時間
	Updater   int64  `gorm:"column:updater;type:bigint;not null"`                                             // 修改者
}

func (UserTable) TableName() string {
	return "user"
}

func NewBlockActionDB(ctx context.Context, cfg BlockActionDBCfg) (*BlockActionDB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True", cfg.UserName, cfg.Password, cfg.Address, cfg.DBName)
	conn, err := ConnGormMySQL(ctx, dsn, cfg.MaxOpenConn, cfg.MaxIdleConn, cfg.ConnMaxLifetime)
	if err != nil {
		return nil, fmt.Errorf("conn mysql fail : %w", err)
	}
	return &BlockActionDB{
		conn: conn,
	}, nil
}

func (db *BlockActionDB) CreateUser(entity UserTable) error {
	err := db.create(entity.TableName(), &entity)
	if err != nil {
		return err
	}

	return nil
}

func (db *BlockActionDB) create(table string, entity interface{}) error {
	err := db.conn.Table(table).
		Create(entity).Error
	if err != nil {
		return err
	}

	return nil
}

func (db *BlockActionDB) GetUser(id int64) (entity UserTable, err error) {
	err = db.conn.Table(entity.TableName()).
		Where("`id` = ?", id).
		First(&entity).
		Error

	return
}

func (db *BlockActionDB) GetUserByAccount(account string) (entity UserTable, err error) {
	err = db.conn.Table(entity.TableName()).
		Where("`account` = ?", account).
		Find(&entity).
		Error

	return
}

func (db *BlockActionDB) IsExistUserAccount(account string) (exist bool, err error) {
	var cnt int64 = 0
	err = db.conn.Table(UserTable{}.TableName()).
		Where("`account` = ?", account).
		Count(&cnt).
		Error
	exist = cnt > 0
	return
}
