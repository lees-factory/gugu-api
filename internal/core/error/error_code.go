package coreerror

// User (G1xxx)
const (
	G1000 = "G1000" // email already exists
)

// Auth (G2xxx)
const (
	G2000 = "G2000" // invalid credentials
	G2001 = "G2001" // email not verified
	G2002 = "G2002" // verification code not found
	G2003 = "G2003" // oauth provider invalid
	G2004 = "G2004" // refresh token invalid
)

// Product & TrackedItem (G3xxx)
const (
	G3000 = "G3000" // unsupported market
	G3001 = "G3001" // tracked item already exists
	G3002 = "G3002" // product not found
	G3003 = "G3003" // tracked item not found
)
