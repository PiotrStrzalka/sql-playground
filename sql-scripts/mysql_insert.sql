INSERT INTO store
    (name, address, shipping_cost)
VALUES
    ('TME', 'www.tme.eu', 15.99),
    ('mouser', 'www.mouser.com', 24.99),
    ('digikey', 'www.digikey.com', 25.99);

INSERT INTO component
    (name, description, package)
VALUES
    ("Atmega8", "8-bit Atmel uC", "DIP-28"),
    ("STM32F405rg", "32-bit ST uC", "tqfp-64"),
    ("Resistor 1k", "+/- 5%", "0805");