-- =============================================================
-- KERIPIK PANGSIT - DATABASE SCHEMA
-- =============================================================

-- Drop existing tables (urutan foreign key)
DROP TABLE IF EXISTS redeem_logs CASCADE;
DROP TABLE IF EXISTS transactions CASCADE;
DROP TABLE IF EXISTS products CASCADE;
DROP TABLE IF EXISTS customers CASCADE;

-- =============================================================
-- TABLE: customers
-- =============================================================
CREATE TABLE customers (
    id              VARCHAR(36) PRIMARY KEY,
    nama_customer   VARCHAR(100) NOT NULL,
    poin            INTEGER NOT NULL DEFAULT 0,
    tanggal_daftar  TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_nama_customer UNIQUE (nama_customer)
);

-- =============================================================
-- TABLE: products
-- =============================================================
CREATE TABLE products (
    id           VARCHAR(36) PRIMARY KEY,
    nama_produk  VARCHAR(100) NOT NULL,
    tipe_produk  VARCHAR(100) NOT NULL,
    rasa         VARCHAR(50) NOT NULL,
    ukuran       VARCHAR(20) NOT NULL CHECK (ukuran IN ('Small', 'Medium', 'Large')),
    harga        NUMERIC(12, 2) NOT NULL,
    stok         INTEGER NOT NULL DEFAULT 0,
    tanggal_buat DATE NOT NULL,
    CONSTRAINT chk_harga_positive CHECK (harga > 0),
    CONSTRAINT chk_stok_positive  CHECK (stok >= 0)
);

-- =============================================================
-- TABLE: transactions
-- =============================================================
CREATE TABLE transactions (
    id                 VARCHAR(36) PRIMARY KEY,
    customer_id        VARCHAR(36) NOT NULL REFERENCES customers(id),
    nama_customer      VARCHAR(100) NOT NULL,
    product_id         VARCHAR(36) NOT NULL REFERENCES products(id),
    nama_produk        VARCHAR(100) NOT NULL,
    ukuran_produk      VARCHAR(20) NOT NULL,
    rasa               VARCHAR(50) NOT NULL,
    quantity           INTEGER NOT NULL DEFAULT 1,
    harga              NUMERIC(12, 2) NOT NULL,
    total_harga        NUMERIC(12, 2) NOT NULL,
    tanggal_transaksi  TIMESTAMP NOT NULL DEFAULT NOW(),
    is_customer_baru   BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT chk_quantity_positive CHECK (quantity > 0)
);

-- =============================================================
-- TABLE: redeem_logs
-- =============================================================
CREATE TABLE redeem_logs (
    id              VARCHAR(36) PRIMARY KEY,
    customer_id     VARCHAR(36) NOT NULL REFERENCES customers(id),
    nama_customer   VARCHAR(100) NOT NULL,
    product_id      VARCHAR(36) NOT NULL REFERENCES products(id),
    nama_produk     VARCHAR(100) NOT NULL,
    ukuran_produk   VARCHAR(20) NOT NULL,
    poin_digunakan  INTEGER NOT NULL,
    tanggal_redeem  TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_poin_positive CHECK (poin_digunakan > 0)
);

-- =============================================================
-- INDEXES
-- =============================================================
CREATE INDEX idx_products_tanggal_buat   ON products(tanggal_buat);
CREATE INDEX idx_products_ukuran         ON products(ukuran);
CREATE INDEX idx_transactions_tanggal    ON transactions(tanggal_transaksi);
CREATE INDEX idx_transactions_customer   ON transactions(customer_id);
CREATE INDEX idx_transactions_product    ON transactions(product_id);
CREATE INDEX idx_customers_nama          ON customers(nama_customer);
CREATE INDEX idx_redeem_customer         ON redeem_logs(customer_id);
CREATE INDEX idx_redeem_tanggal          ON redeem_logs(tanggal_redeem);

-- =============================================================
-- SEED DATA
-- =============================================================

-- Customers
INSERT INTO customers (id, nama_customer, poin, tanggal_daftar) VALUES
    ('c1000000-0000-0000-0000-000000000001', 'Fery',  20, '2025-10-22 15:00:22'),
    ('c1000000-0000-0000-0000-000000000002', 'Fenty', 25, '2025-11-22 13:00:22'),
    ('c1000000-0000-0000-0000-000000000003', 'Kunjo', 35, '2025-12-22 11:00:00');

