package models

type Modeler interface {
	Save() error
	Update() error
	Delete() error
	FindByID(id int) error
}
