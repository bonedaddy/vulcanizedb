-- +goose Up
CREATE TABLE public.headers
(
    id              SERIAL PRIMARY KEY,
    hash            VARCHAR(66) NOT NULL,
    block_number    BIGINT      NOT NULL,
    raw             JSONB,
    block_timestamp NUMERIC,
    eth_node_id     INTEGER     NOT NULL REFERENCES eth_nodes (id) ON DELETE CASCADE,
    created         TIMESTAMP   NOT NULL DEFAULT NOW(),
    updated         TIMESTAMP   NOT NULL DEFAULT NOW(),
    UNIQUE (block_number, eth_node_id)
);

-- +goose StatementBegin
CREATE FUNCTION set_header_updated() RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER header_updated
    BEFORE UPDATE
    ON public.headers
    FOR EACH ROW
EXECUTE PROCEDURE set_header_updated();

CREATE INDEX headers_block_number
    ON public.headers (block_number);
CREATE INDEX headers_block_timestamp_index
    ON public.headers (block_timestamp);
CREATE INDEX headers_eth_node
    ON public.headers (eth_node_id);


-- +goose Down
DROP TRIGGER header_updated ON public.headers;
DROP FUNCTION set_header_updated();

DROP TABLE public.headers;
