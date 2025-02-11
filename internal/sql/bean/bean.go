package bean

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	ScanOutputDirectory      = "/security/devtronimagescan" // This is not configurable due to permissions for devtron-user
	NullProcessIndex         = -1
	JsonOutputFileNameSuffix = "_out.json"
	ScannerTypeClairV4       = "CLAIRV4"
	ScannerTypeClairV2       = "CLAIRV2"
	ScannerTypeTrivy         = "TRIVY"
	ScanToolClair            = "CLAIR"
	ScanToolVersion2         = "V2"
	ScanToolVersion4         = "V4"
)

type ScanExecutionType string

const (
	ScanExecutionTypeHttp ScanExecutionType = "HTTP"
	ScanExecutionTypeCli  ScanExecutionType = "CLI"
)

type ScanExecutionProcessState int

const (
	ScanExecutionProcessStateFailed    ScanExecutionProcessState = iota - 1 //resolved value = -1
	ScanExecutionProcessStateRunning                                        //resolved value =  0
	ScanExecutionProcessStateCompleted                                      //resolved value =  1
)

type ImageScanOutputObject struct {
	Name           string `json:"name"`
	Package        string `json:"package"`
	PackageVersion string `json:"packageVersion"`
	FixedInVersion string `json:"fixedInVersion"`
	Severity       string `json:"severity"`
}

type Severity int

const (
	HIGH     string = "high"
	CRITICAL string = "critical"
	SAFE     string = "safe"
	LOW      string = "low"
	MEDIUM   string = "medium"
	MODERATE string = "moderate"
)

const (
	Low Severity = iota
	Medium
	Critical
	High
	Safe
)

func (sev Severity) String() string {
	return [...]string{"low", "medium", "critical", "high", "safe"}[sev]
}
func ConvertToLowerCase(input string) string {
	return strings.ToLower(input)
}

func ConvertToSeverityUtility(severity string) Severity {
	if severity == LOW || severity == SAFE {
		return Low
	} else if severity == MEDIUM {
		return Medium
	} else if severity == HIGH || severity == CRITICAL {
		return Critical
	}
	return Low
}
func ConvertToStandardSeverityUtility(severity string) Severity {
	if severity == LOW {
		return Low
	} else if severity == MEDIUM {
		return Medium
	} else if severity == HIGH {
		return High
	} else if severity == CRITICAL {
		return Critical
	} else if severity == SAFE {
		return Safe
	}
	return Low
}

type VariableFormat string

const (
	VariableFormatString  VariableFormat = "STRING"
	VariableFormatBoolean VariableFormat = "BOOLEAN"
	VariableFormatNumber  VariableFormat = "NUMBER"
)

func ConvertVariableFormat(value string, varFormat VariableFormat) (interface{}, error) {
	switch varFormat {
	case VariableFormatString:
		return value, nil
	case VariableFormatNumber:
		return strconv.ParseFloat(value, 8)
	case VariableFormatBoolean:
		return strconv.ParseBool(value)
	default:
		return nil, fmt.Errorf("format not supported")
	}
}

type UserInfo struct {
	Id          int32    `json:"id" validate:"number"`
	EmailId     string   `json:"email_id" validate:"required"`
	Roles       []string `json:"roles,omitempty"`
	AccessToken string   `json:"access_token,omitempty"`
	Exist       bool     `json:"-"`
	UserId      int32    `json:"-"` // created or modified user id
	Status      string   `json:"status,omitempty"`
	Groups      []string `json:"groups"`
	SuperAdmin  bool     `json:"superAdmin,notnull"`
}
