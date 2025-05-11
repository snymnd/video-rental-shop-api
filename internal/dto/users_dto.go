package dto

type (
	RegisterReq struct {
		Email    string `json:"email" binding:"email,required"`
		Name     string `json:"name" binding:"required"`
		Password string `json:"password" binding:"required,max=72"`
	}

	RegisterRes struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Role  string `json:"role"`
		Name  string `json:"name"`
	}

	LoginReq struct {
		Email    string `json:"email" binding:"email,required"`
		Password string `json:"password" binding:"required,max=72"`
	}

	LoginRes struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
		Token string `json:"token"`
	}
)
