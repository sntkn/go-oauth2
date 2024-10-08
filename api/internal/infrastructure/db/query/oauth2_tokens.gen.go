// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"github.com/sntkn/go-oauth2/api/internal/infrastructure/db/model"
)

func newOauth2Token(db *gorm.DB, opts ...gen.DOOption) oauth2Token {
	_oauth2Token := oauth2Token{}

	_oauth2Token.oauth2TokenDo.UseDB(db, opts...)
	_oauth2Token.oauth2TokenDo.UseModel(&model.Oauth2Token{})

	tableName := _oauth2Token.oauth2TokenDo.TableName()
	_oauth2Token.ALL = field.NewAsterisk(tableName)
	_oauth2Token.AccessToken = field.NewString(tableName, "access_token")
	_oauth2Token.ClientID = field.NewString(tableName, "client_id")
	_oauth2Token.UserID = field.NewString(tableName, "user_id")
	_oauth2Token.Scope = field.NewString(tableName, "scope")
	_oauth2Token.ExpiresAt = field.NewTime(tableName, "expires_at")
	_oauth2Token.RevokedAt = field.NewTime(tableName, "revoked_at")
	_oauth2Token.CreatedAt = field.NewTime(tableName, "created_at")
	_oauth2Token.UpdatedAt = field.NewTime(tableName, "updated_at")

	_oauth2Token.fillFieldMap()

	return _oauth2Token
}

type oauth2Token struct {
	oauth2TokenDo

	ALL         field.Asterisk
	AccessToken field.String
	ClientID    field.String
	UserID      field.String
	Scope       field.String
	ExpiresAt   field.Time
	RevokedAt   field.Time
	CreatedAt   field.Time
	UpdatedAt   field.Time

	fieldMap map[string]field.Expr
}

func (o oauth2Token) Table(newTableName string) *oauth2Token {
	o.oauth2TokenDo.UseTable(newTableName)
	return o.updateTableName(newTableName)
}

func (o oauth2Token) As(alias string) *oauth2Token {
	o.oauth2TokenDo.DO = *(o.oauth2TokenDo.As(alias).(*gen.DO))
	return o.updateTableName(alias)
}

func (o *oauth2Token) updateTableName(table string) *oauth2Token {
	o.ALL = field.NewAsterisk(table)
	o.AccessToken = field.NewString(table, "access_token")
	o.ClientID = field.NewString(table, "client_id")
	o.UserID = field.NewString(table, "user_id")
	o.Scope = field.NewString(table, "scope")
	o.ExpiresAt = field.NewTime(table, "expires_at")
	o.RevokedAt = field.NewTime(table, "revoked_at")
	o.CreatedAt = field.NewTime(table, "created_at")
	o.UpdatedAt = field.NewTime(table, "updated_at")

	o.fillFieldMap()

	return o
}

func (o *oauth2Token) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := o.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (o *oauth2Token) fillFieldMap() {
	o.fieldMap = make(map[string]field.Expr, 8)
	o.fieldMap["access_token"] = o.AccessToken
	o.fieldMap["client_id"] = o.ClientID
	o.fieldMap["user_id"] = o.UserID
	o.fieldMap["scope"] = o.Scope
	o.fieldMap["expires_at"] = o.ExpiresAt
	o.fieldMap["revoked_at"] = o.RevokedAt
	o.fieldMap["created_at"] = o.CreatedAt
	o.fieldMap["updated_at"] = o.UpdatedAt
}

func (o oauth2Token) clone(db *gorm.DB) oauth2Token {
	o.oauth2TokenDo.ReplaceConnPool(db.Statement.ConnPool)
	return o
}

func (o oauth2Token) replaceDB(db *gorm.DB) oauth2Token {
	o.oauth2TokenDo.ReplaceDB(db)
	return o
}

type oauth2TokenDo struct{ gen.DO }

