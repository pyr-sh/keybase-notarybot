package models

type Document struct {
	Signatories []*Signatory `json:"signatory"`
}

type Signatory struct {
	Name   string  `json:"name"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}
