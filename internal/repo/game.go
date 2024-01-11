package repo

import (
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/ucok-man/pixelrental/internal/contract"
	"github.com/ucok-man/pixelrental/internal/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GameServices struct {
	db *gorm.DB
}

func (s *GameServices) GetByID(gameid int) (*entity.Game, error) {
	var game entity.Game
	err := s.db.Where("game_id = $1", gameid).First(&game).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &game, nil
}

func (s *GameServices) Create(game *entity.Game) error {
	err := s.db.Create(game).Error
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "duplicate key value violates unique constraint"):
			return ErrDuplicateRecord
		default:
			return err
		}
	}

	return nil
}

func (s *GameServices) GetAll(param *contract.ReqGameGetAll) ([]*entity.Game, *contract.Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER() as totalrecords, game_id, title, description, price, year, genres, stock, created_at, version
		FROM games
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '') 
		AND (to_tsvector('simple', array_to_string(genres, ' ')) @@ plainto_tsquery('simple', $2) OR $2 = '{}')   
		ORDER BY %s %s, game_id ASC
		LIMIT $3 OFFSET $4`, param.SortColumn(), param.SortDirection(),
	)

	args := []any{param.Title, pq.Array(param.Genres), param.Limit(), param.Offset()}
	dbsql, err := s.db.DB()
	if err != nil {
		return nil, nil, err
	}

	totalRecords := 0
	games := []*entity.Game{}

	rows, err := dbsql.Query(query, args...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var game entity.Game

		err := rows.Scan(
			&totalRecords,
			&game.GameID,
			&game.Title,
			&game.Description,
			&game.Price,
			&game.Year,
			&game.Genres,
			&game.Stock,
			&game.CreatedAt,
			&game.Version,
		)
		if err != nil {
			return nil, nil, err
		}
		games = append(games, &game)
	}

	if err = rows.Err(); err != nil {
		return nil, nil, err
	}

	metadata := contract.CalculateMetadata(totalRecords, param.Page, param.PageSize)

	return games, &metadata, nil
}

// NOTE: Implement versioning to prevent race update
func (s *GameServices) Update(game *entity.Game) error {
	// versionbefore := user.Version
	err := s.db.
		Clauses(clause.Returning{}).
		Model(game).
		// Where("version = $11", versionbefore).
		Updates(&entity.Game{
			Title:       game.Title,
			Description: game.Description,
			Price:       game.Price,
			Year:        game.Year,
			Genres:      game.Genres,
			Stock:       game.Stock,
			Version:     game.Version + 1,
		}).Error
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "duplicate key value violates unique constraint"):
			return ErrDuplicateRecord
		case errors.Is(err, gorm.ErrRecordNotFound):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}
