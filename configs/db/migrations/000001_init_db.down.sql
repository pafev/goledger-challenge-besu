DROP TRIGGER IF EXISTS update_smart_contracts_updated_at ON smart_contracts;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP INDEX IF EXISTS idx_smart_contracts_updated_at;
DROP INDEX IF EXISTS idx_smart_contracts_created_at;
DROP INDEX IF EXISTS idx_smart_contracts_address;
DROP TABLE IF EXISTS smart_contracts;
