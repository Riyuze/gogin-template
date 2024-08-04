package repository

import (
	"context"
	"gogin-template/baselib/dto"
	"gogin-template/baselib/helper"
	"gogin-template/bootstrap"
	"gogin-template/internal/model"
	"reflect"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
)

type SampleRepository interface {
	GetSamples(c context.Context, obj *model.SampleQueryModel, dtoPage dto.PageRequest) (*[]model.SampleModel, *dto.PageInfo, error)
	GetSample(c context.Context, obj *model.SampleQueryModel) (*model.SampleModel, error)
	GetSampleVersions(c context.Context, obj *model.SampleVersionQueryModel) (*[]model.SampleVersionModel, error)
	GetSampleVersion(c context.Context, obj *model.SampleVersionQueryModel) (*model.SampleVersionModel, error)
	SetSample(c context.Context, action string, obj *model.SampleModel) error
	SetSampleVersions(c context.Context, obj *[]model.SampleVersionModel) error
	SetSampleVersion(c context.Context, action string, obj *model.SampleVersionModel) error
}

type SampleRepositoryImpl struct {
	dbr      *sqlx.DB
	dbw      *sqlx.DB
	pgr      *pgxpool.Pool
	pgw      *pgxpool.Pool
	cfg      *bootstrap.Container
	queryMap map[string]string
	schema   string
}

func NewSampleRepository(dbr *sqlx.DB, dbw *sqlx.DB, pgr *pgxpool.Pool, pgw *pgxpool.Pool, cfg *bootstrap.Container) SampleRepository {
	queryMap := map[string]string{}
	columns := ""
	schema := "sample"

	// Initialize Query Map
	_, columns, _, _, _ = helper.RepoPGGetColumns(reflect.TypeOf(model.SampleModel{}))
	queryMap["GetSamples"] = `SELECT ` + columns + ` `

	_, columns, _, _, _ = helper.RepoPGGetColumns(reflect.TypeOf(model.SampleVersionModel{}))
	queryMap["GetSampleVersions"] = `SELECT ` + columns + ` `

	queryMap["SetSample"] = helper.RepoPGGetInsert(reflect.TypeOf(model.SampleModel{}), schema, "sample")

	queryMap["UpdateSample"] = helper.RepoPGGetUpsert(reflect.TypeOf(model.SampleModel{}), schema, "sample")

	queryMap["DeleteSample"] = helper.RepoPGGetDelete(reflect.TypeOf(model.SampleModel{}), schema, "sample")

	queryMap["SetSampleVersion"] = helper.RepoPGGetInsert(reflect.TypeOf(model.SampleVersionModel{}), schema, "sample_version")

	queryMap["UpdateSampleVersion"] = helper.RepoPGGetUpsert(reflect.TypeOf(model.SampleVersionModel{}), schema, "sample_version")

	queryMap["DeleteSampleVersion"] = helper.RepoPGGetDelete(reflect.TypeOf(model.SampleVersionModel{}), schema, "sample_version")

	return &SampleRepositoryImpl{dbr: dbr, dbw: dbw, cfg: cfg, pgr: pgr, pgw: pgw, queryMap: queryMap, schema: schema}
}

func (r *SampleRepositoryImpl) GetSamples(c context.Context, obj *model.SampleQueryModel, dtoPage dto.PageRequest) (*[]model.SampleModel, *dto.PageInfo, error) {
	// Set Base Query
	var data model.SampleModel
	var result []model.SampleModel

	selectQuery := r.queryMap["GetSamples"]
	baseKey, _, _, _, allowedOrder := helper.RepoPGGetColumns(reflect.TypeOf(data))
	limit, offset := helper.GetLimitAndOffset(dtoPage.PageSize, dtoPage.Page)
	baseQuery := `
	FROM ` + r.schema + `.sample 
	WHERE ($1::text is NULL OR $1::text = '' OR sample_id = $1::text)
	AND ($2::text is NULL OR $2::text = '' OR sample_type = $2::text)
	AND ($3::text is NULL OR $3::text = '' OR $3::text = '%%' 
	OR lower(sample_id) like lower($3::text)
	OR lower(sample_type) like lower($3::text)
	OR lower(sample_description) like lower($3::text))
	`
	orderString := ` ORDER BY ` + dtoPage.GetOrderString(baseKey, allowedOrder) + ` LIMIT $4::int OFFSET $5::int`
	query := selectQuery + baseQuery + orderString

	// Get Max Page
	var totalData int
	queryCount := `SELECT count(1) ` + baseQuery
	err := r.dbw.QueryRowxContext(c, queryCount, obj.SampleId, obj.SampleType, "%"+dtoPage.Query+"%").Scan(&totalData)
	if err != nil {
		return nil, nil, err
	}
	pageInfo := dtoPage.GetPageInfo(totalData)

	// Get Data
	rows, err := r.dbr.QueryxContext(c, query, obj.SampleId, obj.SampleType, "%"+dtoPage.Query+"%", limit, offset)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	// Convert to Struct
	for rows.Next() {
		err = rows.StructScan(&data)
		if err != nil {
			return nil, nil, err
		}

		versionQuery := &model.SampleVersionQueryModel{SampleId: obj.SampleId}

		versions, err := r.GetSampleVersions(c, versionQuery)
		if err != nil {
			return nil, nil, err
		}
		data.SampleVersions = versions

		result = append(result, data)
	}

	return &result, &pageInfo, nil
}

