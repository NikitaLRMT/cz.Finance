package services

import (
	"context"
	"errors"
	"time"

	"cz.Finance/backend/models"
	"cz.Finance/backend/repositories"
)

// WishlistServiceImpl реализация сервиса для работы со списком желаний
type WishlistServiceImpl struct {
	wishlistRepo repositories.WishlistRepository
	userRepo     repositories.UserRepository
}

// NewWishlistService создает новый экземпляр сервиса для работы со списком желаний
func NewWishlistService(wishlistRepo repositories.WishlistRepository, userRepo repositories.UserRepository) WishlistService {
	return &WishlistServiceImpl{
		wishlistRepo: wishlistRepo,
		userRepo:     userRepo,
	}
}

// CreateWishlistItem создает новый элемент списка желаний
func (s *WishlistServiceImpl) CreateWishlistItem(ctx context.Context, userID int64, request *models.CreateWishlistItemRequest) (*models.WishlistItem, error) {
	// Проверяем существование пользователя
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Создаем новый элемент
	item := &models.WishlistItem{
		UserID:      userID,
		Title:       request.Title,
		Price:       request.Price,
		Priority:    request.Priority,
		Description: request.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Сохраняем элемент в базе данных
	id, err := s.wishlistRepo.Create(ctx, item)
	if err != nil {
		return nil, err
	}

	// Получаем созданный элемент
	item.ID = id
	return item, nil
}

// GetWishlistItem получает элемент списка желаний по ID с проверкой принадлежности пользователю
func (s *WishlistServiceImpl) GetWishlistItem(ctx context.Context, id int64, userID int64) (*models.WishlistItem, error) {
	// Получаем элемент из базы данных
	item, err := s.wishlistRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Проверяем принадлежность элемента пользователю
	if item.UserID != userID {
		return nil, errors.New("элемент списка желаний не принадлежит пользователю")
	}

	return item, nil
}

// GetUserWishlist получает список желаний пользователя
func (s *WishlistServiceImpl) GetUserWishlist(ctx context.Context, userID int64) ([]models.WishlistItem, error) {
	return s.wishlistRepo.GetByUserID(ctx, userID)
}

// UpdateWishlistItem обновляет элемент списка желаний
func (s *WishlistServiceImpl) UpdateWishlistItem(ctx context.Context, id int64, userID int64, request *models.UpdateWishlistItemRequest) (*models.WishlistItem, error) {
	// Получаем текущий элемент
	item, err := s.GetWishlistItem(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	// Обновляем поля, если они предоставлены
	if request.Title != nil {
		item.Title = *request.Title
	}
	if request.Price != nil {
		item.Price = *request.Price
	}
	if request.Priority != nil {
		item.Priority = *request.Priority
	}
	if request.Description != nil {
		item.Description = *request.Description
	}

	// Обновляем время изменения
	item.UpdatedAt = time.Now()

	// Сохраняем изменения в базе данных
	err = s.wishlistRepo.Update(ctx, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

// DeleteWishlistItem удаляет элемент списка желаний
func (s *WishlistServiceImpl) DeleteWishlistItem(ctx context.Context, id int64, userID int64) error {
	return s.wishlistRepo.Delete(ctx, id, userID)
}
