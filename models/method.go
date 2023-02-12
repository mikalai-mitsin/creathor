package models

type Method struct {
	Name    string `json:"name" yaml:"name"`
	Args    []*Param
	Results []*Param
}
