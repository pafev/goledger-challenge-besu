CREATE TABLE smart_contracts (
    smart_contract_id BIGSERIAL PRIMARY KEY,
    address VARCHAR(255) NOT NULL UNIQUE,
    value NUMERIC(78, 0) NOT NULL, -- suficiente para inteiro de 256 bits (uint256 na abi)
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_smart_contracts_address ON smart_contracts(address);
CREATE INDEX idx_smart_contracts_created_at ON smart_contracts(created_at);
CREATE INDEX idx_smart_contracts_updated_at ON smart_contracts(updated_at);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_smart_contracts_updated_at
    BEFORE UPDATE ON smart_contracts
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
