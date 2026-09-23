package model

type Revenue struct {
	Revenue float64 `json:"revenue"`
}

type ProductStats struct {
	ProductID string  `json:"product_id"`
	Revenue   float64 `json:"revenue"`
	Orders    int     `json:"orders"`
}

type OrderStats struct {
	Orders int `json:"orders"`
}

type AverageCheck struct {
	Average float64 `json:"average"`
}
