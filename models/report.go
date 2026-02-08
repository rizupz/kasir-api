package models

type Report struct {
	TotalRevenue   int                    `json:"total_revenue"`
	TotalTransaksi int                    `json:"total_transaksi"`
	ProdukTerlaris []ReportProdukTerlaris `json:"produk_terlaris"`
}

type ReportProdukTerlaris struct {
	Nama       string `json:"nama"`
	QtyTerjual int    `json:"qty_terjual"`
}
