
/*
Created: 02/22/2019
Modified: 09/11/2019
Model: PostgreSQL 10
Database: PostgreSQL 10
*/


-- Create tables section -------------------------------------------------

-- Table hash

CREATE TABLE hash(
 hash Text NOT NULL,
 source Text,
 first_seen Timestamp,
 path Text,
 name Text
)
;

-- Create indexes for table hash

CREATE INDEX IX_Relationship10 ON hash (name)
;

-- Add keys for table hash

ALTER TABLE hash ADD CONSTRAINT hashKey PRIMARY KEY (hash)
;

-- Table possibles

CREATE TABLE possibles(
 id Bigint NOT NULL,
 hash Text NOT NULL,
 download Boolean,
 valid Boolean,
 possible Boolean,
 num Bigint,
 "projectName" Text
)
;

-- Create indexes for table possibles

CREATE INDEX IX_Relationship7 ON possibles ("projectName")
;

-- Add keys for table possibles

ALTER TABLE possibles ADD CONSTRAINT possibleKey PRIMARY KEY (hash)
;

-- Table project

CREATE TABLE project(
 "projectName" Text NOT NULL,
 date Timestamp
)
;

-- Add keys for table project

ALTER TABLE project ADD CONSTRAINT projectKey PRIMARY KEY ("projectName")
;

-- Table download

CREATE TABLE download(
 date Timestamp NOT NULL,
 ip Inet NOT NULL,
 port int,
 hash Text NOT NULL
)
;

-- Create indexes for table download

CREATE INDEX IX_Relationship3 ON download (ip)
;

CREATE INDEX IX_Relationship5 ON download (hash)
;

CREATE INDEX IX_Relationship6 ON download (date)
;

-- Add keys for table download

ALTER TABLE download ADD CONSTRAINT downloadKey PRIMARY KEY (ip,hash,date,port)
;

-- Table ip

CREATE TABLE ip(
 ip Inet NOT NULL,
 "projectName" Text
)
;

-- Create indexes for table ip

CREATE INDEX IX_Relationship11 ON ip ("projectName")
;

-- Add keys for table ip

ALTER TABLE ip ADD CONSTRAINT ipKey PRIMARY KEY (ip)
;

-- Table monitor

CREATE TABLE monitor(
 hash Text NOT NULL,
 userName Text,
 "projectName" Text
)
;

-- Create indexes for table monitor

CREATE INDEX IX_Relationship8 ON monitor ("projectName")
;

-- Add keys for table monitor

ALTER TABLE monitor ADD CONSTRAINT monitorKey PRIMARY KEY (hash)
;

-- Table alert

CREATE TABLE alert(
 list Text NOT NULL,
 userName Text,
 ip Inet NOT NULL,
 "projectName" Text
)
;

-- Create indexes for table alert

CREATE INDEX IX_Relationship9 ON alert ("projectName")
;

-- Add keys for table alert

ALTER TABLE alert ADD CONSTRAINT alertKey PRIMARY KEY (ip)
;
-- Create foreign keys (relationships) section ------------------------------------------------- 

ALTER TABLE download ADD CONSTRAINT Relationship3 FOREIGN KEY (ip) REFERENCES ip (ip) ON DELETE NO ACTION ON UPDATE NO ACTION
;

ALTER TABLE download ADD CONSTRAINT Relationship5 FOREIGN KEY (hash) REFERENCES hash (hash) ON DELETE NO ACTION ON UPDATE NO ACTION
;

ALTER TABLE alert ADD CONSTRAINT Relationship6 FOREIGN KEY (ip) REFERENCES ip (ip) ON DELETE NO ACTION ON UPDATE NO ACTION
;

ALTER TABLE possibles ADD CONSTRAINT Relationship7 FOREIGN KEY ("projectName") REFERENCES project ("projectName") ON DELETE NO ACTION ON UPDATE NO ACTION
;

ALTER TABLE monitor ADD CONSTRAINT Relationship8 FOREIGN KEY ("projectName") REFERENCES project ("projectName") ON DELETE NO ACTION ON UPDATE NO ACTION
;

ALTER TABLE alert ADD CONSTRAINT Relationship9 FOREIGN KEY ("projectName") REFERENCES project ("projectName") ON DELETE NO ACTION ON UPDATE NO ACTION
;

--ALTER TABLE hash ADD CONSTRAINT Relationship10 FOREIGN KEY ("projectName") REFERENCES project ("projectName") ON DELETE NO ACTION ON UPDATE NO ACTION
--;

ALTER TABLE ip ADD CONSTRAINT Relationship11 FOREIGN KEY ("projectName") REFERENCES project ("projectName") ON DELETE NO ACTION ON UPDATE NO ACTION
;


--alert

CREATE OR REPLACE FUNCTION notify_event() RETURNS TRIGGER AS $$

    DECLARE 
        data json;
        notification json;
        counter   int;
        response TEXT;
    
    BEGIN
        --check if is in alert
        select COUNT(*) into counter from alert where ip >>= NEW.ip;

        --if its in alert notify
        IF (counter != 0 ) THEN
            response = 'DATE: ' || NEW.date::text || ', IP: ' ||NEW.ip::text || ', HASH: ' || NEW.hash::text ;
            PERFORM pg_notify('events', response);
        END IF;

        RETURN NULL; 
    END;
    
$$ LANGUAGE plpgsql;






CREATE TRIGGER download_notify_event
AFTER INSERT ON download
    FOR EACH ROW EXECUTE PROCEDURE notify_event();
