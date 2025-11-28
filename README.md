📦 Product Inventory System

A simple full-stack Inventory Management System built with:

Backend: Go (Fiber) + PostgreSQL + GORM

Frontend: Svelte + Vite

Features:

Product creation

Variants & sub-variants

SKU-based stock management

Stock-in / stock-out

Stock report with date filter

Pagination

🚀 1. Backend Setup (Go + Fiber)
Requirements

Go 1.20+

PostgreSQL

Git

Clone the project
git clone https://github.com/Anvarsha-k/Product-Inventory-System.git
cd Product-Inventory-System

Create .env file
DATABASE_DSN=host=localhost user=postgres password=yourpass dbname=inventory port=5432 sslmode=disable
PORT=8080

Create PostgreSQL database
CREATE DATABASE inventory;

Run backend
go run ./cmd/server


Backend runs at:

http://localhost:8080

🎨 2. Frontend Setup (Svelte)
Go to frontend folder
cd frontend

Install dependencies
npm install

Run development server
npm run dev


Frontend runs at:

http://localhost:5173

📁 Project Structure (Minimal Explanation)
Product-Inventory-System/
│
├── cmd/
│   └── server/          # Entry point (main.go)
│
├── internal/
│   ├── db/              # Database connection + migration
│   ├── handlers/        # All API controllers
│   ├── models/          # GORM models (Product, Variant, SubVariant, Stock)
│
├── frontend/            # Svelte frontend
│   ├── src/pages/       # ProductList, ProductForm, StockManage, StockReport
│   └── src/App.svelte   # Routing + layout
│
└── README.md

🔑 API Endpoints (Short Overview)
Product
POST /api/products           -> Create product with variants + subvariants
GET  /api/Listproducts       -> List products (pagination)

Stock
POST /api/stock/adjust       -> Add or remove stock
GET  /api/stock/report       -> Stock report with date filter

🧪 Testing

Import this JSON for creating a sample product:

{
  "product_id": 2001,
  "product_code": "TS2001",
  "product_name": "Premium T-Shirt",
  "product_image": "https://example.com/tshirt.png",
  "variants": [
    {
      "name": "Size",
      "options": ["S", "M", "L"]
    },
    {
      "name": "Color",
      "options": ["Red", "Blue"]
    }
  ],
  "sub_variants": [
    { "sku": "TS-S-RED", "stock": "20", "option_ids": ["size-s", "color-red"] },
    { "sku": "TS-M-RED", "stock": "30", "option_ids": ["size-m", "color-red"] }
  ]
}