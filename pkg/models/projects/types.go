package projects

import "fmt"

type ProjectID string

var (
	_ fmt.Stringer = ProjectID("")
)

func (p ProjectID) String() string {
	return string(p)
}

func NewProjectID(intId int) ProjectID {
	return ProjectID(fmt.Sprintf("%d", intId))
}

func NewProjectIDStr(s string) ProjectID {
	return ProjectID(s)
}

type ProjectType string

func (p ProjectType) String() string {
	return string(p)
}

var (
	ProxyProjectType    ProjectType = "proxy"
	RegularProjectType  ProjectType = "regular"
	ExternalProjectType ProjectType = "external"
)

type ModelProject struct {
	ID   ProjectID   `gorm:"column:id;primarykey;type:string;"`
	Type ProjectType `gorm:"column:type;type:string;"`
}

func (m ModelProject) TableName() string {
	return "restcol-projects"
}
