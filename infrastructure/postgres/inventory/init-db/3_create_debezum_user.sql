CREATE USER dbzuser WITH PASSWORD 'dbzpw';
ALTER USER dbzuser WITH REPLICATION;
GRANT CONNECT ON DATABASE inventory_db TO dbzuser;
GRANT USAGE ON SCHEMA public TO dbzuser;
GRANT SELECT ON public.outbox TO dbzuser;