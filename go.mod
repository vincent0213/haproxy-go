module forward

go 1.16

require (
	github.com/juju/ratelimit v1.0.1
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	tool v0.0.0-00010101000000-000000000000
)

//replace bridge => ../bridge
replace tool => ../tool
