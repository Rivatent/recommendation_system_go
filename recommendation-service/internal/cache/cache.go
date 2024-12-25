package cache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"recommendation-service/internal/model"
	"time"
)

// RedisCache представляет структуру кеша с клиентом Redis и контекстом.
type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisCache создает новое подключение к Redis и возвращает указатель на RedisCache.
//
// Эта функция инициализирует новый клиент Redis, подключается к серверу
// с использованием адреса, указанного в переменной окружения REDIS_ADDR.
// Если подключение не удалось, функция завершает выполнение приложения с ошибкой.
//
// Возвращает указатель на созданный RedisCache.
func NewRedisCache() *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0,
	})
	log.Print("NewRedisCache created")

	ctx := context.Background()

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	log.Printf("Redis connected: %v", pong)

	return &RedisCache{
		client: client,
		ctx:    ctx,
	}
}

// GetRecommendationByID извлекает рекомендацию из кеша по идентификатору.
//
// Эта функция выполняет запрос к Redis, чтобы получить значение,
// связанное с указанным идентификатором. Если значение отсутствует,
// возвращается ошибка redis.Nil. В случае другой ошибки она также будет возвращена.
//
// Параметры:
// - id: уникальный идентификатор рекомендации.
//
// Возвращаемые значения:
//   - model.Recommendation: найденная рекомендация или нулевая структура,
//     если не найдено.
//   - error: ошибка запроса, если таковая произошла, или nil, если операция прошла успешно.
func (r *RedisCache) GetRecommendationByID(id string) (model.Recommendation, error) {
	var recommendation model.Recommendation
	data, err := r.client.Get(r.ctx, id).Result()
	if errors.Is(err, redis.Nil) {
		return recommendation, redis.Nil
	} else if err != nil {
		return recommendation, err
	}

	err = json.Unmarshal([]byte(data), &recommendation)
	if err != nil {
		return recommendation, err
	}

	return recommendation, nil
}

// GetRecommendationsByUserID извлекает список рекомендаций из кеша по идентификатору пользователя.
//
// Эта функция выполняет запрос к Redis, чтобы получить записи,
// связанные с указанным идентификатором пользователя. Если значение отсутствует,
// возвращается ошибка redis.Nil. В случае другой ошибки она также будет возвращена.
//
// Параметры:
// - id: уникальный идентификатор пользователя.
//
// Возвращаемые значения:
//   - []model.Recommendation: список найденных рекомендаций или пустой срез,
//     если не найдено.
//   - error: ошибка запроса, если таковая произошла, или nil, если операция прошла успешно.
func (r *RedisCache) GetRecommendationsByUserID(id string) ([]model.Recommendation, error) {
	var recommendations []model.Recommendation
	data, err := r.client.Get(r.ctx, id).Result()
	if errors.Is(err, redis.Nil) {
		return recommendations, redis.Nil
	} else if err != nil {
		return recommendations, err
	}

	err = json.Unmarshal([]byte(data), &recommendations)
	if err != nil {
		return recommendations, err
	}

	return recommendations, nil
}

// SetRecommendationByID сохраняет рекомендацию в кеш по идентификатору с указанным временем истечения.
//
// Эта функция сериализует рекомендацию в JSON и сохраняет ее в Redis с указанным идентификатором.
// Если возникает ошибка сериализации, она будет возвращена.
//
// Параметры:
// - id: уникальный идентификатор рекомендации.
// - recommendation: структура рекомендации для сохранения.
// - expiration: время истечения кеша.
//
// Возвращает ошибку, если операция не удалась, или nil, если операция прошла успешно.
func (r *RedisCache) SetRecommendationByID(id string, recommendation model.Recommendation, expiration time.Duration) error {
	data, err := json.Marshal(recommendation)
	if err != nil {
		return err
	}

	return r.client.Set(r.ctx, id, data, expiration).Err()
}

// SetRecommendationsByUserID сохраняет список рекомендаций в кеш по идентификатору пользователя с указанным временем истечения.
//
// Эта функция сериализует список рекомендаций в JSON и сохраняет его в Redis с указанным идентификатором пользователя.
// Если возникает ошибка сериализации, она будет возвращена.
//
// Параметры:
// - id: уникальный идентификатор пользователя.
// - recommendations: срез рекомендаций для сохранения.
// - expiration: время истечения кеша.
//
// Возвращает ошибку, если операция не удалась, или nil, если операция прошла успешно.
func (r *RedisCache) SetRecommendationsByUserID(id string, recommendations []model.Recommendation, expiration time.Duration) error {
	data, err := json.Marshal(recommendations)
	if err != nil {
		return err
	}

	return r.client.Set(r.ctx, id, data, expiration).Err()
}

// Close закрывает соединение с клиентом Redis.
//
// Эта функция вызывает метод Close у клиента Redis, что освобождает
// все ресурсы, используемые соединением. Возвращает ошибку, если операция не удалась, или nil, если завершена успешно.
func (r *RedisCache) Close() error {
	return r.client.Close()
}
