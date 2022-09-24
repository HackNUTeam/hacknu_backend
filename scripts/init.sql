CREATE TABLE IF NOT EXISTS users (
  id bigserial,
  name text NOT NULL ,
  _created_at timestamp DEFAULT now() NOT NULL ,
  PRIMARY KEY(id)
);
CREATE TABLE IF NOT EXISTS positions (
  id bigserial,
  latitude float NOT NULL ,
  longitude float NOT NULL ,
  altitude float NOT NULL ,
  floorLabel text ,
  activity text ,
  h_accuracy float NOT NULL,
  v_accuracy float NOT NULL,
  _created_at bigint ,
   user_id bigint,
  PRIMARY KEY(id),
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);