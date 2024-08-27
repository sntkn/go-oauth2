// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"
	"database/sql"

	"gorm.io/gorm"

	"gorm.io/gen"

	"gorm.io/plugin/dbresolver"
)

var (
	Q                  = new(Query)
	Oauth2Client       *oauth2Client
	Oauth2Code         *oauth2Code
	Oauth2RefreshToken *oauth2RefreshToken
	Oauth2Token        *oauth2Token
	Post               *post
	User               *user
)

func SetDefault(db *gorm.DB, opts ...gen.DOOption) {
	*Q = *Use(db, opts...)
	Oauth2Client = &Q.Oauth2Client
	Oauth2Code = &Q.Oauth2Code
	Oauth2RefreshToken = &Q.Oauth2RefreshToken
	Oauth2Token = &Q.Oauth2Token
	Post = &Q.Post
	User = &Q.User
}

func Use(db *gorm.DB, opts ...gen.DOOption) *Query {
	return &Query{
		db:                 db,
		Oauth2Client:       newOauth2Client(db, opts...),
		Oauth2Code:         newOauth2Code(db, opts...),
		Oauth2RefreshToken: newOauth2RefreshToken(db, opts...),
		Oauth2Token:        newOauth2Token(db, opts...),
		Post:               newPost(db, opts...),
		User:               newUser(db, opts...),
	}
}

type Query struct {
	db *gorm.DB

	Oauth2Client       oauth2Client
	Oauth2Code         oauth2Code
	Oauth2RefreshToken oauth2RefreshToken
	Oauth2Token        oauth2Token
	Post               post
	User               user
}

func (q *Query) Available() bool { return q.db != nil }

func (q *Query) clone(db *gorm.DB) *Query {
	return &Query{
		db:                 db,
		Oauth2Client:       q.Oauth2Client.clone(db),
		Oauth2Code:         q.Oauth2Code.clone(db),
		Oauth2RefreshToken: q.Oauth2RefreshToken.clone(db),
		Oauth2Token:        q.Oauth2Token.clone(db),
		Post:               q.Post.clone(db),
		User:               q.User.clone(db),
	}
}

func (q *Query) ReadDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Read))
}

func (q *Query) WriteDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Write))
}

func (q *Query) ReplaceDB(db *gorm.DB) *Query {
	return &Query{
		db:                 db,
		Oauth2Client:       q.Oauth2Client.replaceDB(db),
		Oauth2Code:         q.Oauth2Code.replaceDB(db),
		Oauth2RefreshToken: q.Oauth2RefreshToken.replaceDB(db),
		Oauth2Token:        q.Oauth2Token.replaceDB(db),
		Post:               q.Post.replaceDB(db),
		User:               q.User.replaceDB(db),
	}
}

type queryCtx struct {
	Oauth2Client       IOauth2ClientDo
	Oauth2Code         IOauth2CodeDo
	Oauth2RefreshToken IOauth2RefreshTokenDo
	Oauth2Token        IOauth2TokenDo
	Post               IPostDo
	User               IUserDo
}

func (q *Query) WithContext(ctx context.Context) *queryCtx {
	return &queryCtx{
		Oauth2Client:       q.Oauth2Client.WithContext(ctx),
		Oauth2Code:         q.Oauth2Code.WithContext(ctx),
		Oauth2RefreshToken: q.Oauth2RefreshToken.WithContext(ctx),
		Oauth2Token:        q.Oauth2Token.WithContext(ctx),
		Post:               q.Post.WithContext(ctx),
		User:               q.User.WithContext(ctx),
	}
}

func (q *Query) Transaction(fc func(tx *Query) error, opts ...*sql.TxOptions) error {
	return q.db.Transaction(func(tx *gorm.DB) error { return fc(q.clone(tx)) }, opts...)
}

func (q *Query) Begin(opts ...*sql.TxOptions) *QueryTx {
	tx := q.db.Begin(opts...)
	return &QueryTx{Query: q.clone(tx), Error: tx.Error}
}

type QueryTx struct {
	*Query
	Error error
}

func (q *QueryTx) Commit() error {
	return q.db.Commit().Error
}

func (q *QueryTx) Rollback() error {
	return q.db.Rollback().Error
}

func (q *QueryTx) SavePoint(name string) error {
	return q.db.SavePoint(name).Error
}

func (q *QueryTx) RollbackTo(name string) error {
	return q.db.RollbackTo(name).Error
}