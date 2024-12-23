package repository

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"recommendation-service/internal/model"
)

func (r *Repo) GetRecommendationsRepo() ([]model.Recommendation, error) {
	var recommendations []model.Recommendation

	rows, err := r.db.Query("SELECT * FROM recommendations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var rec model.Recommendation
		if err := rows.Scan(&rec.ID, &rec.UserID, &rec.ProductID, &rec.Score, &rec.CreatedAt, &rec.UpdatedAt); err != nil {
			return nil, err
		}
		recommendations = append(recommendations, rec)
	}
	return recommendations, nil
}

func (r *Repo) GetRecommendationByIDRepo(id string) (model.Recommendation, error) {
	var rec model.Recommendation
	row := r.db.QueryRow("SELECT * FROM recommendations WHERE id = $1", id)
	if err := row.Scan(&rec.ID, &rec.UserID, &rec.ProductID, &rec.Score, &rec.CreatedAt, &rec.UpdatedAt); err != nil {
		return model.Recommendation{}, err
	}

	return rec, nil
}

func (r *Repo) GetRecommendationsByUserIDRepo(id string) ([]model.Recommendation, error) {
	var recs []model.Recommendation
	rows, err := r.db.Query("SELECT * FROM recommendations WHERE user_id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var rec model.Recommendation
		if err := rows.Scan(&rec.ID, &rec.UserID, &rec.ProductID, &rec.Score, &rec.CreatedAt, &rec.UpdatedAt); err != nil {
			return nil, err
		}
		recs = append(recs, rec)
	}

	return recs, nil

}
func (r *Repo) ProductUpdateMsgRepo(msg kafka.Message) error {
	var updatedProduct map[string]interface{}

	if err := json.Unmarshal(msg.Value, &updatedProduct); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}
	product, ok := updatedProduct["product"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("product field not found or is not a map")
	}
	productRating, ok := product["rating"].(float64)
	if !ok {
		return fmt.Errorf("rating not found or is not a float")
	}
	if productRating > 4.5 {
		query := `
        WITH top_products AS (
            SELECT 
                id,
                rating
            FROM products
            ORDER BY rating DESC
            LIMIT 3
        )
        INSERT INTO recommendations (user_id, product_id, score)
        SELECT u.id, tp.id, tp.rating
        FROM users u
        CROSS JOIN top_products tp
        ON CONFLICT (user_id, product_id) DO NOTHING;`

		_, err := r.db.Exec(query)
		if err != nil {
			log.Printf("Error updating recommendations: %v", err)
			return fmt.Errorf("could not update recommendations: %w", err)
		}
	}

	return nil
}

func (r *Repo) UserNewMsgRepo(msg kafka.Message) error {
	var newUser map[string]interface{}
	if err := json.Unmarshal(msg.Value, &newUser); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	user, ok := newUser["user"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("user field not found or is not a map")
	}

	userID, ok := user["id"].(string)
	if !ok {
		return fmt.Errorf("user not found or is not a string")
	}

	query := `
        WITH top_products AS (
            SELECT 
                id,
                rating,
                sales_count
            FROM products
            ORDER BY rating DESC
            LIMIT 3
        )
        INSERT INTO recommendations (user_id, product_id, score)
        SELECT $1::UUID, id, rating
        FROM top_products
        ON CONFLICT (user_id, product_id) DO NOTHING;`

	_, err := r.db.Exec(query, userID)
	if err != nil {
		log.Printf("Error updating recommendations: %v", err)
		return fmt.Errorf("could not update recommendations: %w", err)
	}

	return nil
}

func (r *Repo) ProductNewMsgRepo(msg kafka.Message) error {
	var newProduct map[string]interface{}

	if err := json.Unmarshal(msg.Value, &newProduct); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}
	product, ok := newProduct["product"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("product field not found or is not a map")
	}
	productID, ok := product["id"].(string)
	if !ok {
		return fmt.Errorf("id not found or is not a string")
	}
	productRating, ok := product["rating"].(float64)
	if !ok {
		return fmt.Errorf("rating not found or is not a float")
	}
	if productRating > 4.5 {
		query := `
	INSERT INTO recommendations (user_id, product_id, score)
	SELECT id, $1, $2 FROM users
	ON CONFLICT (user_id, product_id) DO NOTHING;`

		_, err := r.db.Exec(query, productID, productRating)
		if err != nil {
			log.Printf("Error updating recommendations: %v", err)
			return fmt.Errorf("could not update recommendations: %w", err)
		}
	}

	return nil
}
