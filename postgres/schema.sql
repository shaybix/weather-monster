CREATE TABLE cities (
    ID int NOT NULL AUTO_INCREMENT,
    city varchar(255),
    latitude float(24) NOT NULL, 
    longitude float(24) NOT NULL,
    created_at timestamp() NOT NULL,
    updated_at timestamp(),
    version varchar(255) NOT NULL,

    PRIMARY KEY (ID)
);

CREATE TABLE temperatures (
    ID int NOT NULL AUTO_INCREMENT,
    min int NOT NULL,
    max int NOT NULL,
    city_id int,
    _timestamp timestamp NOT NULL,

    PRIMARY KEY (ID),
    FOREIGN KEY (city_id)

);

CREATE TABLE forecasts (
    ID int NOT NULL AUTO_INCREMENT,
    min int NOT NULL,
    max int NOT NULL,
    sample int NOT NULL,

    PRIMARY KEY (ID),
);