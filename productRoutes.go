package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/khushi2411/zocket/database"
	"github.com/khushi2411/zocket/models"
	"github.com/lib/pq" // Import pq for pq.Array
	
)

// CreateProduct handles POST requests to create a new product
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	// Decode JSON into the product struct
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		log.Printf("Error decoding product: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}


	// Save the product to the database
	query := `
		INSERT INTO products (user_id, product_name, product_description, product_images, product_price)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := database.DB.QueryRow(
		query,
		product.UserID,
		product.ProductName,
		product.ProductDescription,
		pq.Array(product.ProductImages), // Handle TEXT[] arrays
		product.ProductPrice,
	).Scan(&product.ID)
	if err != nil {
		log.Printf("Error saving product: %v", err)
		http.Error(w, "Unable to save product to database", http.StatusInternalServerError)
		return
	}

	// Return the saved product as a response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Query to fetch paginated products
	query := `
		SELECT id, user_id, product_name, product_description, product_images, product_price
		FROM products
		LIMIT $1 OFFSET $2`
	rows, err := database.DB.Query(query, limit, offset)
	if err != nil {
		log.Printf("Error fetching products: %v", err)
		http.Error(w, "Unable to fetch products from database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		var images []string
		if err := rows.Scan(
			&product.ID,
			&product.UserID,
			&product.ProductName,
			&product.ProductDescription,
			pq.Array(&images),
			&product.ProductPrice,
		); err != nil {
			log.Printf("Error scanning product: %v", err)
			http.Error(w, "Error scanning product data", http.StatusInternalServerError)
			return
		}
		product.ProductImages = images
		products = append(products, product)
	}

	// Return paginated products as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}


// GetProductByID handles GET requests to fetch a product by its ID
func GetProductByID(w http.ResponseWriter, r *http.Request) {
	// Extract the product ID from the URL
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Trim any whitespace or newline characters
	idStr = strings.TrimSpace(idStr)

	// Convert the ID to an integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid product ID: %v", err)
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Query the database for the product
	var product models.Product
	var images []string // Temporary variable for images

	query := `
		SELECT id, user_id, product_name, product_description, product_images, product_price
		FROM products
		WHERE id = $1`
	err = database.DB.QueryRow(query, id).Scan(
		&product.ID,
		&product.UserID,
		&product.ProductName,
		&product.ProductDescription,
		pq.Array(&images), // Handle TEXT[] arrays
		&product.ProductPrice,
	)
	if err != nil {
		log.Printf("Error fetching product by ID: %v", err)
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	product.ProductImages = images // Assign the images to the product

	// Return the product as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// Extract product ID from the URL
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(strings.TrimSpace(idStr))
	if err != nil {
		log.Printf("Invalid product ID: %v", err)
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Decode the updated product details from the request body
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		log.Printf("Error decoding product: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update the product in the database
	query := `
		UPDATE products
		SET user_id = $1, product_name = $2, product_description = $3, product_images = $4, product_price = $5
		WHERE id = $6`
	_, err = database.DB.Exec(
		query,
		product.UserID,
		product.ProductName,
		product.ProductDescription,
		pq.Array(product.ProductImages),
		product.ProductPrice,
		id,
	)
	if err != nil {
		log.Printf("Error updating product: %v", err)
		http.Error(w, "Unable to update product in database", http.StatusInternalServerError)
		return
	}

	// Return a success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Product updated successfully"))
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	// Extract product ID from the URL
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(strings.TrimSpace(idStr))
	if err != nil {
		log.Printf("Invalid product ID: %v", err)
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Delete the product from the database
	query := `DELETE FROM products WHERE id = $1`
	_, err = database.DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting product: %v", err)
		http.Error(w, "Unable to delete product from database", http.StatusInternalServerError)
		return
	}

	// Return a success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Product deleted successfully"))
}
