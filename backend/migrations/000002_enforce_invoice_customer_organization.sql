ALTER TABLE customers
    ADD CONSTRAINT uq_customers_id_organization
    UNIQUE (id, organization_id);

ALTER TABLE invoices
    ADD CONSTRAINT fk_invoices_customer_organization
    FOREIGN KEY (customer_id, organization_id)
    REFERENCES customers (id, organization_id);