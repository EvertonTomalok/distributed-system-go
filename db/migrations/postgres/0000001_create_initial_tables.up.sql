CREATE TABLE IF NOT EXISTS methods
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    installment INT NOT NULL
);

CREATE TABLE IF NOT EXISTS orders 
(
    id uuid PRIMARY KEY,
    value DECIMAL(10,2),
    method_id INT,
    user_id VARCHAR(64),
    status BOOLEAN,
    created_at timestamp,
    updated_at timestamp,
    FOREIGN KEY(method_id) 
        REFERENCES methods(id)
);

CREATE INDEX if not exists idx_name_installment ON methods(name, installment);
CREATE INDEX if not exists idx_order_user_id on orders(user_id, id);