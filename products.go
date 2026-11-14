package main

import (
	"html/template"
	"net/http"

	"lb4/models"
)

func ProductsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := database.QueryContext(r.Context(), `
        SELECT product_id, name, description, price::text
        FROM products
        WHERE active = TRUE
        ORDER BY product_id
    `)

	if err != nil {
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []models.Product

	for rows.Next() {
		var product models.Product

		err := rows.Scan(
			&product.ProductID,
			&product.Name,
			&product.Description,
			&product.Price,
		)

		if err != nil {
			http.Error(w, "Failed to read product", http.StatusInternalServerError)
			return
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Failed to read products", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/products.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, products); err != nil {
		return
	}
}