type IOauth2TokenDo interface {
	gen.SubQuery
	Debug() IOauth2TokenDo
	WithContext(ctx context.Context) IOauth2TokenDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IOauth2TokenDo
	WriteDB() IOauth2TokenDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IOauth2TokenDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IOauth2TokenDo
	Not(conds ...gen.Condition) IOauth2TokenDo
	Or(conds ...gen.Condition) IOauth2TokenDo
	Select(conds ...field.Expr) IOauth2TokenDo
	Where(conds ...gen.Condition) IOauth2TokenDo
	Order(conds ...field.Expr) IOauth2TokenDo
	Distinct(cols ...field.Expr) IOauth2TokenDo
	Omit(cols ...field.Expr) IOauth2TokenDo
	Join(table schema.Tabler, on ...field.Expr) IOauth2TokenDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IOauth2TokenDo
	RightJoin(table schema.Tabler, on ...field.Expr) IOauth2TokenDo
	Group(cols ...field.Expr) IOauth2TokenDo
	Having(conds ...gen.Condition) IOauth2TokenDo
	Limit(limit int) IOauth2TokenDo
	Offset(offset int) IOauth2TokenDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IOauth2TokenDo
	Unscoped() IOauth2TokenDo
	Create(values ...*model.Oauth2Token) error
	CreateInBatches(values []*model.Oauth2Token, batchSize int) error
	Save(values ...*model.Oauth2Token) error
	First() (*model.Oauth2Token, error)
	Take() (*model.Oauth2Token, error)
	Last() (*model.Oauth2Token, error)
	Find() ([]*model.Oauth2Token, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Oauth2Token, err error)
	FindInBatches(result *[]*model.Oauth2Token, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.Oauth2Token) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IOauth2TokenDo
	Assign(attrs ...field.AssignExpr) IOauth2TokenDo
	Joins(fields ...field.RelationField) IOauth2TokenDo
	Preload(fields ...field.RelationField) IOauth2TokenDo
	FirstOrInit() (*model.Oauth2Token, error)
	FirstOrCreate() (*model.Oauth2Token, error)
	FindByPage(offset int, limit int) (result []*model.Oauth2Token, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IOauth2TokenDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (o oauth2TokenDo) Debug() IOauth2TokenDo {
	return o.withDO(o.DO.Debug())
}

func (o oauth2TokenDo) WithContext(ctx context.Context) IOauth2TokenDo {
	return o.withDO(o.DO.WithContext(ctx))
}

func (o oauth2TokenDo) ReadDB() IOauth2TokenDo {
	return o.Clauses(dbresolver.Read)
}

func (o oauth2TokenDo) WriteDB() IOauth2TokenDo {
	return o.Clauses(dbresolver.Write)
}

func (o oauth2TokenDo) Session(config *gorm.Session) IOauth2TokenDo {
	return o.withDO(o.DO.Session(config))
}

func (o oauth2TokenDo) Clauses(conds ...clause.Expression) IOauth2TokenDo {
	return o.withDO(o.DO.Clauses(conds...))
}

func (o oauth2TokenDo) Returning(value interface{}, columns ...string) IOauth2TokenDo {
	return o.withDO(o.DO.Returning(value, columns...))
}

func (o oauth2TokenDo) Not(conds ...gen.Condition) IOauth2TokenDo {
	return o.withDO(o.DO.Not(conds...))
}

func (o oauth2TokenDo) Or(conds ...gen.Condition) IOauth2TokenDo {
	return o.withDO(o.DO.Or(conds...))
}

func (o oauth2TokenDo) Select(conds ...field.Expr) IOauth2TokenDo {
	return o.withDO(o.DO.Select(conds...))
}

func (o oauth2TokenDo) Where(conds ...gen.Condition) IOauth2TokenDo {
	return o.withDO(o.DO.Where(conds...))
}

func (o oauth2TokenDo) Order(conds ...field.Expr) IOauth2TokenDo {
	return o.withDO(o.DO.Order(conds...))
}

func (o oauth2TokenDo) Distinct(cols ...field.Expr) IOauth2TokenDo {
	return o.withDO(o.DO.Distinct(cols...))
}

func (o oauth2TokenDo) Omit(cols ...field.Expr) IOauth2TokenDo {
	return o.withDO(o.DO.Omit(cols...))
}

func (o oauth2TokenDo) Join(table schema.Tabler, on ...field.Expr) IOauth2TokenDo {
	return o.withDO(o.DO.Join(table, on...))
}

func (o oauth2TokenDo) LeftJoin(table schema.Tabler, on ...field.Expr) IOauth2TokenDo {
	return o.withDO(o.DO.LeftJoin(table, on...))
}

func (o oauth2TokenDo) RightJoin(table schema.Tabler, on ...field.Expr) IOauth2TokenDo {
	return o.withDO(o.DO.RightJoin(table, on...))
}

func (o oauth2TokenDo) Group(cols ...field.Expr) IOauth2TokenDo {
	return o.withDO(o.DO.Group(cols...))
}

func (o oauth2TokenDo) Having(conds ...gen.Condition) IOauth2TokenDo {
	return o.withDO(o.DO.Having(conds...))
}

func (o oauth2TokenDo) Limit(limit int) IOauth2TokenDo {
	return o.withDO(o.DO.Limit(limit))
}

func (o oauth2TokenDo) Offset(offset int) IOauth2TokenDo {
	return o.withDO(o.DO.Offset(offset))
}

func (o oauth2TokenDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IOauth2TokenDo {
	return o.withDO(o.DO.Scopes(funcs...))
}

func (o oauth2TokenDo) Unscoped() IOauth2TokenDo {
	return o.withDO(o.DO.Unscoped())
}

func (o oauth2TokenDo) Create(values ...*model.Oauth2Token) error {
	if len(values) == 0 {
		return nil
	}
	return o.DO.Create(values)
}

func (o oauth2TokenDo) CreateInBatches(values []*model.Oauth2Token, batchSize int) error {
	return o.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (o oauth2TokenDo) Save(values ...*model.Oauth2Token) error {
	if len(values) == 0 {
		return nil
	}
	return o.DO.Save(values)
}

func (o oauth2TokenDo) First() (*model.Oauth2Token, error) {
	if result, err := o.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Oauth2Token), nil
	}
}

func (o oauth2TokenDo) Take() (*model.Oauth2Token, error) {
	if result, err := o.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Oauth2Token), nil
	}
}

