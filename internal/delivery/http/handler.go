package handler

// echojwt "github.com/labstack/echo-jwt/v4"

type response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omiempty"`
}

// type jwtCustomClaims struct {
// 	Id       string `json:"id"`
// 	Username string `json:"username"`
// 	Role     string `json:"role"`
// 	jwt.RegisteredClaims
// }

// func jwtConfig() echojwt.Config {
// 	return echojwt.Config{
// 		NewClaimsFunc: func(c echo.Context) jwt.Claims {
// 			return new(jwtCustomClaims)
// 		},
// 		SigningKey: []byte(config.GetJwtSecret()),
// 	}
// }

// func signJwtToken(id int, username string, role string) (string, error) {
// 	claims := &jwtCustomClaims{
// 		id,
// 		username,
// 		role,
// 		jwt.RegisteredClaims{},
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

// 	t, err := token.SignedString([]byte(config.GetJwtSecret()))
// 	if err != nil {
// 		return "", err
// 	}

// 	return t, nil
// }

// func claimsSession(c echo.Context) *jwtCustomClaims {
// 	user := c.Get("user").(*jwt.Token)
// 	claims := user.Claims.(*jwtCustomClaims)
// 	return claims
// }
