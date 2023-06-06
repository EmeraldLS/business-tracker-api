package token

func RefreshAccessToken(refreshToken string) (string, string, int64, error, JWTClaim) {

	user, err := ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", "", 0, err, JWTClaim{}
	}
	token, newRefreshToken, expires_at, err, claims := GenerateToken(user.FirstName, user.LastName, user.Email)
	if err != nil {
		return "", "", 0, err, JWTClaim{}
	}

	return token, newRefreshToken, expires_at, nil, claims

}
