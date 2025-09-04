<h1>MySQL Align</h1>

This is MySQL script for algin information between source and destination. This is script be able to open connection to remote MySQL server, help minimize access to sesitive production server.

<h2>Features</h2>

1. Reconcile total rows of table for every schema
2. Reconcile every object within database
3. Check compatibility before migration such as CHAR_SET, ENGINE, ROW_FORMAT, PRIMARY KEY, and so on.
4. Get all configurations of MySQL Server
5. Get size of each schema in database


<h2>Usage</h2>
This script run with:

```bash
my-align <command> [options]
```
