package services

import "errors"

var (
	ErrValidation = errors.New("ошибка валидации")
	ErrNotFound   = errors.New("запись не найдена")
)