func (o oauth2TokenDo) Last() (*model.Oauth2Token, error) {
	if result, err := o.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Oauth2Token), nil
	}
}

func (o oauth2TokenDo) Find() ([]*model.Oauth2Token, error) {
	result, err := o.DO.Find()
	return result.([]*model.Oauth2Token), err
}

func (o oauth2TokenDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Oauth2Token, err error) {
	buf := make([]*model.Oauth2Token, 0, batchSize)
	err = o.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (o oauth2TokenDo) FindInBatches(result *[]*model.Oauth2Token, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return o.DO.FindInBatches(result, batchSize, fc)
}

func (o oauth2TokenDo) Attrs(attrs ...field.AssignExpr) IOauth2TokenDo {
	return o.withDO(o.DO.Attrs(attrs...))
}

func (o oauth2TokenDo) Assign(attrs ...field.AssignExpr) IOauth2TokenDo {
	return o.withDO(o.DO.Assign(attrs...))
}

func (o oauth2TokenDo) Joins(fields ...field.RelationField) IOauth2TokenDo {
	for _, _f := range fields {
		o = *o.withDO(o.DO.Joins(_f))
	}
	return &o
}

func (o oauth2TokenDo) Preload(fields ...field.RelationField) IOauth2TokenDo {
	for _, _f := range fields {
		o = *o.withDO(o.DO.Preload(_f))
	}
	return &o
}

func (o oauth2TokenDo) FirstOrInit() (*model.Oauth2Token, error) {
	if result, err := o.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Oauth2Token), nil
	}
}

func (o oauth2TokenDo) FirstOrCreate() (*model.Oauth2Token, error) {
	if result, err := o.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Oauth2Token), nil
	}
}

func (o oauth2TokenDo) FindByPage(offset int, limit int) (result []*model.Oauth2Token, count int64, err error) {
	result, err = o.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = o.Offset(-1).Limit(-1).Count()
	return
}

func (o oauth2TokenDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = o.Count()
	if err != nil {
		return
	}

	err = o.Offset(offset).Limit(limit).Scan(result)
	return
}

func (o oauth2TokenDo) Scan(result interface{}) (err error) {
	return o.DO.Scan(result)
}

func (o oauth2TokenDo) Delete(models ...*model.Oauth2Token) (result gen.ResultInfo, err error) {
	return o.DO.Delete(models)
}

func (o *oauth2TokenDo) withDO(do gen.Dao) *oauth2TokenDo {
	o.DO = *do.(*gen.DO)
	return o
}
