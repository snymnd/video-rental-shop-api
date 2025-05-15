BEGIN;

CREATE TABLE IF NOT EXISTS roles(
  id SERIAL PRIMARY KEY,
  name VARCHAR NOT NULL,
  description VARCHAR,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMP
);

CREATE TABLE  IF NOT EXISTS permissions(
  id SERIAL PRIMARY KEY,
  name VARCHAR NOT NULL,
  description VARCHAR,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMP
);

CREATE TABLE  IF NOT EXISTS resources(
  id SERIAL PRIMARY KEY,
  name VARCHAR NOT NULL,
  description VARCHAR,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMP
);

CREATE TABLE  IF NOT EXISTS rbac(
  id SERIAL PRIMARY KEY,
  role_id INT NOT NULL,
  permission_id INT NOT NULL,
  resource_id INT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMP,
  FOREIGN KEY (role_id) REFERENCES roles(id),
  FOREIGN KEY (permission_id) REFERENCES permissions(id),
  FOREIGN KEY (resource_id) REFERENCES resources(id)
);


COMMIT;