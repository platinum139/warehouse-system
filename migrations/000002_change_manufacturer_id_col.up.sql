ALTER TABLE orders DROP CONSTRAINT orders_manufacturer_id_fkey;

ALTER TABLE orders
RENAME manufacturer_id TO product_id;

ALTER TABLE orders ADD CONSTRAINT orders_product_id_fkey
FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE ON UPDATE CASCADE;
