// The entry file of your WebAssembly module.
declare function query(sql: string): void;

export function init(): void {
  query("CREATE TABLE IF NOT EXISTS person (id SERIAL PRIMARY KEY, name TEXT)");
}

export function run(): void {
  query("INSERT INTO person (name) VALUES ('Alice')");
  query("INSERT INTO person (name) VALUES ('Bob')");
  query("INSERT INTO person (name) VALUES ('Charlie')");
}
