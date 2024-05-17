package repositories

import (
	"context"
	"database/sql"

	"github.com/tsanaativa/sehatin-backend-v0.1/custom_errors"
	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
)

type ChatRepoOpts struct {
	Db *sql.DB
}

type ChatRepository interface {
	FindAllConsultationChat(ctx context.Context, consultationId int64) ([]entities.Chat, error)
	CreateOne(ctx context.Context, chat entities.Chat) error
}

type ChatRepositoryPostgres struct {
	db *sql.DB
}

func NewChatRepositoryPostgres(cOpts *ChatRepoOpts) ChatRepository {
	return &ChatRepositoryPostgres{
		db: cOpts.Db,
	}
}

func (r *ChatRepositoryPostgres) FindAllConsultationChat(ctx context.Context, consultationId int64) ([]entities.Chat, error) {
	chats := []entities.Chat{}

	var err error
	var rows *sql.Rows

	tx := extractTx(ctx)
	if tx != nil {
		rows, err = tx.QueryContext(ctx, qFindAllChatByConsultationId, consultationId)
	} else {
		rows, err = r.db.QueryContext(ctx, qFindAllChatByConsultationId, consultationId)
	}
	defer rows.Close()

	for rows.Next() {
		chat := entities.Chat{}
		rows.Scan(&chat.Id, &chat.IsFromUser, &chat.Content, &chat.Type, &chat.CreatedAt)
		chats = append(chats, chat)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return chats, nil
}

func (r *ChatRepositoryPostgres) CreateOne(ctx context.Context, chat entities.Chat) error {
	var err error

	values := []interface{}{}
	values = append(values, chat.ConsultationId)
	values = append(values, chat.IsFromUser)
	values = append(values, chat.Content)
	values = append(values, chat.Type)

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qCreateOneChat, values...).Scan(&chat.Id)
	} else {
		err = r.db.QueryRowContext(ctx, qCreateOneChat, values...).Scan(&chat.Id)
	}

	if err != nil {
		return err
	}

	return nil
}
