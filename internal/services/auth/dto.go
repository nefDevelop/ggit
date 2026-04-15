package auth

// LoginOptions contiene los datos para iniciar sesión.
type LoginOptions struct {
	Username    string
	Password    string
	LoginSource int64
	Remember    bool
}

// SignupOptions contiene los datos para registrar un nuevo usuario.
type SignupOptions struct {
	Username string
	Email    string
	Password string
}

// ResetPasswordOptions contiene los datos para restablecer la contraseña.
type ResetPasswordOptions struct {
	Code        string
	NewPassword string
}
