CREATE TABLE "Events" (
    "id" BIGSERIAL PRIMARY KEY,
    "id_creatore" BIGINT NULL,
    "titolo" VARCHAR(255) NULL,
    "descrizione" VARCHAR(255) NOT NULL,
    "data" TIMESTAMP NOT NULL,
    "ruoli" BIGINT NOT NULL
);

CREATE TABLE "Members" (
    "id" BIGSERIAL PRIMARY KEY,
    "email" VARCHAR(255) NOT NULL,
    "showname" VARCHAR(255) NOT NULL,
    "is_admin" BOOLEAN NOT NULL
);

ALTER TABLE "Members" ADD CONSTRAINT "members_email_uniqueemail" UNIQUE ("email");
ALTER TABLE "Members" ADD CONSTRAINT "members_showname_uniqueshowname" UNIQUE ("showname");

CREATE TABLE "Partecipation" (
    "id" BIGSERIAL PRIMARY KEY,
    "id_evento" BIGINT NULL,
    "id_partecipante" BIGINT NULL,
    "ruolo" BIGINT NULL
);

CREATE TABLE "Roles" (
    "id" BIGINT NOT NULL,
    "nome" VARCHAR(255) NOT NULL,
    PRIMARY KEY ("id")
);

CREATE TABLE "event_roles" (
    "id" BIGSERIAL PRIMARY KEY,
    "id_evento" BIGINT NOT NULL,
    "id_ruolo" BIGINT NOT NULL
    "max" BIGINT NOT NULL
);

ALTER TABLE "event_roles" ADD CONSTRAINT "event_roles_id_evento_foreign" FOREIGN KEY ("id_evento") REFERENCES "Events" ("id");
ALTER TABLE "Partecipation" ADD CONSTRAINT "partecipation_ruolo_foreign" FOREIGN KEY ("ruolo") REFERENCES "Roles" ("id");
ALTER TABLE "Partecipation" ADD CONSTRAINT "partecipation_id_partecipante_foreign" FOREIGN KEY ("id_partecipante") REFERENCES "Members" ("id");
ALTER TABLE "Events" ADD CONSTRAINT "events_id_creatore_foreign" FOREIGN KEY ("id_creatore") REFERENCES "Members" ("id");
ALTER TABLE "event_roles" ADD CONSTRAINT "event_roles_id_ruolo_foreign" FOREIGN KEY ("id_ruolo") REFERENCES "Roles" ("id");
ALTER TABLE "Partecipation" ADD CONSTRAINT "partecipation_id_evento_foreign" FOREIGN KEY ("id_evento") REFERENCES "Events" ("id");
