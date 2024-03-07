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

type ProjectType string

func (p ProjectType) projectType() string {
	return string(p)
}

var (
	ProxyProjectType   ProjectType = "proxy"
	RegularProjectType ProjectType = "regular"
)

type ModelProject struct {
	ID   ProjectID   `gorm:"column:id;primarykey;type:string;"`
	Type ProjectType `gorm:"column:type;type:string;"`
}
