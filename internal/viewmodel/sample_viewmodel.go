package viewmodel

import "time"

// SampleRqViewModel info
// @Description Sample Data Request
type SampleRqViewModel struct {
	SampleId            string     `json:"sampleId,omitempty" example:"SampleId00001"`                                // Idetification for Sample
	SampleType          string     `json:"sampleType,omitempty" example:"Type_A" enums:"Type_A,Type_B,Type_C"`        // Type of Sample (Type_A,Type_B,Type_C)
	SampleName          string     `json:"sampleName,omitempty" example:"Sample Name 1"`                              // Name of Sample (freetext)
	SampleDescription   string     `json:"sampleDescription,omitempty" example:"Description Description Description"` // Description of Sample (freetext)
	SampleActiveVersion string     `json:"sampleActiveVersion,omitempty" example:"1.23.32.1"`                         // Current Active Version of Sample
	CreateDate          *time.Time `json:"createDate,omitempty" example:"2001-01-01 01:01:01"`                        // Created Date & Time
	CreateUser          string     `json:"createUser,omitempty" example:"11111"`                                      // Created User ID
	CreateApprover      string     `json:"createApprover,omitempty" example:"22222"`                                  // Created Approver ID
	UpdateDate          *time.Time `json:"updateDate,omitempty" example:"2002-02-02 02:02:02"`                        // Last Updated Date & Time
	UpdateUser          string     `json:"updateUser,omitempty" example:"33333"`                                      // Last Updated User ID
	UpdateApprover      string     `json:"updateApprover,omitempty" example:"44444"`
}

// SampleRsViewModel info
// @Description Sample Data Response
type SampleRsViewModel struct {
	SampleId            string                      `json:"sampleId,omitempty" example:"SampleId00001"`                                // Idetification for Sample
	SampleType          string                      `json:"sampleType,omitempty" example:"Type_A" enums:"Type_A,Type_B,Type_C"`        // Type of Sample (Type_A,Type_B,Type_C)
	SampleName          string                      `json:"sampleName,omitempty" example:"Sample Name 1"`                              // Name of Sample (freetext)
	SampleDescription   string                      `json:"sampleDescription,omitempty" example:"Description Description Description"` // Description of Sample (freetext)
	SampleActiveVersion string                      `json:"sampleActiveVersion,omitempty" example:"1.23.32.1"`                         // Current Active Version of Sample
	SampleVersions      *[]SampleVersionRsViewModel `json:"sampleVersions,omitempty"`                                                  // List of All Sample Version
	CreateDate          *time.Time                  `json:"createDate,omitempty" example:"2001-01-01 01:01:01"`                        // Created Date & Time
	CreateUser          string                      `json:"createUser,omitempty" example:"11111"`                                      // Created User ID
	CreateApprover      string                      `json:"createApprover,omitempty" example:"22222"`                                  // Created Approver ID
	UpdateDate          *time.Time                  `json:"updateDate,omitempty" example:"2002-02-02 02:02:02"`                        // Last Updated Date & Time
	UpdateUser          string                      `json:"updateUser,omitempty" example:"33333"`                                      // Last Updated User ID
	UpdateApprover      string                      `json:"updateApprover,omitempty" example:"44444"`                                  // Last Updated Approver ID
}

// SampleVersionRqViewModel info
// @Description Sample Version Data Request
type SampleVersionRqViewModel struct {
	SampleId       string     `json:"sampleId,omitempty" example:"SampleId00001"`         // Idetification for Sample
	VersionNumber  string     `json:"versionNumber,omitempty" example:"1.23.32.1"`        // Version of Sample
	CreateDate     *time.Time `json:"createDate,omitempty" example:"2001-01-01 01:01:01"` // Created Date & Time
	CreateUser     string     `json:"createUser,omitempty" example:"11111"`               // Created User ID
	CreateApprover string     `json:"createApprover,omitempty" example:"22222"`           // Created Approver ID
	UpdateDate     *time.Time `json:"updateDate,omitempty" example:"2002-02-02 02:02:02"` // Last Updated Date & Time
	UpdateUser     string     `json:"updateUser,omitempty" example:"33333"`               // Last Updated User ID
	UpdateApprover string     `json:"updateApprover,omitempty" example:"44444"`           // Last Updated Approver ID
}

// SampleVersionRsViewModel info
// @Description Sample Version Data Respponse
type SampleVersionRsViewModel struct {
	SampleId       string     `json:"sampleId,omitempty" example:"SampleId00001"`         // Idetification for Sample
	VersionNumber  string     `json:"versionNumber,omitempty" example:"1.23.32.1"`        // Version of Sample
	CreateDate     *time.Time `json:"createDate,omitempty" example:"2001-01-01 01:01:01"` // Created Date & Time
	CreateUser     string     `json:"createUser,omitempty" example:"11111"`               // Created User ID
	CreateApprover string     `json:"createApprover,omitempty" example:"22222"`           // Created Approver ID
	UpdateDate     *time.Time `json:"updateDate,omitempty" example:"2002-02-02 02:02:02"` // Last Updated Date & Time
	UpdateUser     string     `json:"updateUser,omitempty" example:"33333"`               // Last Updated User ID
	UpdateApprover string     `json:"updateApprover,omitempty" example:"44444"`           // Last Updated Approver ID
}
