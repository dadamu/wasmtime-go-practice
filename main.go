package main

import (
	"encoding/binary"
	"fmt"
	"log"

	"github.com/bytecodealliance/wasmtime-go/v24"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sqlx.Connect("postgres", "user=test dbname=test password=test sslmode=disable")
	check(err)

	// Set up Wasmtime Engine and Store
	engine := wasmtime.NewEngine()
	store := wasmtime.NewStore(engine)

	// Define the SQL query function
	sqlQuery :=
		func(caller *wasmtime.Caller, sqlPtr int32) {
			// Access the WASM memory
			memory := caller.GetExport("memory").Memory()

			// Read the SQL string from WASM memory
			data := memory.UnsafeData(store)

			// Find the null terminator to determine the length of the string
			var sqlBytes []byte
			// Read bytes until we find the null terminator (0, 0 in UTF-16)
			for i := sqlPtr; binary.LittleEndian.Uint16(data[i:i+2]) != 0; i += 2 {
				// Append the first byte of the UTF-16 character to sqlBytes
				sqlBytes = append(sqlBytes, data[i])
			}
			sql := string(sqlBytes)

			// Execute the SQL query
			fmt.Printf("Executing SQL: %s\n", sql)
			_, err := db.Exec(sql)
			if err != nil {
				log.Printf("SQL execution error: %v\n", err)
			}
		}

	abort := func(arg1, arg2, arg3, arg4 int32) {
		fmt.Printf("WASM called abort with args: %d, %d, %d, %d\n", arg1, arg2, arg3, arg4)
	}

	// Set up the imports mapping
	linker := wasmtime.NewLinker(engine)
	linker.DefineFunc(store, "index", "query", sqlQuery)
	linker.DefineFunc(store, "env", "abort", abort)

	// Load your WASM module
	module, err := wasmtime.NewModuleFromFile(engine, "release.wasm")
	if err != nil {
		panic(err)
	}

	instance, err := linker.Instantiate(store, module)
	check(err)

	init := instance.GetFunc(store, "init")
	if init == nil {
		panic("init function not found")
	}
	_, err = init.Call(store)
	check(err)

	run := instance.GetFunc(store, "run")
	if run == nil {
		panic("run function not found")
	}

	_, err = run.Call(store)
	check(err)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
