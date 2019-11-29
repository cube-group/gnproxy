package etcd

const (
	FRONTEND           = "frontend"
	F_PASS_HOST_HEADER = "passhostheader"

	F_AUTH                     = "auth"
	F_AUTH_FORWARD             = "forward"
	F_AUTH_FORWARD_RES_HEADERS = ".authresponseheaders"
	F_AUTH_FORWARD_ADDRESS     = "address"

	F_ROUTES               = "routes"
	F_ROUTES_SERVICES      = "services"
	F_ROUTES_SERVICES_RULE = "rule"

	BACKEND = "backend"
	B_CIRCUIBREAKER = "circuitbreaker/expression"
)
