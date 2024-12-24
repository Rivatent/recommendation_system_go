package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"recommendation-service/internal/model"
	"testing"
	"time"
)

func TestGetRecommendationsRepo(t *testing.T) {

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &Repo{db: db}

	t.Run("success", func(t *testing.T) {

		rows := sqlmock.NewRows([]string{"id", "user_id", "product_id", "score", "created_at", "updated_at"}).
			AddRow(1, 1, 2, 4.5, time.Now(), time.Now()).
			AddRow(2, 2, 3, 3.5, time.Now(), time.Now())

		mock.ExpectQuery("SELECT \\* FROM recommendations").WillReturnRows(rows)

		recommendations, err := repo.GetRecommendationsRepo()
		require.NoError(t, err)
		assert.NotNil(t, recommendations)
		assert.Len(t, recommendations, 2)

		assert.Equal(t, "1", recommendations[0].ID)
		assert.Equal(t, "1", recommendations[0].UserID)
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM recommendations").WillReturnError(sql.ErrNoRows)

		recommendations, err := repo.GetRecommendationsRepo()
		require.Error(t, err)
		assert.Nil(t, recommendations)
	})

	t.Run("scan error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "user_id", "product_id", "score", "created_at", "updated_at"}).
			AddRow(1, 1, 2, "InvalidScore", time.Now(), time.Now())

		mock.ExpectQuery("SELECT \\* FROM recommendations").WillReturnRows(rows)

		recommendations, err := repo.GetRecommendationsRepo()
		require.Error(t, err)
		assert.Nil(t, recommendations)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetRecommendationByIDRepo(t *testing.T) {

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &Repo{db: db}

	t.Run("success", func(t *testing.T) {

		expectedRec := model.Recommendation{
			ID:        "1",
			UserID:    "1",
			ProductID: "2",
			Score:     4.5,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		rows := sqlmock.NewRows([]string{"id", "user_id", "product_id", "score", "created_at", "updated_at"}).
			AddRow(expectedRec.ID, expectedRec.UserID, expectedRec.ProductID, expectedRec.Score, expectedRec.CreatedAt, expectedRec.UpdatedAt)

		mock.ExpectQuery("SELECT \\* FROM recommendations WHERE id = \\$1").
			WithArgs(expectedRec.ID).
			WillReturnRows(rows)

		rec, err := repo.GetRecommendationByIDRepo("1")
		require.NoError(t, err)
		assert.Equal(t, expectedRec, rec)
	})

	t.Run("error scanning", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM recommendations WHERE id = \\$1").
			WithArgs("1").
			WillReturnError(errors.New("scan error"))

		_, err := repo.GetRecommendationByIDRepo("1")
		assert.Error(t, err)
		assert.Equal(t, "scan error", err.Error())
	})

	t.Run("no rows", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM recommendations WHERE id = \\$1").
			WithArgs("999").
			WillReturnRows(sqlmock.NewRows([]string{}))

		rec, err := repo.GetRecommendationByIDRepo("999")
		assert.Error(t, err)
		assert.Equal(t, model.Recommendation{}, rec)
	})
}

func TestGetRecommendationsByUserIDRepo(t *testing.T) {

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &Repo{db: db}

	t.Run("success", func(t *testing.T) {
		userID := "1"

		expectedRecs := []model.Recommendation{
			{ID: "1", UserID: userID, ProductID: "2", Score: 4.5, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{ID: "2", UserID: userID, ProductID: "3", Score: 5.0, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		}

		rows := sqlmock.NewRows([]string{"id", "user_id", "product_id", "score", "created_at", "updated_at"}).
			AddRow(expectedRecs[0].ID, expectedRecs[0].UserID, expectedRecs[0].ProductID, expectedRecs[0].Score, expectedRecs[0].CreatedAt, expectedRecs[0].UpdatedAt).
			AddRow(expectedRecs[1].ID, expectedRecs[1].UserID, expectedRecs[1].ProductID, expectedRecs[1].Score, expectedRecs[1].CreatedAt, expectedRecs[1].UpdatedAt)

		mock.ExpectQuery("SELECT \\* FROM recommendations WHERE user_id = \\$1").WithArgs(userID).WillReturnRows(rows)

		recs, err := repo.GetRecommendationsByUserIDRepo(userID)
		require.NoError(t, err)
		assert.Equal(t, expectedRecs, recs)

		err = mock.ExpectationsWereMet()
		require.NoError(t, err)
	})

	t.Run("db_error", func(t *testing.T) {
		userID := "1"

		mock.ExpectQuery("SELECT \\* FROM recommendations WHERE user_id = \\$1").WithArgs(userID).WillReturnError(errors.New("database error"))

		recs, err := repo.GetRecommendationsByUserIDRepo(userID)
		assert.Nil(t, recs)
		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())

		err = mock.ExpectationsWereMet()
		require.NoError(t, err)
	})
}

func TestProductUpdateMsgRepo(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &Repo{db: db}

	t.Run("TestNoProductMap", func(t *testing.T) {
		err := repo.ProductUpdateMsgRepo(map[string]interface{}{})
		require.Error(t, err)
		require.Equal(t, "product field not found or is not a map", err.Error())
	})

	t.Run("TestNoRating", func(t *testing.T) {
		err := repo.ProductUpdateMsgRepo(map[string]interface{}{
			"product": map[string]interface{}{},
		})
		require.Error(t, err)
		require.Equal(t, "rating not found or is not a float", err.Error())
	})

	t.Run("TestRatingLow", func(t *testing.T) {
		err := repo.ProductUpdateMsgRepo(map[string]interface{}{
			"product": map[string]interface{}{
				"rating": 4.0,
			},
		})
		require.NoError(t, err)
	})
	t.Run("TestSuccessfulUpdate", func(t *testing.T) {
		expectedQuery := "WITH top_products AS \\( SELECT id, rating FROM products ORDER BY rating DESC LIMIT 3 \\) INSERT INTO recommendations \\(user_id, product_id, score\\) SELECT u.id, tp.id, tp.rating FROM users u CROSS JOIN top_products tp ON CONFLICT \\(user_id, product_id\\) DO NOTHING;"
		mock.ExpectExec(expectedQuery).WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.ProductUpdateMsgRepo(map[string]interface{}{
			"product": map[string]interface{}{
				"rating": 4.6,
			},
		})

		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("TestUpdateError", func(t *testing.T) {
		mock.ExpectExec(queryProductUpdateMsg).WillReturnError(fmt.Errorf("some error"))
		err := repo.ProductUpdateMsgRepo(map[string]interface{}{
			"product": map[string]interface{}{
				"rating": 4.6,
			},
		})
		require.Error(t, err)
	})
}
