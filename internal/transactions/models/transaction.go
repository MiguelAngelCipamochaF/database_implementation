package models

type Transaction struct {
	ID              int    `json:"id"`
	Codigo          string `json:"codigo"`
	Moneda          string `json:"moneda"`
	Monto           int    `json:"monto"`
	Emisor          string `json:"emisor"`
	Receptor        string `json:"receptor"`
	Fecha           string `json:"fecha"`
	Id_warehouse    int    `json:"id_warehouse"`
	Warehouse       string `json:"warehouse"`
	WarehouseAdress string `json:"warehouseAdress"`
}
