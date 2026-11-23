package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"items_service/model"
	"items_service/rabbitmq"
	"items_service/repository"
	"log"
)

type ItemService struct {
	ItemRepository   *repository.ItemRepository
	LoggingPublisher *rabbitmq.LoggingPublisher
}

func NewItemsService(
	itemRepository *repository.ItemRepository,
	loggingPublisher *rabbitmq.LoggingPublisher,
) *ItemService {
	return &ItemService{
		ItemRepository:   itemRepository,
		LoggingPublisher: loggingPublisher,
	}
}

func (s *ItemService) CreateItem(item *model.Item, userID int64) (*model.Item, error) {
	existing, err := s.ItemRepository.FindByItemName(
		item.ItemName,
		item.CategoryID,
	)

	if err == nil && existing != nil {
		s.publishUserLog(
			"WARN",
			"CREATE_ITEM",
			"FAILED",
			"Item Sudah Terdaftar pada group",
			&userID,
			map[string]any{
				"reason": "ITEM_NAME_ALREADY_EXIST",
			},
		)
		return nil, fmt.Errorf("Category sudah terdaftar")
	}

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.publishUserLog(
			"WARN",
			"CREATE_ITEM",
			"FAILED",
			"Tidak dapat menemukan Item",
			&userID,
			map[string]any{
				"reason": "FAILED CHECK_ITEM",
			},
		)
		return nil, fmt.Errorf(
			"gagal memeriksa Item: %w", err,
		)
	}
	s.publishUserLog(
		"INFO",
		"CREATE_ITEM",
		"SUCCESS",
		"Item Berhasil Didaftarkan",
		&userID,
		map[string]any{
			"reason": "CREATED_ITEM_SUCCESS",
		},
	)

	return s.ItemRepository.CreateItem(item)

}

func (s *ItemService) FindItemByID(id int64, userID int64) (*model.Item, error) {
	value, err := s.ItemRepository.FindByItemID(id)
	if err != nil {
		s.publishUserLog(
			"WARN",
			"FIND_ITEM",
			"FAILED",
			"Tidak dapat menemukan ITEM",
			&userID,
			map[string]any{
				"reason": "FAILED_FIND_ITEM",
			},
		)
		return nil, fmt.Errorf("gagal menemukan Item berdasarkan ID %w", err)
	}
	s.publishUserLog(
		"INFO",
		"FIND_ITEM",
		"SUCCESS",
		"ITEM Berhasil Ditemukan",
		&userID,
		map[string]any{
			"reason": "FIND_ITEM_SUCCESS",
		},
	)
	return value, nil
}

func (s *ItemService) FindItemByCategoryId(categoryId int64, userID int64) ([]model.Item, error) {
	value, err := s.ItemRepository.FindItemByCategoryId(categoryId)
	if err != nil {
		s.publishUserLog(
			"WARN",
			"FIND_ITEM_BY_CATEGORY",
			"FAILED",
			"Tidak dapat menemukan ITEM",
			&userID,
			map[string]any{
				"reason": "FAILED_FIND_ITEM_BY_CATEGORY",
			},
		)
		return nil, fmt.Errorf("gagal menemukan Item berdasarkan Category %w", err)
	}
	s.publishUserLog(
		"INFO",
		"FIND_ITEM_BY_CATEGORY",
		"SUCCESS",
		"ITEM Berhasil Ditemukan",
		&userID,
		map[string]any{
			"reason": "FIND_ITEM_BY_CATEGORY_SUCCESS",
		},
	)
	return value, nil
}

func (s *ItemService) DeleteItem(id int64, userID int64) error {

	return s.ItemRepository.DeleteItem(id)
}

func (s *ItemService) UpdateItem(item *model.Item, userID int64) (*model.Item, error) {
	value, err := s.ItemRepository.UpdateItem(item)
	if err != nil {
		s.publishUserLog(
			"WARN",
			"UPDAYE_ITEM",
			"FAILED",
			"Gagal Meupdate Item",
			&userID,
			map[string]any{
				"reason": "FAILED_UPDATED_ITEM",
			},
		)
		return nil, fmt.Errorf("Gagal melakukan Update Item: %w", err)
	}
	s.publishUserLog(
		"INFO",
		"FIND_ITEM_BY_CATEGORY",
		"SUCCESS",
		"ITEM Berhasil DiUpdate",
		&userID,
		map[string]any{
			"reason": "UPDATED_ITEM_SUCCESS",
		},
	)

	return value, nil
}

func (s *ItemService) publishUserLog(
	level string,
	eventName string,
	status string,
	message string,
	userID *int64,
	metadata map[string]any,
) {
	if s.LoggingPublisher == nil {
		log.Println("Logging Publisher belum tersedia")
		return
	}

	logData := rabbitmq.LoggingMessage{
		ServiceName:  "Item_service",
		UserID:       userID,
		Level:        level,
		EventName:    eventName,
		Status:       status,
		Message:      message,
		ResourceType: stringPointer("user"),
		Metadata:     metadata,
	}

	err := s.LoggingPublisher.PublishLogging(
		context.Background(),
		logData,
	)

	if err != nil {
		log.Println("Gagal mengirim log registrasi:", err)
	}
}

func stringPointer(value string) *string {
	return &value
}
