-- Drop indexes
DROP INDEX IF EXISTS idx_clients_updated_at;
DROP INDEX IF EXISTS idx_clients_created_at;
DROP INDEX IF EXISTS idx_clients_status;
DROP INDEX IF EXISTS idx_clients_email;
DROP INDEX IF EXISTS idx_clients_name;

-- Drop table
DROP TABLE IF EXISTS clients;
