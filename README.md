dbclient demonstrates a way to make DB struct in pg package of cc-utils able to switch db endpoint dynamically. 
Main changes to DB struct:
1. Give sqlx.DB field a name so clients can't NOT call sqlx methods directly.
2. Wrap 3 methods of sqlx (BeginTxx, MustBegin and QueryRowsContext) for the demo purpose. Each method
   1) check if db is in stop mode. If true the call will return an error or panic
   2) calling sql or sqlx methods directly.
3. Add switchDB method
   1) Set DB to stop mode
   2) Closes current sqlx DB instance and creates a new one.
   3) Set DB to start mode
4. Add a goroutine to periodically trigger switching. 
  
Note: for change #4, in the real case it should periodically reads db endpoint from launchdarkly and once found any change then triggers the switch.

The demo:
1. Make one db call per second. 
2. Trigger DB switch every 5 seonds. (it mimics db switch, actually there is only one db instance so db endpoint does not change)
3. To simulate the case where DB.close gets blocked during high load, switchDB returns in 3 seconds.

```
2024/03/10 20:33:32 Committed a tx
2024/03/10 20:33:33 Query rows succeeded
2024/03/10 20:33:33 Committed a tx
2024/03/10 20:33:34 Query rows succeeded
2024/03/10 20:33:34 Committed a tx
2024/03/10 20:33:35 Query rows succeeded
2024/03/10 20:33:35 Committed a tx
2024/03/10 20:33:36 Query rows succeeded
2024/03/10 20:33:36 Committed a tx
2024/03/10 20:33:37 Time to switch db 
2024/03/10 20:33:37 switch started
2024/03/10 20:33:37 Actual closing db spent 71.748µs
2024/03/10 20:33:37 QueryxContext failed: QueryxContext failed due to dbclient in stop mode
2024/03/10 20:33:37 BeginTxx failed: Tx begin failed due to dbclient in stop mode
2024/03/10 20:33:38 QueryxContext failed: QueryxContext failed due to dbclient in stop mode
2024/03/10 20:33:38 BeginTxx failed: Tx begin failed due to dbclient in stop mode
2024/03/10 20:33:39 QueryxContext failed: QueryxContext failed due to dbclient in stop mode
2024/03/10 20:33:39 BeginTxx failed: Tx begin failed due to dbclient in stop mode
2024/03/10 20:33:40 switch done
2024/03/10 20:33:40 Query rows succeeded
2024/03/10 20:33:40 Committed a tx
2024/03/10 20:33:41 Query rows succeeded
2024/03/10 20:33:41 Committed a tx
2024/03/10 20:33:42 Time to switch db 
2024/03/10 20:33:42 switch started
2024/03/10 20:33:42 Actual closing db spent 122.686µs
2024/03/10 20:33:42 QueryxContext failed: QueryxContext failed due to dbclient in stop mode
2024/03/10 20:33:42 BeginTxx failed: Tx begin failed due to dbclient in stop mode
2024/03/10 20:33:43 QueryxContext failed: QueryxContext failed due to dbclient in stop mode
2024/03/10 20:33:43 BeginTxx failed: Tx begin failed due to dbclient in stop mode
2024/03/10 20:33:44 QueryxContext failed: QueryxContext failed due to dbclient in stop mode
2024/03/10 20:33:44 BeginTxx failed: Tx begin failed due to dbclient in stop mode

```

From the output, we could see db calls could fail when db in stop mode and later will continue when db switch is done.
