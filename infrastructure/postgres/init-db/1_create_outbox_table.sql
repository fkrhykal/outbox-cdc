CREATE TABLE public.outbox (
    id UUID PRIMARY KEY,
    aggregateid VARCHAR(255) NOT NULL,
    aggregatetype VARCHAR(255) NOT NULL,
    type VARCHAR(255) NOT NULL,
    payload JSONB NOT NULL,
    timestamp TIMESTAMP WITHOUT TIME ZONE DEFAULT now() NOT NULL
);
CREATE INDEX outbox_aggregateid_idx ON public.outbox (aggregateid);