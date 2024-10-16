package repository

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/Kei-K23/mock-commerce-api/db"
	"github.com/Kei-K23/mock-commerce-api/models"
	"github.com/Kei-K23/mock-commerce-api/utils"
	"github.com/jackc/pgx/v5"
)

var ErrCategoryNotFound = errors.New("category not found")

type CategoryRepository interface {
	CreateCategory(ctx context.Context, category *models.Category) (*models.Category, error)
	UpdateCategory(ctx context.Context, id int, category *models.Category) (*models.Category, error)
	GetCategoryById(ctx context.Context, id int) (*models.Category, error)
	GetAllCategories(ctx context.Context, title, limitStr, sortBy string) ([]models.Category, error)
	DeleteCategory(ctx context.Context, id int) (int, error)
}

type categoryRepository struct{}

func NewCategoryRepository() CategoryRepository {
	return &categoryRepository{}
}

func (p *categoryRepository) CreateCategory(ctx context.Context, category *models.Category) (*models.Category, error) {
	// Simulate the create data. This process will not actually create data in to database
	return &models.Category{
		ID:          11,
		Title:       category.Title,
		Description: category.Description,
		Image:       category.Image,
	}, nil
}

func (p *categoryRepository) UpdateCategory(ctx context.Context, id int, category *models.Category) (*models.Category, error) {
	return &models.Category{
		ID:          id,
		Title:       category.Title,
		Description: category.Description,
		Image:       category.Image,
	}, nil
}

func (p *categoryRepository) GetCategoryById(ctx context.Context, id int) (*models.Category, error) {
	query := `SELECT id, title, description, image FROM categories WHERE id=$1 LIMIT 1`
	row := db.Pool.QueryRow(ctx, query, id)

	var category models.Category
	// Get the category
	if err := row.Scan(
		&category.ID,
		&category.Title,
		&category.Description,
		&category.Image,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrCategoryNotFound
		}

		log.Printf("Error when fetching category by id: %v\n", err)
		return nil, err
	}

	return &category, nil
}

func (p *categoryRepository) GetAllCategories(ctx context.Context, title, limitStr, sortBy string) ([]models.Category, error) {

	// Base query
	baseQuery := "SELECT id, title, description, image FROM categories"

	qb := utils.NewQueryBuilder(baseQuery)

	if title != "" {
		qb.AddCondition("title ILIKE $%d", "%"+title+"%")
	}

	if limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			log.Fatalln("Error when parsing limit string to int: ", err)
			return nil, err
		}

		qb.SetLimit(limit)
	}

	if sortBy != "" {
		qb.SetSortBy(sortBy)
	}

	query, params := qb.Build()

	rows, err := db.Pool.Query(ctx, query, params...)
	if err != nil {
		log.Fatalln("Error fetching all categories: ", err)
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category

	for rows.Next() {
		var category models.Category
		// Get the category
		if err := rows.Scan(
			&category.ID,
			&category.Title,
			&category.Description,
			&category.Image,
		); err != nil {
			if err == pgx.ErrNoRows {
				return nil, ErrCategoryNotFound
			}
			log.Printf("Error when fetching categories: %v\n", err)
			return nil, err
		}

		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		log.Println("Error with category rows: ", err)
		return nil, err
	}

	return categories, nil
}

func (p *categoryRepository) DeleteCategory(ctx context.Context, id int) (int, error) {
	return id, nil
}
