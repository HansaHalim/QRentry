CREATE DATABASE qrEntryDatabase;

CREATE TABLE GATE (
  gateID INT NOT NULL,
  gateName TEXT,
  PRIMARY KEY (gateID)
)

CREATE TABLE AUTHENTICATED_USERS (
  gateID INT NOT NULL,
  email TEXT NOT NULL,
  FOREIGN KEY (gateID) REFERENCES GATE(gateID)
);