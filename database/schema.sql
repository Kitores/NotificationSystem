CREATE TABLE users
(
    id SERIAL PRIMARY KEY,
    firstname CHARACTER VARYING(30),
    lastname CHARACTER VARYING(30),
    phonenumber CHARACTER VARYING(15),
    telegramid INTEGER,
    mail CHARACTER VARYING(30)
)