CREATE TYPE company_type AS ENUM ('Corporations', 'NonProfit', 'Cooperative', 'Sole Proprietorship');

CREATE TABLE "companies"
(
    "id"               uuid PRIMARY KEY,
    "name"             varchar(15)  NOT NULL UNIQUE,
    "description"      varchar(3000),
    "employees_amount" int          NOT NULL,
    "registered"       bool         NOT NULL,
    "type"             company_type NOT NULL,
    "created_at"       timestamptz  NOT NULL DEFAULT now(),
    "updated_at"       timestamptz
);

-- CREATE INDEX ON "operations" ("wallet", "type", "created_at");

/*
INSERT INTO operations (wallet, type, amount)
VALUES ('wallet01', 'deposit', 100000),
       ('wallet02', 'deposit', 100000);
*/
