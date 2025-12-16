CREATE TABLE public.orders (
    id UUID PRIMARY KEY,
    item_id INT NOT NULL,
    estimated_price BIGINT NOT NULL,
    quantity INT NOT NULL,
    placed_at TIMESTAMP NOT NULL DEFAULT NOW()
);