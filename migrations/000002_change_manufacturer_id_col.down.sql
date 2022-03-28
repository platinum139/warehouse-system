ALTER TABLE orders DROP CONSTRAINT orders_product_id_fkey;

ALTER TABLE orders
RENAME product_id TO manufacturer_id;

ALTER TABLE orders ADD CONSTRAINT orders_manufacturer_id_fkey
FOREIGN KEY (manufacturer_id) REFERENCES manufacturers(id) ON DELETE CASCADE ON UPDATE CASCADE;
