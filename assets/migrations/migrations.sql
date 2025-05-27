-- +migrate Up

CREATE TABLE Client (
                        id SERIAL PRIMARY KEY,
                        first_name VARCHAR(50) NOT NULL,
                        last_name VARCHAR(50) NOT NULL,
                        photo VARCHAR(255),
                        contact VARCHAR(100)
);

CREATE TABLE Location (
                          id SERIAL PRIMARY KEY,
                          name VARCHAR(100) NOT NULL,
                          description TEXT,
                          photo VARCHAR(255)
);

CREATE TABLE ApplicationStatus (
                        id SERIAL PRIMARY KEY,
                        status_name VARCHAR(20) NOT NULL UNIQUE
);


CREATE TABLE Application (
                             id SERIAL PRIMARY KEY,
                             fishing_date DATE NOT NULL,
                             client_id INTEGER NOT NULL,
                             location_id INTEGER NOT NULL,
                             status_id INTEGER NOT NULL,
                             created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                             CONSTRAINT fk_application_client FOREIGN KEY (client_id) REFERENCES Client(id),
                             CONSTRAINT fk_application_location FOREIGN KEY (location_id) REFERENCES Location(id),
                             CONSTRAINT fk_application_status FOREIGN KEY (status_id) REFERENCES ApplicationStatus(id)
);

-- +migrate Down

DROP TABLE Client;
DROP TABLE Location;
DROP TABLE ApplicationStatus
DROP TABLE Application;