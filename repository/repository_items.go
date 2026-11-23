package repository

import (
	"database/sql"
	"fmt"
	"items_service/model"
)

type ItemRepository struct {
	DB *sql.DB
}

func NewItemsRepository(db *sql.DB) *ItemRepository {
	return &ItemRepository{
		DB: db,
	}
}

func (r *ItemRepository) CreateItem(
	item *model.Item,
) (*model.Item, error) {

	query := `
		INSERT INTO items(id, category_id, item_name, is_active)
    	VALUES (?, ?, ?, ?)
	`
	item.IsActive = true
	_, err := r.DB.Exec(
		query,
		item.ID,
		item.CategoryID,
		item.ItemName,
		item.IsActive,
	)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat Item: %w", err)
	}

	return item, nil
}

func (r *ItemRepository) DeleteItem(id int64) error {
	query := `
        DELETE FROM Items
        WHERE id = ?
    `

	_, err := r.DB.Exec(query, id)

	if err != nil {
		return err
	}

	return nil
}

func (r *ItemRepository) FindByItemID(
	id int64,
) (*model.Item, error) {

	item := &model.Item{}
	err := r.DB.QueryRow(
		`
		SELECT id, category_id, item_name, is_active,created_at, updated_at, deleted_at
		FROM items
		WHERE id = ?
		`,
		id,
	).Scan(
		&item.ID,
		&item.CategoryID,
		&item.ItemName,
		&item.IsActive,
		&item.CreatedAt,
		&item.UpdatedAt,
		&item.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return item, nil
}

func (r *ItemRepository) FindItemByCategoryId(
	categoryId int64,
) ([]model.Item, error) {

	rows, err := r.DB.Query(
		`
        SELECT id, category_id, item_name, is_active,
				created_at, updated_at, deleted_at
        FROM items
        WHERE category_id = ?
        ORDER BY id ASC
        `,
		categoryId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.Item{}

	for rows.Next() {
		var item model.Item

		err := rows.Scan(
			&item.ID,
			&item.CategoryID,
			&item.ItemName,
			&item.IsActive,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *ItemRepository) FindByItemName(
	itemName string,
	categoryId int64,
) (*model.Item, error) {

	item := &model.Item{}
	err := r.DB.QueryRow(
		`
		SELECT id, category_id, item_name, is_active,created_at, updated_at, deleted_at
		FROM items
		WHERE item_name= ? AND category_id = ?
		`,
		itemName, categoryId,
	).Scan(
		&item.ID,
		&item.CategoryID,
		&item.ItemName,
		&item.IsActive,
		&item.CreatedAt,
		&item.UpdatedAt,
		&item.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return item, nil
}

func (r *ItemRepository) FindByCategoryIdAndItemName(
	name string,
	categoryId int64,
) (*model.Item, error) {

	item := &model.Item{}
	query := `
		SELECT id, category_id, item_name, is_active,
				created_at, updated_at, deleted_at
		FROM items
		WHERE item_name = ? AND category_id = ?
		`

	err := r.DB.QueryRow(
		query, name, categoryId,
	).Scan(
		&item.ID,
		&item.CategoryID,
		&item.ItemName,
		&item.IsActive,
		&item.CreatedAt,
		&item.UpdatedAt,
		&item.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return item, nil
}

func (r *ItemRepository) UpdateItem(item *model.Item) (*model.Item, error) {
	query := `
        UPDATE items
        SET category_id = ?, item_name = ?, is_active = ?
        WHERE id = ?
    `

	_, err := r.DB.Exec(
		query,
		item.CategoryID,
		item.ItemName,
		item.IsActive,
		item.ID,
	)

	if err != nil {
		return nil, err
	}

	return item, nil
}
