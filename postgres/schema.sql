CREATE TABLE cities (
    ID BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    latitude REAL NOT NULL, 
    longitude REAL NOT NULL,
    version VARCHAR(40) NOT NULL
);

CREATE TABLE temperatures (
    ID SERIAL PRIMARY KEY,
    min INT NOT NULL,
    max INT NOT NULL,
    city_id BIGINT NOT NULL REFERENCES cities (ID) ON DELETE CASCADE,
    timestamp int NOT NULL
);

CREATE TABLE webhooks (
    ID SERIAL PRIMARY KEY,
    callback_url VARCHAR(255) NOT NULL, 
    city_id BIGINT NOT NULL REFERENCES cities (ID) ON DELETE CASCADE, 
    UNIQUE (callback_url, city_id)
);