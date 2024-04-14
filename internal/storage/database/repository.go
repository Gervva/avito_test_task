package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/Gervva/avito_test_task/internal/model"
)

const errCodeUniqueViolation = "23505"

type Repository struct {
	db Database
}

func NewRepository(db Database) *Repository {
	return &Repository{
		db: db,
	}
}

func (r Repository) AddBanner(ctx context.Context, banner model.Banner) (*int64, error) {
	bannerQuery := `insert into Banner (id, feature_id, is_active, content) values ($1, $2, $3, $4) returning group_id`
	tagQuery := `insert into Tag (tag_id, banner_id) values ($1, $2)`
	bannerId := uuid.New()
	getBannerQuery := `
	SELECT b.content, b.is_active 
	FROM Banner b
	JOIN Tag t ON b.id = t.banner_id
	WHERE 
		b.feature_id = $1
		AND t.tag_id = $2 
		AND (b.feature_id, t.tag_id, b.version) IN (
			SELECT feature_id, tag_id, MAX(version) 
			FROM Banner
			JOIN Tag ON Banner.id = Tag.banner_id
			WHERE feature_id = $1 AND tag_id = $2
			GROUP BY feature_id, tag_id
		);
	`

	var groupId int64

	tagIds := *banner.TagIds

	err := r.db.WithTransaction(ctx, func(ctx context.Context) error {
		rows, _ := r.db.QueryContext(ctx, getBannerQuery, *banner.FeatureId, tagIds[0])
		if rows.Next() {
			return fmt.Errorf("%w: banner already exists", ErrRepoDB)
		}

		err := r.db.QueryRowContext(ctx, bannerQuery, bannerId, banner.FeatureId, banner.IsActive, banner.Content).Scan(&groupId)
		if err != nil {
			var pqErr *pq.Error
			if errors.As(err, &pqErr) && pqErr.Code == errCodeUniqueViolation {
				return fmt.Errorf("%w: banner already exists: %w", ErrRepoDB, err)
			}

			return fmt.Errorf("%w: error while inserting into Banner: %w", ErrRepoDB, err)
		}

		for _, tagId := range *banner.TagIds {
			_, err = r.db.ExecContext(ctx, tagQuery, tagId, bannerId)

			if err != nil {
				var pqErr *pq.Error
				if errors.As(err, &pqErr) && pqErr.Code == errCodeUniqueViolation {
					return fmt.Errorf("%w: tag already exists: %w", ErrRepoDB, err)
				}

				return fmt.Errorf("%w: error while inserting into Tag: %w", ErrRepoDB, err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("%w: error while beginning transaction while banner add: %w", ErrRepoDB, err)
	}

	return &groupId, nil
}

func (r Repository) DeleteBanner(ctx context.Context, id int64) error {
	query := `delete from Banner where group_id = $1`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: error while deleting banner: %w", ErrRepoDB, err)
	}

	numberOfDeletedSegments, err := res.RowsAffected()
	if err != nil {
		return ErrBannerNotExist
	}

	if numberOfDeletedSegments == 0 {
		return ErrBannerNotExist
	}

	return nil
}

func (r Repository) GetBanner(ctx context.Context, req model.GetBannerReq) (*[]model.GetBannerResp, error) {
	query := `
	SELECT DISTINCT
		group_id, feature_id, version, content, is_active, created_at, updated_at, 
		ARRAY(SELECT tag_id FROM Tag WHERE Tag.banner_id = Banner.id) AS tag_ids
	FROM 
		Banner
	JOIN 
		Tag ON Banner.id = Tag.banner_id
	WHERE`

	var rows *sql.Rows
	var queryError string

	if req.TagId == nil {
		query = fmt.Sprintf(`%s Banner.feature_id = %d
			AND (Banner.feature_id, Banner.version) IN (
				SELECT feature_id, MAX(version) 
				FROM Banner
				JOIN Tag ON Banner.id = Tag.banner_id
				WHERE feature_id = %d
				GROUP BY feature_id
			)`,
			query, *req.FeatureId, *req.FeatureId)
		queryError = "%w: error while getting banners by feature_id: "
	} else if req.FeatureId == nil {
		query = fmt.Sprintf(`%s Tag.tag_id = %d
			AND (Tag.tag_id, Banner.version) IN (
				SELECT tag_id, MAX(version) 
				FROM Banner
				JOIN Tag ON Banner.id = Tag.banner_id
				WHERE tag_id = %d
				GROUP BY tag_id
			)`,
			query, *req.TagId, *req.TagId)
		queryError = "%w: error while getting banners by tag_id: "
	} else {
		query = fmt.Sprintf(`%s (Tag.tag_id = %d OR Banner.feature_id = %d)
			AND (Banner.feature_id, Tag.tag_id, Banner.version) IN (
				SELECT feature_id, tag_id, MAX(version) 
				FROM Banner
				JOIN Tag ON Banner.id = Tag.banner_id
				WHERE feature_id = %d OR tag_id = %d
				GROUP BY feature_id, tag_id
			)`,
			query, *req.TagId, *req.FeatureId, *req.FeatureId, *req.TagId)
		queryError = "%w: error while getting banners by feature_id: "
	}

	if req.Offset != nil {
		query = fmt.Sprintf(`%s offset = %d;`, query, *req.Offset)
	}

	if req.Limit != nil {
		query = fmt.Sprintf(`%s limit = %d;`, query, *req.Limit)
	}

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprint(queryError, "%w"), ErrRepoDB, err)
	}

	defer func() { _ = rows.Close() }()

	banners := make([]model.GetBannerResp, 0)
	pqTagIds := pq.Int64Array{}

	for rows.Next() {
		banner := model.GetBannerResp{}

		err = rows.Scan(&banner.GroupId, &banner.FeatureId, &banner.Version, &banner.Content, &banner.IsActive, &banner.CreatedAt, &banner.UpdatedAt, &pqTagIds)
		if err != nil {
			return nil, fmt.Errorf("%w: error while scanning banner when get banner: %w", ErrRepoDB, err)
		}
		banner.TagIds = []int64(pqTagIds)

		banners = append(banners, banner)
	}

	return &banners, nil
}

func (r Repository) GetUserBanner(ctx context.Context, req model.GetUserBannerReq) (*model.GetUserBannerResp, error) {
	query := `
		SELECT b.content, b.is_active 
		FROM Banner b
		JOIN Tag t ON b.id = t.banner_id
		WHERE 
			b.feature_id = $1
			AND t.tag_id = $2 
			AND (b.feature_id, t.tag_id, b.version) IN (
				SELECT feature_id, tag_id, MAX(version) 
				FROM Banner
				JOIN Tag ON Banner.id = Tag.banner_id
				WHERE feature_id = $1 AND tag_id = $2
				GROUP BY feature_id, tag_id
			);
	`

	rows := r.db.QueryRowContext(ctx, query, *req.FeatureId, *req.TagId)

	var banner model.GetUserBannerResp
	err := rows.Scan(&banner.Content, &banner.IsActive)
	if err != nil {
		return nil, fmt.Errorf("%w: error while scanning baner when get user banner: %w", ErrRepoDB, err)
	}

	return &banner, nil
}

func (r Repository) UpdateBanner(ctx context.Context, req model.UpdateBannerReq) error {
	query := `
		SELECT DISTINCT
			Banner.*, 
			ARRAY(SELECT tag_id FROM Tag WHERE Tag.banner_id = Banner.id) AS tag_ids
		FROM 
			Banner
		JOIN 
			Tag ON Banner.id = Tag.banner_id
		WHERE group_id = $1 
		ORDER BY Banner.version DESC
	`
	bannerQuery := `insert into Banner (id, version, group_id, feature_id, is_active, content, created_at) values ($1, $2, $3, $4, $5, $6, $7)`
	tagQuery := `insert into Tag (tag_id, banner_id) values ($1, $2)`

	err := r.db.WithTransaction(ctx, func(ctx context.Context) error {
		rows, err := r.db.QueryContext(ctx, query, req.GroupId)
		if err != nil {
			return ErrBannerNotExist
		}

		defer func() { _ = rows.Close() }()

		var lastVersion model.BannerWithPK
		var OldVersion model.BannerWithPK
		numberOfVersions := 0

		pqTagIds := pq.Int64Array{}

		if rows.Next() {
			err = rows.Scan(&lastVersion.Id, &lastVersion.GroupId, &lastVersion.FeatureId, &lastVersion.Version, &lastVersion.Content, &lastVersion.IsActive, &lastVersion.CreatedAt, &lastVersion.UpdatedAt, &pqTagIds)
			if err != nil {
				return fmt.Errorf("%w: error while scanning baner when get user banner: %w", ErrRepoDB, err)
			}

			tids := []int64(pqTagIds)
			lastVersion.TagIds = &tids

			numberOfVersions++
		}

		if *req.Version != *lastVersion.Version {
			return fmt.Errorf("%w: the version of the request does not match the latest existing version", ErrRepoDB)
		}

		for rows.Next() {
			err = rows.Scan(&OldVersion.Id, &OldVersion.GroupId, &OldVersion.FeatureId, &OldVersion.Version, &OldVersion.Content, &OldVersion.IsActive, &OldVersion.CreatedAt, &OldVersion.UpdatedAt, &pqTagIds)
			if err != nil {
				return fmt.Errorf("%w: error while scanning banner when get user banner: %w", ErrRepoDB, err)
			}

			tids := []int64(pqTagIds)
			OldVersion.TagIds = &tids

			numberOfVersions++
		}
		_ = rows.Close()

		// если достигнуто максимально возможное количество версий банера, то удаляем самую старую
		if numberOfVersions == 4 {
			deleteQuery := fmt.Sprintf(`delete from Banner where id = '%s'`, OldVersion.Id.String())

			_, err := r.db.QueryContext(ctx, deleteQuery)
			if err != nil {
				return fmt.Errorf("%w: error while deleting old banner version: %w", ErrRepoDB, err)
			}

			// numberOfDeletedSegments, err := res.RowsAffected()
			// if err != nil {
			// 	return fmt.Errorf("%w: error while getting affected rows when delete old banner version: %w", ErrRepoDB, err)
			// }

			// if numberOfDeletedSegments == 0 {
			// 	return fmt.Errorf("%w: error while deleting old version: %w", ErrRepoDB, err)
			// }
		}

		newVersionId := uuid.New()
		newBannerVersion := UpdateRow(req, lastVersion)

		nv := *lastVersion.Version + 1
		newBannerVersion.Version = &nv

		// добавляем новую версию баннера
		_, err = r.db.ExecContext(
			ctx,
			bannerQuery,
			newVersionId,
			newBannerVersion.Version,
			newBannerVersion.GroupId,
			newBannerVersion.FeatureId,
			newBannerVersion.IsActive,
			newBannerVersion.Content,
			lastVersion.CreatedAt,
		)
		if err != nil {
			var pqErr *pq.Error
			if errors.As(err, &pqErr) && pqErr.Code == errCodeUniqueViolation {
				return ErrBannerAlreadyExists
			}

			return fmt.Errorf("%w: error while inserting into Banner when add new version: %w", ErrRepoDB, err)
		}

		for _, tagId := range *newBannerVersion.TagIds {
			_, err = r.db.ExecContext(ctx, tagQuery, tagId, newVersionId)

			if err != nil {
				var pqErr *pq.Error
				if errors.As(err, &pqErr) && pqErr.Code == errCodeUniqueViolation {
					return ErrBannerAlreadyExists
				}

				return fmt.Errorf("%w: error while inserting into Tag when add new version: %w", ErrRepoDB, err)
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("%w: error while beginning transaction when banner apdate: %w", ErrRepoDB, err)
	}

	return nil
}

func (r Repository) GetBannerVersion(ctx context.Context, req *model.GetBannerVersionReq) (*model.GetBannerVersionResp, error) {
	query := `
		SELECT DISTINCT
			group_id, feature_id, version, content, is_active, created_at, updated_at, 
			ARRAY(SELECT tag_id FROM Tag WHERE Tag.banner_id = Banner.id) AS tag_ids
		FROM 
			Banner
		JOIN 
			Tag ON Banner.id = Tag.banner_id
		WHERE Banner.version = $1 AND Banner.Group_id = $2
	`

	var banner model.GetBannerVersionResp
	pqTagIds := pq.Int64Array{}

	row := r.db.QueryRowContext(ctx, query, *req.Version, *req.GroupId)
	err := row.Scan(&banner.GroupId, &banner.FeatureId, &banner.Version, &banner.Content, &banner.IsActive, &banner.CreatedAt, &banner.UpdatedAt, &pqTagIds)
	if err != nil {
		return nil, ErrBannerNotExist
	}

	tids := []int64(pqTagIds)
	banner.TagIds = &tids

	return &banner, nil
}

func (r Repository) GetAllVersions(ctx context.Context, req *model.GetAllVersionsReq) (*[]model.GetAllVersionsResp, error) {
	query := `
		SELECT DISTINCT
			group_id, feature_id, version, content, is_active, created_at, updated_at, 
			ARRAY(SELECT tag_id FROM Tag WHERE Tag.banner_id = Banner.id) AS tag_ids
		FROM 
			Banner
		JOIN 
			Tag ON Banner.id = Tag.banner_id
		WHERE Banner.Group_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, *req.GroupId)
	if err != nil {
		return nil, ErrBannerNotExist
	}

	defer func() { _ = rows.Close() }()

	banners := make([]model.GetAllVersionsResp, 0, 4)
	pqTagIds := pq.Int64Array{}

	for rows.Next() {
		var banner model.GetAllVersionsResp

		err = rows.Scan(&banner.GroupId, &banner.FeatureId, &banner.Version, &banner.Content, &banner.IsActive, &banner.CreatedAt, &banner.UpdatedAt, &pqTagIds)
		if err != nil {
			return nil, fmt.Errorf("%w: error while scanning banner when get banner: %w", ErrRepoDB, err)
		}

		tids := []int64(pqTagIds)
		banner.TagIds = &tids

		banners = append(banners, banner)
	}

	return &banners, nil
}

func (r Repository) DeleteByFeatureTag(ctx context.Context, req *model.DeleteByFeatureTagReq) error {
	var query string
	var queryError string

	if req.TagId == nil {
		query = fmt.Sprintf(`
			DELETE FROM Banner
			WHERE feature_id = %d`,
			*req.FeatureId)
		queryError = "%w: error while delete banners by feature_id: "
	} else {
		query = fmt.Sprintf(`
			DELETE FROM Banner
			WHERE id IN (
				SELECT b.id
				FROM Banner b
				INNER JOIN Tag t ON b.id = t.banner_id
				WHERE t.tag_id = %d
			)
		`, *req.TagId)
		queryError = "%w: error while delete banners by tag_id: "
	}

	res, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf(fmt.Sprint(queryError, "%w"), ErrRepoDB, err)
	}

	numberOfDeletedSegments, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: error while getting affected rows when delete banner: %w", ErrRepoDB, err)
	}

	if numberOfDeletedSegments == 0 {
		return ErrBannerNotExist
	}

	return nil
}
