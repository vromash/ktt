CREATE TYPE marital_status_enum AS ENUM ('SINGLE', 'MARRIED', 'DIVORCED', 'COHABITING');
CREATE TYPE offer_status_enum AS ENUM ('DRAFT', 'PROCESSED');

CREATE TABLE IF NOT EXISTS applications
(
    id                         UUID PRIMARY KEY,
    created_at                 TIMESTAMPTZ         NOT NULL DEFAULT NOW(),
    updated_at                 TIMESTAMPTZ         NOT NULL DEFAULT NOW(),
    deleted_at                 TIMESTAMPTZ,
    phone                      VARCHAR(32)         NOT NULL,
    email                      VARCHAR(255)        NOT NULL,
    monthly_income             NUMERIC(15, 2)      NOT NULL,
    monthly_expenses           NUMERIC(15, 2)      NOT NULL,
    monthly_credit_liabilities NUMERIC(15, 2)      NOT NULL,
    marital_status             marital_status_enum NOT NULL,
    dependents                 INT                 NOT NULL,
    agree_to_data_sharing      BOOLEAN             NOT NULL,
    agree_to_be_scored         BOOLEAN             NOT NULL,
    amount                     NUMERIC(15, 2)      NOT NULL
);

CREATE TABLE IF NOT EXISTS offers
(
    id                     UUID PRIMARY KEY,
    created_at             TIMESTAMPTZ       NOT NULL DEFAULT NOW(),
    updated_at             TIMESTAMPTZ       NOT NULL DEFAULT NOW(),
    deleted_at             TIMESTAMPTZ,
    application_id         UUID REFERENCES applications (id),
    external_id            VARCHAR(64)       NOT NULL,
    bank                   VARCHAR(64)       NOT NULL,
    status                 offer_status_enum NOT NULL DEFAULT 'DRAFT',
    monthly_payment_amount NUMERIC(15, 2)    NOT NULL,
    total_repayment_amount NUMERIC(15, 2)    NOT NULL,
    number_of_payments     INT               NOT NULL,
    annual_percentage_rate NUMERIC(7, 3)     NOT NULL,
    first_repayment_date   DATE              NOT NULL
);
