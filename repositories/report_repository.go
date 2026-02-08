package repositories

import (
	"database/sql"
	"kasir-api/models"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) GetReportToday() (*models.Report, error) {
	report := models.Report{
		ProdukTerlaris: make([]models.ReportProdukTerlaris, 0),
	}

	queryStats := `
		SELECT 
			COALESCE(SUM(total_amount), 0), 
			COUNT(id) 
		FROM transactions
		WHERE created_at BETWEEN 
            CURRENT_DATE::timestamp 
            AND 
            (CURRENT_DATE + INTERVAL '1 day' - INTERVAL '1 second')::timestamp
	`

	err := r.db.QueryRow(queryStats).Scan(&report.TotalRevenue, &report.TotalTransaksi)
	if err != nil {
		return nil, err
	}

	queryProducts := `
		SELECT 
			p.name,
			SUM(td.quantity) as total_qty
		FROM transaction_details td
		JOIN products p ON td.product_id = p.id
		JOIN transactions t ON td.transaction_id = t.id
		WHERE t.created_at BETWEEN 
            CURRENT_DATE::timestamp 
            AND 
            (CURRENT_DATE + INTERVAL '1 day' - INTERVAL '1 second')::timestamp
		GROUP BY p.id, p.name
		ORDER BY total_qty DESC
		LIMIT 1
	`

	rows, err := r.db.Query(queryProducts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p models.ReportProdukTerlaris
		if err := rows.Scan(&p.Nama, &p.QtyTerjual); err != nil {
			return nil, err
		}
		report.ProdukTerlaris = append(report.ProdukTerlaris, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &report, nil
}
