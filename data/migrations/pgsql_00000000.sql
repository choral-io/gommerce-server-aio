-- DROP DATABASE IF EXISTS "gommerce";
-- CREATE DATABASE "gommerce";

-- realms definition

CREATE TABLE "realms" (
    "id" VARCHAR(16) NOT NULL,
    "disabled" BOOLEAN NOT NULL DEFAULT FALSE,
    "immutable" BOOLEAN NOT NULL DEFAULT FALSE,
    "created_at" TIMESTAMP(6) NOT NULL,
    "updated_at" TIMESTAMP(6) DEFAULT NULL,
    "deleted_at" TIMESTAMP(6) DEFAULT NULL,
    "flags" BIGINT NOT NULL,
    "name" VARCHAR(64) NOT NULL,
    "title" VARCHAR(64) NOT NULL,
    "description" VARCHAR(255) DEFAULT NULL,
    CONSTRAINT "pk_realms" PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX "ix_realms_name" ON "realms" ("name");

-- realms data

INSERT INTO "realms" VALUES ('030a67b921005000', FALSE, TRUE, '2023-08-28 22:31:26.596', NULL, NULL, B'0000'::int8, 'admin', 'Admin', NULL);


-- users definition

CREATE TABLE "users" (
    "id" VARCHAR(16) NOT NULL,
    "realm_id" VARCHAR(16) NOT NULL,
    "creator_id" VARCHAR(16) DEFAULT NULL,
    "disabled" BOOLEAN NOT NULL DEFAULT FALSE,
    "approved" BOOLEAN NOT NULL DEFAULT FALSE,
    "verified" BOOLEAN NOT NULL DEFAULT FALSE,
    "immutable" BOOLEAN NOT NULL DEFAULT FALSE,
    "created_at" TIMESTAMP(6) NOT NULL,
    "updated_at" TIMESTAMP(6) DEFAULT NULL,
    "deleted_at" TIMESTAMP(6) DEFAULT NULL,
    "expires_at" TIMESTAMP(6) DEFAULT NULL,
    "first_login_time" TIMESTAMP(6) DEFAULT NULL,
    "last_active_time" TIMESTAMP(6) DEFAULT NULL,
    "flags" BIGINT NOT NULL,
    "attributes" jsonb DEFAULT NULL,
    "display_name" VARCHAR(128) DEFAULT NULL,
    "gender" VARCHAR(32) DEFAULT NULL,
    "phone_number" VARCHAR(64) DEFAULT NULL,
    "email_address" VARCHAR(64) DEFAULT NULL,
    "description" VARCHAR(255) DEFAULT NULL,
    CONSTRAINT "pk_users" PRIMARY KEY ("id"),
    CONSTRAINT "fk_users_realms_realm_id" FOREIGN KEY ("realm_id") REFERENCES "realms" ("id") ON DELETE RESTRICT ON UPDATE RESTRICT,
    CONSTRAINT "fk_users_users_creator_id" FOREIGN KEY ("creator_id") REFERENCES "users" ("id") ON DELETE RESTRICT ON UPDATE RESTRICT
);

CREATE UNIQUE INDEX "ix_users_realm_id_phone_number" ON "users" ("realm_id", "phone_number");
CREATE UNIQUE INDEX "ix_users_realm_id_email_address" ON "users" ("realm_id", "email_address");

-- users data

INSERT INTO "users" VALUES ('030a67b921005000', '030a67b921005000', NULL, FALSE, TRUE, TRUE, TRUE, '2023-08-28 22:31:26.596', NULL, NULL, NULL, NULL, NULL, 0, '{"profile.display_name": "Admin"}', 'Admin', NULL, NULL, NULL);


-- clients definition
CREATE TABLE "clients" (
    "id" VARCHAR(16) NOT NULL,
    "realm_id" VARCHAR(16) NOT NULL,
    "disabled" BOOLEAN NOT NULL DEFAULT FALSE,
    "immutable" BOOLEAN NOT NULL DEFAULT FALSE,
    "created_at" TIMESTAMP(6) NOT NULL,
    "updated_at" TIMESTAMP(6) DEFAULT NULL,
    "deleted_at" TIMESTAMP(6) DEFAULT NULL,
    "expires_at" TIMESTAMP(6) DEFAULT NULL,
    "secret_key" VARCHAR(32) NOT NULL,
    "secret_code" VARCHAR(64) DEFAULT NULL,
    "description" VARCHAR(255) DEFAULT NULL,
    CONSTRAINT "pk_clients" PRIMARY KEY ("id"),
    CONSTRAINT "fk_clients_realms_realm_id" FOREIGN KEY ("realm_id") REFERENCES "realms" ("id") ON DELETE RESTRICT ON UPDATE RESTRICT
);

CREATE UNIQUE INDEX "ix_clients_secret_key" ON "clients" ("secret_key");

-- clients data

INSERT INTO "clients" VALUES ('030a67b921005000', '030a67b921005000', FALSE, TRUE, '2023-08-28 22:31:26.596', NULL, NULL, NULL, '030a67b921005000', NULL, NULL);


-- client_users definition

CREATE TABLE "client_users" (
    "client_id" VARCHAR(16) NOT NULL,
    "user_id" VARCHAR(16) NOT NULL,
    "immutable" BOOLEAN NOT NULL DEFAULT FALSE,
    "created_at" TIMESTAMP(6) NOT NULL,
    "updated_at" TIMESTAMP(6) DEFAULT NULL,
    "deleted_at" TIMESTAMP(6) DEFAULT NULL,
    CONSTRAINT "pk_client_users" PRIMARY KEY ("client_id", "user_id"),
    CONSTRAINT "fk_client_users_clients_client_id" FOREIGN KEY ("client_id") REFERENCES "clients" ("id") ON DELETE RESTRICT ON UPDATE RESTRICT,
    CONSTRAINT "fk_client_users_users_user_id" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE RESTRICT ON UPDATE RESTRICT
);


-- roles definition

CREATE TABLE "roles" (
    "id" VARCHAR(16) NOT NULL,
    "realm_id" VARCHAR(16) NOT NULL,
    "disabled" BOOLEAN NOT NULL DEFAULT FALSE,
    "immutable" BOOLEAN NOT NULL DEFAULT FALSE,
    "created_at" TIMESTAMP(6) NOT NULL,
    "updated_at" TIMESTAMP(6) DEFAULT NULL,
    "deleted_at" TIMESTAMP(6) DEFAULT NULL,
    "name" VARCHAR(64) NOT NULL,
    "description" VARCHAR(255) DEFAULT NULL,
    CONSTRAINT "pk_roles" PRIMARY KEY ("id"),
    CONSTRAINT "fk_roles_realms_realm_id" FOREIGN KEY ("realm_id") REFERENCES "realms" ("id") ON DELETE RESTRICT ON UPDATE RESTRICT
);

CREATE UNIQUE INDEX "ix_roles_realm_id_name" ON "roles" ("realm_id", "name");

-- roles data

INSERT INTO "roles" VALUES ('030a67b921005000', '030a67b921005000', FALSE, TRUE, '2023-08-28 22:31:26.596', NULL, NULL, 'ADMIN', NULL);


-- role_users definition

CREATE TABLE "role_users" (
    "role_id" VARCHAR(16) NOT NULL,
    "user_id" VARCHAR(16) NOT NULL,
    "immutable" BOOLEAN NOT NULL DEFAULT FALSE,
    "created_at" TIMESTAMP(6) NOT NULL,
    "updated_at" TIMESTAMP(6) DEFAULT NULL,
    "deleted_at" TIMESTAMP(6) DEFAULT NULL,
    CONSTRAINT "pk_role_users" PRIMARY KEY ("role_id", "user_id"),
    CONSTRAINT "fk_role_users_roles_role_id" FOREIGN KEY ("role_id") REFERENCES "roles" ("id") ON DELETE RESTRICT ON UPDATE RESTRICT,
    CONSTRAINT "fk_role_users_users_user_id" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE RESTRICT ON UPDATE RESTRICT
);

-- role_users data

INSERT INTO "role_users" VALUES ('030a67b921005000', '030a67b921005000', TRUE, '2023-08-28 22:31:26.596', NULL, NULL);


-- logins definition

CREATE TABLE "logins" (
    "id" VARCHAR(16) NOT NULL,
    "realm_id" VARCHAR(16) NOT NULL,
    "user_id" VARCHAR(16) DEFAULT NULL,
    "disabled" BOOLEAN NOT NULL DEFAULT FALSE,
    "immutable" BOOLEAN NOT NULL DEFAULT FALSE,
    "created_at" TIMESTAMP(6) NOT NULL,
    "updated_at" TIMESTAMP(6) DEFAULT NULL,
    "deleted_at" TIMESTAMP(6) DEFAULT NULL,
    "expires_at" TIMESTAMP(6) DEFAULT NULL,
    "provider" VARCHAR(16) NOT NULL,
    "identifier" VARCHAR(64) NOT NULL,
    "credential" VARCHAR(64) DEFAULT NULL,
    "metadata" jsonb DEFAULT NULL,
    CONSTRAINT "pk_logins" PRIMARY KEY ("id"),
    CONSTRAINT "fk_logins_realms_realm_id" FOREIGN KEY ("realm_id") REFERENCES "realms" ("id") ON DELETE RESTRICT ON UPDATE RESTRICT,
    CONSTRAINT "fk_logins_users_user_id" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE RESTRICT ON UPDATE RESTRICT
);

CREATE UNIQUE INDEX "ix_logins_realm_id_provider_identifier" ON "logins" ("realm_id", "provider", "identifier");

-- logins data

INSERT INTO "logins" VALUES ('030a67b921005000', '030a67b921005000', '030a67b921005000', FALSE, FALSE, '2023-08-28 22:31:26.596', NULL, NULL, NULL, 'FORM_PASSWORD', 'admin', NULL, NULL);


-- devices definition

CREATE TABLE "devices" (
    "id" VARCHAR(16) NOT NULL,
    "user_id" VARCHAR(16) DEFAULT NULL,
    "client_id" VARCHAR(16) NOT NULL,
    "created_at" TIMESTAMP(6) NOT NULL,
    "updated_at" TIMESTAMP(6) DEFAULT NULL,
    "trace_code" VARCHAR(64) NOT NULL,
    "push_token" VARCHAR(64) DEFAULT NULL,
    "metadata" jsonb DEFAULT NULL,
    CONSTRAINT "pk_devices" PRIMARY KEY ("id"),
    CONSTRAINT "fk_devices_users_user_id" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE RESTRICT ON UPDATE RESTRICT,
    CONSTRAINT "fk_devices_users_client_id" FOREIGN KEY ("client_id") REFERENCES "clients" ("id") ON DELETE RESTRICT ON UPDATE RESTRICT
);

CREATE UNIQUE INDEX "ix_devices_trace_code" ON "devices" ("trace_code");


-- user_devices definition

CREATE TABLE "user_devices" (
    "user_id" VARCHAR(16) NOT NULL,
    "device_id" VARCHAR(16) NOT NULL,
    "created_at" TIMESTAMP(6) NOT NULL,
    "updated_at" TIMESTAMP(6) DEFAULT NULL,
    "deleted_at" TIMESTAMP(6) DEFAULT NULL,
    CONSTRAINT "pk_user_devices" PRIMARY KEY ("user_id", "device_id"),
    CONSTRAINT "fk_user_devices_users_user_id" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE RESTRICT ON UPDATE RESTRICT,
    CONSTRAINT "fk_user_devices_devices_device_id" FOREIGN KEY ("device_id") REFERENCES "devices" ("id") ON DELETE RESTRICT ON UPDATE RESTRICT
);
