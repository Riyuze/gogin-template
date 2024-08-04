package model

import "time"

type SampleQueryModel struct {
	SampleId   string
	SampleType string
}

type SampleModel struct {
	SampleId            string                `db:"sample_id" dbx:"key,sort" validate:"omitempty,max=20"`
	SampleType          string                `db:"sample_type" dbx:"sort" validate:"omitempty,max=20" param:"SampleType"`
	SampleName          string                `db:"sample_name" dbx:"sort" validate:"omitempty,max=50"`
	SampleDescription   string                `db:"sample_description" validate:"omitempty,max=1000"`
	SampleActiveVersion string                `db:"sample_active_version" validate:"omitempty,max=10"`
	SampleVersions      *[]SampleVersionModel `validate:"omitempty,dive"`
	CreateDate          *time.Time            `db:"create_date"`
	CreateUser          string                `db:"create_user" validate:"omitempty,max=10"`
	CreateApprover      string                `db:"create_approver" validate:"omitempty,max=10"`
	UpdateDate          *time.Time            `db:"update_date"`
	UpdateUser          string                `db:"update_user" validate:"omitempty,max=10"`
	UpdateApprover      string                `db:"update_approver" validate:"omitempty,max=10"`
}

type SampleVersionQueryModel struct {
	SampleId       string
	SampleVersions string
}

type SampleVersionModel struct {
	SampleId       string     `db:"sample_id" dbx:"key,foreign" validate:"omitempty,max=20"`
	VersionNumber  string     `db:"version_number" dbx:"key" validate:"omitempty,max=10"`
	CreateDate     *time.Time `db:"create_date"`
	CreateUser     string     `db:"create_user" validate:"omitempty,max=10"`
	CreateApprover string     `db:"create_approver" validate:"omitempty,max=10"`
	UpdateDate     *time.Time `db:"update_date"`
	UpdateUser     string     `db:"update_user" validate:"omitempty,max=10"`
	UpdateApprover string     `db:"update_approver" validate:"omitempty,max=10"`
}
