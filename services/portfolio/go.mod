module github.com/aske/go_fi_chart/services/portfolio

go 1.24.0

require (
	github.com/aske/go_fi_chart/internal/common/errors v0.0.0
	github.com/aske/go_fi_chart/pkg v0.0.0
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1
	github.com/stretchr/testify v1.10.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/aske/go_fi_chart/pkg => ../../pkg

replace github.com/aske/go_fi_chart/internal/common/errors => ../../internal/common/errors