-- Products (stok besar untuk demo)
INSERT INTO products (id, nama_produk, tipe_produk, rasa, ukuran, harga, stok, tanggal_buat) VALUES
    ('p1000000-0000-0000-0000-000000000001', 'Keripik Pangsit', 'Keripik Pangsit', 'Jagung Bakar', 'Small',  10000, 100, '2025-10-01'),
    ('p1000000-0000-0000-0000-000000000002', 'Keripik Pangsit', 'Keripik Pangsit', 'Jagung Bakar', 'Medium', 25000, 100, '2025-10-01'),
    ('p1000000-0000-0000-0000-000000000003', 'Keripik Pangsit', 'Keripik Pangsit', 'Jagung Bakar', 'Large',  35000, 100, '2025-10-01'),
    ('p1000000-0000-0000-0000-000000000004', 'Keripik Pangsit', 'Keripik Pangsit', 'Rumput Laut',  'Small',  10000, 100, '2025-10-01'),
    ('p1000000-0000-0000-0000-000000000005', 'Keripik Pangsit', 'Keripik Pangsit', 'Rumput Laut',  'Medium', 25000, 100, '2025-10-01'),
    ('p1000000-0000-0000-0000-000000000006', 'Keripik Pangsit', 'Keripik Pangsit', 'Rumput Laut',  'Large',  35000, 100, '2025-10-01'),
    ('p1000000-0000-0000-0000-000000000007', 'Keripik Pangsit', 'Keripik Pangsit', 'Original',     'Small',  10000, 100, '2025-10-01'),
    ('p1000000-0000-0000-0000-000000000008', 'Keripik Pangsit', 'Keripik Pangsit', 'Original',     'Medium', 25000, 100, '2025-10-01'),
    ('p1000000-0000-0000-0000-000000000009', 'Keripik Pangsit', 'Keripik Pangsit', 'Original',     'Large',  35000, 100, '2025-10-01'),
    ('p1000000-0000-0000-0000-000000000010', 'Keripik Pangsit', 'Keripik Pangsit', 'Jagung Manis', 'Small',  10000, 100, '2025-10-01'),
    ('p1000000-0000-0000-0000-000000000011', 'Keripik Pangsit', 'Keripik Pangsit', 'Jagung Manis', 'Medium', 25000, 100, '2025-10-01'),
    ('p1000000-0000-0000-0000-000000000012', 'Keripik Pangsit', 'Keripik Pangsit', 'Jagung Manis', 'Large',  35000, 100, '2025-10-01'),
    ('p1000000-0000-0000-0000-000000000013', 'Keripik Pangsit', 'Keripik Pangsit', 'Keju Asin',    'Small',  10000, 100, '2025-10-01'),
    ('p1000000-0000-0000-0000-000000000014', 'Keripik Pangsit', 'Keripik Pangsit', 'Keju Asin',    'Medium', 25000, 100, '2025-10-01'),
    ('p1000000-0000-0000-0000-000000000015', 'Keripik Pangsit', 'Keripik Pangsit', 'Keju Asin',    'Large',  35000, 100, '2025-10-01'),
    ('p1000000-0000-0000-0000-000000000016', 'Keripik Pangsit', 'Keripik Pangsit', 'Keju Manis',   'Small',  10000, 100, '2025-10-01'),
    ('p1000000-0000-0000-0000-000000000017', 'Keripik Pangsit', 'Keripik Pangsit', 'Keju Manis',   'Medium', 25000, 100, '2025-10-01'),
    ('p1000000-0000-0000-0000-000000000018', 'Keripik Pangsit', 'Keripik Pangsit', 'Keju Manis',   'Large',  35000, 100, '2025-10-01'),
    ('p1000000-0000-0000-0000-000000000019', 'Keripik Pangsit', 'Keripik Pangsit', 'Pedas',        'Small',  10000, 100, '2025-10-01'),
    ('p1000000-0000-0000-0000-000000000020', 'Keripik Pangsit', 'Keripik Pangsit', 'Pedas',        'Medium', 25000, 100, '2025-10-01'),
    ('p1000000-0000-0000-0000-000000000021', 'Keripik Pangsit', 'Keripik Pangsit', 'Pedas',        'Large',  35000, 100, '2025-10-01');

-- Sample Transactions (sesuai requirement)
INSERT INTO transactions (id, customer_id, nama_customer, product_id, nama_produk, ukuran_produk, rasa, quantity, harga, total_harga, tanggal_transaksi, is_customer_baru) VALUES
    ('t1000000-0000-0000-0000-000000000001',
     'c1000000-0000-0000-0000-000000000001', 'Fery',
     'p1000000-0000-0000-0000-000000000001', 'Keripik Pangsit',
     'Small', 'Jagung Bakar', 2, 10000, 20000,
     '2025-10-22 15:00:22', TRUE),
    ('t1000000-0000-0000-0000-000000000002',
     'c1000000-0000-0000-0000-000000000002', 'Fenty',
     'p1000000-0000-0000-0000-000000000002', 'Keripik Pangsit',
     'Medium', 'Jagung Bakar', 1, 25000, 25000,
     '2025-11-22 13:00:22', FALSE),
    ('t1000000-0000-0000-0000-000000000003',
     'c1000000-0000-0000-0000-000000000003', 'Kunjo',
     'p1000000-0000-0000-0000-000000000003', 'Keripik Pangsit',
     'Large', 'Jagung Bakar', 1, 35000, 35000,
     '2025-12-22 11:00:00', FALSE);

-- Update stok sesuai transaksi
UPDATE products SET stok = stok - 2 WHERE id = 'p1000000-0000-0000-0000-000000000001';
UPDATE products SET stok = stok - 1 WHERE id = 'p1000000-0000-0000-0000-000000000002';
UPDATE products SET stok = stok - 1 WHERE id = 'p1000000-0000-0000-0000-000000000003';
