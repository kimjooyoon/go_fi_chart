module github.com/aske/go_fi_chart/services/portfolio

go 1.24.0

require (
	github.com/aske/go_fi_chart/pkg v0.0.0
	github.com/go-chi/chi/v5 v5.2.1
	github.com/google/uuid v1.6.0
)

replace github.com/aske/go_fi_chart/pkg => ../../pkg
