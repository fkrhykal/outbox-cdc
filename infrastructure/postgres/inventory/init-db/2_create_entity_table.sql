CREATE TABLE public.products (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    stock INT NOT NULL,
    price INT NOT NULL
);

CREATE TABLE public.reservations (
    id UUID PRIMARY KEY,
    reservation_key UUID NOT NULL,
    product_id UUID NOT NULL REFERENCES public.products(id),
    price_level INT NOT NULL,
    quantity INT NOT NULL
);

INSERT INTO public.products (id, name, stock, price) VALUES
('550e8400-e29b-41d4-a716-446655440000', 'Laptop Pro', 50, 1500),
('6ba7b810-9dad-11d1-80b4-00c04fd430c8', 'Wireless Mouse', 200, 25),
('789e4567-e89b-12d3-a456-426614174000', 'Mechanical Keyboard', 75, 120),
('a1b2c3d4-e5f6-4a5b-bcde-f12345678901', 'UltraWide Monitor', 30, 450),
('b2c3d4e5-f6a7-4b6c-cdef-0123456789ab', 'USB-C Docking Station', 100, 180),
('c3d4e5f6-a7b8-4c7d-def0-123456789abc', 'Noise Cancelling Headphones', 60, 300),
('d4e5f6a7-b8c9-4d8e-ef01-23456789abcd', 'Ergonomic Desk Chair', 40, 350),
('e5f6a7b8-c9d0-4e9f-f012-3456789abcde', 'Desk Lamp with Charger', 120, 45);