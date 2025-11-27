CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS products (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  product_id bigint UNIQUE NOT NULL,
  product_code text UNIQUE NOT NULL,
  product_name text NOT NULL,
  product_image text,
  created_date timestamptz DEFAULT now(),
  updated_date timestamptz,
  created_user uuid,
  is_favourite boolean DEFAULT false,
  active boolean DEFAULT true,
  hsn_code text,
  total_stock numeric(20,8) DEFAULT 0
);

CREATE TABLE IF NOT EXISTS variants (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  product_id uuid NOT NULL REFERENCES products(id) ON DELETE CASCADE,
  name text NOT NULL
);

CREATE TABLE IF NOT EXISTS variant_options (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  variant_id uuid NOT NULL REFERENCES variants(id) ON DELETE CASCADE,
  value text NOT NULL
);

CREATE TABLE IF NOT EXISTS sub_variants (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  product_id uuid NOT NULL REFERENCES products(id) ON DELETE CASCADE,
  option_ids text[] NOT NULL,
  sku text UNIQUE NOT NULL,
  stock numeric(20,8) DEFAULT 0
);

CREATE TABLE IF NOT EXISTS stock_transactions (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  product_id uuid REFERENCES products(id),
  sub_variant_id uuid REFERENCES sub_variants(id),
  quantity numeric(20,8) NOT NULL,
  transaction_type text NOT NULL,
  transaction_date timestamptz DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_products_product_id ON products(product_id);
CREATE INDEX IF NOT EXISTS idx_subvariants_product_id ON sub_variants(product_id);
CREATE INDEX IF NOT EXISTS idx_stock_transactions_product_id ON stock_transactions(product_id);
CREATE INDEX IF NOT EXISTS idx_stock_transactions_date ON stock_transactions(transaction_date);
