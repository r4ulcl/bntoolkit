/*

http://coussej.github.io/2015/09/15/Listening-to-generic-JSON-notifications-from-PostgreSQL-in-Go/


CREATE OR REPLACE FUNCTION notify_event() RETURNS TRIGGER AS $$

    DECLARE
        data json;
        notification json;

    BEGIN

        -- Convert the old or new row to JSON, based on the kind of action.
        -- Action = DELETE?             -> OLD row
        -- Action = INSERT or UPDATE?   -> NEW row
        IF (TG_OP = 'DELETE') THEN
            data = row_to_json(OLD);
        ELSE
            data = row_to_json(NEW);
        END IF;

        -- Contruct the notification as a JSON string.
        notification = json_build_object(
                          'table',TG_TABLE_NAME,
                          'data', data);


        -- Execute pg_notify(channel, notification)
        PERFORM pg_notify('events',notification::text);

        -- Result is ignored since this is an AFTER trigger
        RETURN NULL;
    END;

$$ LANGUAGE plpgsql;




CREATE TRIGGER download_notify_event
AFTER INSERT ON download
    FOR EACH ROW EXECUTE PROCEDURE notify_event();
*/

package utils

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
)

func waitForNotification(l *pq.Listener, verbose bool) {
	for {
		select {
		case n := <-l.Notify:
			fmt.Println(string(n.Extra))
			return
		case <-time.After(120 * time.Second):
			if verbose {
				fmt.Println("Received no events for 120 seconds, checking connection")
			}
			go func() {
				err := l.Ping()
				if err != nil {
					fmt.Println("Error conecting to DB")
					return
				}
			}()
			return
		}
	}
}

//MonitorAlert create a listener for the PostgreSQL database for alerts and print them.
func MonitorAlert(configfile string, debug bool, verbose bool) {

	config, err := GetConfig(configfile, debug, verbose)
	if err != nil {
		panic(err)
	}

	conninfo := "dbname=" + string(config.Dbname) + " user=" + string(config.User) + " password=" + string(config.Password+" sslmode=disable")

	_, err = sql.Open("postgres", conninfo)
	if err != nil {
		panic(err)
	}

	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	listener := pq.NewListener(conninfo, 10*time.Second, time.Minute, reportProblem)
	err = listener.Listen("events")
	if err != nil {
		panic(err)
	}

	fmt.Println("Start monitoring PostgreSQL...")
	for {
		waitForNotification(listener, verbose)
	}
}
