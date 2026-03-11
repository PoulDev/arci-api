CREATE TABLE Events (
    id BIGSERIAL PRIMARY KEY,
    titolo VARCHAR(255) NOT NULL,
    descrizione VARCHAR(255) NULL,
    data TIMESTAMP NOT NULL
);

CREATE TABLE Members (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    is_admin BOOLEAN NOT NULL,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE Roles (
    id SERIAL NOT NULL,
    nome VARCHAR(255) NOT NULL UNIQUE,
    PRIMARY KEY (id)
);

CREATE TABLE Partecipation (
    id BIGSERIAL PRIMARY KEY,
    id_evento BIGINT NOT NULL,
    id_partecipante BIGINT NOT NULL,
    ruolo VARCHAR(255) NOT NULL,
    CONSTRAINT partecipation_id_evento_foreign FOREIGN KEY (id_evento) REFERENCES Events(id),
    CONSTRAINT partecipation_id_partecipante_foreign FOREIGN KEY (id_partecipante) REFERENCES Members(id),
    CONSTRAINT partecipation_ruolo_foreign FOREIGN KEY (ruolo) REFERENCES Roles(nome)
);

CREATE TABLE EventRoles (
    id BIGSERIAL PRIMARY KEY,
    id_evento BIGINT NOT NULL,
    nome_ruolo VARCHAR(255) NOT NULL,
    max BIGINT NOT NULL,
    CONSTRAINT eventroles_id_evento_foreign FOREIGN KEY (id_evento) REFERENCES Events(id)
);

GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO arci;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO arci;
GRANT ALL PRIVILEGES ON DATABASE arcidb TO arci;
