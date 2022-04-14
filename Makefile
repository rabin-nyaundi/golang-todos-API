#! run the application
RUN: 
	go run ./cmd/api


#? create table migrations
CREATE_TODOS_TABLE_MIGRATIONS:
	migrate create -seq -ext=.sql -dir=./migrations create_todos_table


#? adding table constraints
ADD_TODOS_CHECK_CONSTRAINTS:
	# migrate create -seq -ext=.sql -dir=./migrations add_todos_check_constraints

# ? Run migrations
EXECUTE_UP_MIGRATIONS:
	migrate -path=./migrations -database=$$TODO_DB_DSN up

EXECUTE_DOWN_MIGRATIONS:
	migrate -path=./migrations -database=$$TODO_DB_DSN down

ADD_TODOS_INDEXES:
	migrate create -seq -ext .sql -dir=./migrations add_todos_indexes