func (r *SampleRepositoryImpl) GetSample(c context.Context, obj *model.SampleQueryModel) (*model.SampleModel, error) {
	var result model.SampleModel
	var list *[]model.SampleModel
	var err error

	list, _, err = r.GetSamples(c, obj, dto.PageRequest{PageSize: 1})
	if err != nil {
		return nil, err
	}

	if len(*list) > 0 {
		result = (*list)[0]
	}

	return &result, err
}

func (r *SampleRepositoryImpl) GetSampleVersions(c context.Context, obj *model.SampleVersionQueryModel) (*[]model.SampleVersionModel, error) {
	// Set Base Query
	data := model.SampleVersionModel{}
	result := []model.SampleVersionModel{}
	selectQuery := r.queryMap["GetSampleVersions"]
	baseQuery := `
	FROM ` + r.schema + `.sample_version 
		WHERE ($1::text is NULL OR $1::text = '' OR sample_id = $1::text) 
		AND ($2::text is NULL OR $2::text = '' OR sample_version = $2::text)
	`
	query := selectQuery + baseQuery

	// Get Data
	rows, err := r.dbr.QueryxContext(c, query, obj.SampleId, obj.SampleVersions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Convert to Struct
	for rows.Next() {
		err = rows.StructScan(&data)
		if err != nil {
			return nil, err
		}
		result = append(result, data)
	}

	return &result, nil
}

func (r *SampleRepositoryImpl) GetSampleVersion(c context.Context, obj *model.SampleVersionQueryModel) (*model.SampleVersionModel, error) {
	var result model.SampleVersionModel
	var list *[]model.SampleVersionModel
	var err error

	list, err = r.GetSampleVersions(c, obj)
	if err != nil {
		return nil, err
	}

	if len(*list) > 0 {
		result = (*list)[0]
	}

	return &result, err
}

func (r *SampleRepositoryImpl) SetSample(c context.Context, action string, obj *model.SampleModel) error {
	var err error

	if strings.HasPrefix(action, "I") {
		query := r.queryMap["SetSample"]
		values := helper.RepoPGGetTypeArgValue(*obj)
		_, err = r.dbw.ExecContext(c, query, values...)
	} else if strings.HasPrefix(action, "U") {
		query := r.queryMap["UpdateSample"]
		values := helper.RepoPGGetTypeArgValue(*obj)
		_, err = r.dbw.ExecContext(c, query, values...)
	} else if strings.HasPrefix(action, "D") {
		query := r.queryMap["DeleteSample"]
		_, err = r.dbw.ExecContext(c, query, obj.SampleId)
	}

	return err
}

func (r *SampleRepositoryImpl) SetSampleVersions(c context.Context, obj *[]model.SampleVersionModel) error {
	var err error
	batch := pgx.Batch{}

	tx, err := r.pgw.BeginTx(c, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(c)

	for _, data := range *obj {
		query := r.queryMap["SetSampleVersion"]
		values := helper.RepoPGGetTypeArgValue(data)
		batch.Queue(query, values...)
	}

	if batch.Len() > 0 {
		err := r.pgw.SendBatch(c, &batch).Close()
		if err != nil {
			return err
		}
	}

	return err
}

func (r *SampleRepositoryImpl) SetSampleVersion(c context.Context, action string, obj *model.SampleVersionModel) error {
	var err error

	if strings.HasPrefix(action, "I") {
		query := r.queryMap["SetSampleVersion"]
		values := helper.RepoPGGetTypeArgValue(*obj)
		_, err = r.dbw.ExecContext(c, query, values...)
	} else if strings.HasPrefix(action, "U") {
		query := r.queryMap["UpdateSampleVersion"]
		values := helper.RepoPGGetTypeArgValue(*obj)
		_, err = r.dbw.ExecContext(c, query, values...)
	} else if strings.HasPrefix(action, "D") {
		query := r.queryMap["DeleteSampleVersion"]
		_, err = r.dbw.ExecContext(c, query, obj.SampleId, obj.VersionNumber)
	}

	return err
}
