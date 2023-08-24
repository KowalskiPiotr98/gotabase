# gotabase
A simple database connection and migration handler for go applications.
It has 2 primary applications:
1. serving as a go-between (pun not intended) between `database/sql` and your program by providing an abstraction layer, thus somewhat decoupling from a single database handling library,
2. making connection handling easier by providing migration functionality and hiding connection management with easy access methods.

## Usage
### Connection handling

At the start of your program, run
```go
err := gotabase.InitialiseConnection(connectionString, driverName)
```
to prepare the database connection.

After that, you can access the connection pool by using the `gotabase.GetConnection()` method.

### Migrations

You can execute migrations by calling the `Migrate` method in the `migrations` package.

There are a couple of important things to remember:
1. You should store your migrations in a file with `.sql` extension. These files should be contained in a folder called `sql`. The files should have names consisting of a single number, starting at 0 and going **SEQUENTIALLY** up (as in `0.sql`, `1.sql`, `2.sql`...). If you need to remove a migration after a next number has been used, leave the file empty instead of removing it.
2. Your database model must, from the very beginning, include a migrations table (it **MUST** be created in your `0.sql` migration file). By default, it should contain a single integer column called `id` as a primary key. Note that this behaviour can be modified or adjusted to a different DBMS by overwriting the `MigrationCreator` and `LatestMigrationSelectorSql` variables of the `migrations` package.
3. The interface for providing database migrations files must be implemented, as there's no default implementation. You can, however, use the `embed.FS` struct, as it fulfills the conditions of this interface. See below for more details.

#### Usage of `embed.FS` as migration provider
The easiest way to provide migration files is to use the `embed.FS` struct.
First, create `sql` folder somewhere within your project directory structure.
Then, in the same directory as the `sql` folder, declare variable in the following way:
```go
var (
	//go:embed sql
	migrations embed.FS
)
```
This will embed the contents of the `sql` directory in your output binary file, making it easy to distribute those files.
You can also use the `migrations` variable in place of the provider interface.