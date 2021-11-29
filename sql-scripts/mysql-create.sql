USE db;

CREATE TABLE store (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    address VARCHAR(100),
    shipping_cost FLOAT
);

CREATE TABLE component (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(1000),
    package VARCHAR(100)
);

CREATE TABLE price (
    ID INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    storeID INT UNSIGNED NOT NULL,
    componentID INT UNSIGNED NOT NULL,
    price_1 FLOAT NOT NULL,
    price_10 FLOAT,
    price_100 FLOAT,
    price_1000 FLOAT,
    FOREIGN KEY (storeID) REFERENCES store(id),
    FOREIGN KEY (componentID) REFERENCES component(id)
);
