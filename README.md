### testing on Ubuntu using CockroachDB and Golang example code

This is an attempt to implement the Go tutorial on a CockroachDB Cluster using pgx driver

#### setup environment for CockroachDB single node cluster

using multipass create Ubuntu VM

```
multipass launch --name ubuntu --cpus 4 --mem 8G --disk 40G
```

attach persistent mount

```
multipass mount /home/admin/Documents/multipass_env/Self-Hosted_CockroachDB_Golang cockroach ubuntu:/cockroach
```

enter ubuntu shell

```
multipass shell ubuntu
```

get go

```
sudo apt-get update
sudo apt install golang-go -y
```

install cockroachdb

```
curl https://binaries.cockroachdb.com/cockroach-v22.1.2.linux-amd64.tgz | tar -xz && sudo cp -i cockroach-v22.1.2.linux-amd64/cockroach /usr/local/bin/

sudo mkdir -p /usr/local/lib/cockroach

sudo cp -i cockroach-v22.1.2.linux-amd64/lib/libgeos.so /usr/local/lib/cockroach/

sudo cp -i cockroach-v22.1.2.linux-amd64/lib/libgeos_c.so /usr/local/lib/cockroach/

which cockroach
/usr/local/bin/cockroach
```

#### Setup CockroachDB cluster and Golang environment

CockroachDB documentation for single node cluster setup: https://www.cockroachlabs.com/docs/stable/cockroach-start-single-node.html?filters=insecure

In my case, as I am running Ubuntu in a virtualized environment, if I expose the ports on localhost then this will not allow me to connect from my virtualization host where I have access to the GUI. For this reason, I use a bind command with the IP of the VM instances

```
cockroach start-single-node \
--insecure \
--listen-addr=10.166.133.153:26257 \
--http-addr=10.166.133.153:8080 \
--background
```

Once you have used your IP to bind the process you'll need to reference it again in your conection string

```
cockroach sql --insecure --host=localhost:26257
```

This could be helpful for awareness of how Go approaches this
Go example with a database: https://go.dev/doc/tutorial/database-access

Even more relevant - Go drivers CockroachDB documentation:https://www.cockroachlabs.com/docs/stable/connect-to-the-database.html?filters=core&filters=go

First thing you'll want to do is ensure you have your connection string configured, then you can export as a variable, for example

```
export DATABASE_URL="postgresql://root@10.166.133.153:26257/test?sslmode=disable"
```

Once your environment is setup you can then setup the database.

#### Create table in CockroachDB  

```
# create a database first
CREATE DATABASE album;
USE album;


CREATE TABLE album (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  title      VARCHAR(128) NOT NULL,
  artist     VARCHAR(255) NOT NULL,
  price      DECIMAL(5,2) NOT NULL
);

INSERT INTO album
  (title, artist, price)
VALUES
  ('Blue Train', 'John Coltrane', 56.99),
  ('Giant Steps', 'John Coltrane', 63.99),
  ('Jeru', 'Gerry Mulligan', 17.99),
  ('Sarah Vaughan', 'Sarah Vaughan', 34.98);
```

#### Golang example application

First create a test database

```
CREATE DATABASE test;
USE test;
```

This connects to the database

```
package main

import (
    "context"
    "log"

    "github.com/jackc/pgx/v4"
)

func main() {
    conn, err := pgx.Connect(context.Background(), "<connection-string>")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close(context.Background())
}
```

From quickstart, this is a Hello World example

```
package main

import (
        "context"
        "log"
        "os"

        "github.com/jackc/pgx/v4"
)

func execute(conn *pgx.Conn, stmt string) error {
        rows, err := conn.Query(context.Background(), stmt)
        if err != nil {
                log.Fatal(err)
        }
        for rows.Next() {
                var message string
                if err := rows.Scan(&message); err != nil {
                        log.Fatal(err)
                }
                log.Printf(message)
        }
        return nil
}

func main() {
        // Read in connection string
        config, err := pgx.ParseConfig(os.Getenv("DATABASE_URL"))
        if err != nil {
                log.Fatal(err)
        }
        config.RuntimeParams["application_name"] = "$ docs_quickstart_go"
        connection, err := pgx.ConnectConfig(context.Background(), config)
        if err != nil {
                log.Fatal(err)
        }
        defer connection.Close(context.Background())

        statements := [3]string{
                // CREATE the messages table
                "CREATE TABLE IF NOT EXISTS messages (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), message STRING)",
                // INSERT a row into the messages table
                "INSERT INTO messages (message) VALUES ('Hello world!')",
                // SELECT a row from the messages table
                "SELECT message FROM messages"}

        for _, s := range statements {
                execute(connection, s)
        }
}
```

May want to experiment with adding rows like this Golang example: https://github.com/cockroachlabs-field/stock-data-gen/blob/main/gen_data.go
