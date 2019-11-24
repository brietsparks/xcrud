This is an example CRUD data layer with supporting concerns such as testing, schema migrations, 
test data fixtures, validation, errors/logging, and a standalone CLI.  

## Database Setup
1. clone the repo:
 
    ```git clone https://github.com/brietsparks/xcrud```

2. install the binary: in the repo, run 

    ```go install```

3. setup a postgres database instance. If you already one, you can skip this. Note, this step uses the Dockerfile provided in this project, but you can connect however you want.
   
    Build the image:
    ```
    docker build ./docker -t xcrud
    ```
    
    Run the container:
    ```
    docker run --env-file ./.env.dev -p 5432:5432 xcrud
    ```
4. With a successful database connection, run migrations to create the tables:
    ```
    xcrud --env ./.env.dev migrate up   
    ```
    and
    ```
    xcrud --env ./.env.test migrate up   
    ``` 
    
5. optional: copy `.env.dev` over to `.env` to avoid having to pass in `--env` to every CLI command. 

## Usage
The data layer can be accessed via a standalone CLI or via Go code.

### CLI
Before using the CLI commands, you will need 
- a database instance with the correct tables
- a .env file that points to the database
  
See above section "Database Setup".

**Create a user:**

```
xcrud resources user:create --FirstName Bo --LastName Peep
```

Output: `{"id":1,"firstName":"Bo","lastName":"Peep"}`

**Get a user:**

```
xcrud resources user:get 1
```

Output: `{"id":1,"firstName":"Bo","lastName":"Peep"}`

**Update a user:**

```
xcrud resources user:update 1 --LastName Jackson
```

**Create a group:**

```
xcrud resources group:create --Name groupA
```

Output: `{"id":1,"name":"groupA"}`

**Add a user to a group:**

```
xcrud resources group:add-user --GroupId 1 --UserId 1
```

**Get users by group ID:**

```
xcrud resources users:get --GroupId 1
```

Output: `[{"id":1,"firstName":"Bo","lastName":"Peep"}]`

**Get groups by user ID:**
   
```
xcrud resources groups:get --UserId 1
```

Output: `[{"id":1,"name":"groupA"}]`

**Remove a user from a group:**

```
xcrud resources group:remove-user --GroupId 1 --UserId 1
```
    
### Go

1. install: ```go get -u github.com/brietsparks/xcrud```

2. see [example usage code](https://github.com/brietsparks/xcrud/blob/master/example/example.go)

## Testing
Before running the test, you will need 
- a database instance with the correct tables
- a .env file that points to the database (.env.test is provided, which points to the docker test database)
  
See above section "Database Setup".

To run the tests, from the project directory run:
```
go test ./data/tests/ --env=$(pwd)/.env.test
```
