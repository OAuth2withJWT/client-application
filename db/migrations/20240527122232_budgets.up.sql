CREATE TYPE budget_category AS ENUM (
    'monthly',
    'groceries',
    'healthcare',
    'clothing',
    'entertainment',
    'dining',
    'transport',
    'utilities'
);

CREATE TABLE budgets (
    id SERIAL PRIMARY KEY,
    user_id int NOT NULL,
    category budget_category NOT NULL,
    amount NUMERIC(10, 2) NOT NULL,
    month DATE NOT NULL,
    update_stamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);
