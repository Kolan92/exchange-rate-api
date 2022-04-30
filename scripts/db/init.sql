CREATE USER exchange_rate_api_user WITH PASSWORD '123';

CREATE DATABASE exchange_rate_api OWNER exchange_rate_api_user;

CREATE TABLE currencies_codes (
    id serial PRIMARY KEY,
    code VARCHAR(3) UNIQUE NOT NULL
);

INSERT INTO
    currencies_codes (code)
VALUES
    ('USD'),
    ('CHF'),
    ('CNY'),
    ('JPY'),
    ('KRW'),
    ('NOK'),
    ('SEK'),
    ('THB'),
    ('TWD');

CREATE TABLE exchange_rates (
    source_currency_id INT NOT NULL,
    destination_currency_id INT NOT NULL,
    date TIMESTAMP NOT NULL,
    rate NUMERIC(15, 6),
    PRIMARY KEY (
        source_currency_id,
        destination_currency_id,
        date
    ),
    FOREIGN KEY (source_currency_id) REFERENCES currencies_codes (id),
    FOREIGN KEY (destination_currency_id) REFERENCES currencies_codes (id)
);