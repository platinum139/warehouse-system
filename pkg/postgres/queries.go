package postgres

import (
	"log"
	"warehouse-system/pkg/models"
)

func (client *Client) CreateUser(user models.User) error {
	query := "INSERT INTO clients(external_id, username, phone, email) VALUES (?,?,?,?)"
	stmt, err := client.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.ExternalId, user.Username, user.Phone, user.Email)
	return err
}

func (client *Client) CreateManufacturer(manufacturer models.Manufacturer) {

}

func (client *Client) CreateProduct(product models.Product) {

}

func (client *Client) CreateOrder(order models.Order) {

}

func (client *Client) GetBoughtProductsQuantity() ([]models.BoughtProductsQuantity, error) {
	log.SetPrefix("[Client.GetBoughtProductsQuantity]")

	queryStr := `
        SELECT manufacturer, COUNT(DISTINCT product) AS bought_products_quantity FROM
		(SELECT orders.id, orders.quantity, products.name AS product,
		manufacturers.name AS manufacturer, clients.username AS client
		FROM orders JOIN products ON orders.product_id=products.id
		JOIN clients ON orders.client_id=clients.id
		JOIN manufacturers ON products.manufacturer_id=manufacturers.id)
		AS orders_list GROUP BY manufacturer;`

	rows, err := client.db.Query(queryStr)
	if err != nil {
		client.log.Printf("unable to query: %s\n", err)
		return nil, err
	}
	defer rows.Close()

	var quantities []models.BoughtProductsQuantity
	for rows.Next() {
		var (
			manufacturer string
			quantity     int
		)
		if err := rows.Scan(&manufacturer, &quantity); err != nil {
			log.Printf("unable to scan result: %s\n", err)
			return nil, err
		}
		quantities = append(quantities, models.BoughtProductsQuantity{
			Manufacturer:           manufacturer,
			BoughtProductsQuantity: quantity,
		})
	}
	return quantities, nil
}

func (client *Client) GetBoughtItemsQuantity() ([]models.BoughtItemsQuantity, error) {
	log.SetPrefix("[Client.GetBoughtItemsQuantity]")

	queryStr := `
        SELECT manufacturer, SUM(quantity) AS bought_items_quantity FROM
		(SELECT orders.id, orders.quantity, products.name AS product,
		manufacturers.name AS manufacturer, clients.username AS client
		FROM orders JOIN products ON orders.product_id=products.id
		JOIN clients ON orders.client_id=clients.id
		JOIN manufacturers ON products.manufacturer_id=manufacturers.id)
		AS orders_list GROUP BY manufacturer;`

	rows, err := client.db.Query(queryStr)
	if err != nil {
		client.log.Printf("unable to query: %s\n", err)
		return nil, err
	}
	defer rows.Close()

	var quantities []models.BoughtItemsQuantity
	for rows.Next() {
		var (
			manufacturer string
			quantity     int
		)
		if err := rows.Scan(&manufacturer, &quantity); err != nil {
			log.Printf("unable to scan result: %s\n", err)
			return nil, err
		}
		quantities = append(quantities, models.BoughtItemsQuantity{
			Manufacturer:        manufacturer,
			BoughtItemsQuantity: quantity,
		})
	}
	return quantities, nil
}

func (client *Client) GetExpiredProductsQuantity() {

}

func (client *Client) GetOrderedProductItemsQuantity() {

}

// SELECT orders.id, orders.quantity, products.name AS product, manufacturers.name AS manufacturer,
// clients.username AS client FROM orders JOIN products ON orders.product_id=products.id JOIN
// clients ON orders.client_id=clients.id JOIN manufacturers ON products.manufacturer_id=manufacturers.id;
