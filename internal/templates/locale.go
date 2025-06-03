package templates

// LocaleEnTemplate is the English localization template
const LocaleEnTemplate = `{
  "validation": {
    "required": "This field is required",
    "email": "Please enter a valid email address",
    "min_length": "Must be at least %d characters",
    "max_length": "Must not exceed %d characters"
  },
  "auth": {
    "login_success": "Login successful",
    "login_failed": "Invalid email or password",
    "register_success": "Registration successful",
    "register_failed": "Registration failed",
    "token_invalid": "Invalid token",
    "token_expired": "Token expired",
    "logout_success": "Logout successful"
  },
  "general": {
    "created": "Created successfully",
    "updated": "Updated successfully", 
    "deleted": "Deleted successfully",
    "retrieved": "Retrieved successfully",
    "not_found": "Resource not found",
    "internal_error": "Internal server error",
    "bad_request": "Bad request",
    "unauthorized": "Unauthorized",
    "forbidden": "Forbidden"
  }
}
`

// LocaleIdTemplate is the Indonesian localization template
const LocaleIdTemplate = `{
  "validation": {
    "required": "Field ini wajib diisi",
    "email": "Silakan masukkan alamat email yang valid",
    "min_length": "Minimal %d karakter",
    "max_length": "Maksimal %d karakter"
  },
  "auth": {
    "login_success": "Login berhasil",
    "login_failed": "Email atau password salah",
    "register_success": "Registrasi berhasil",
    "register_failed": "Registrasi gagal",
    "token_invalid": "Token tidak valid",
    "token_expired": "Token kedaluwarsa",
    "logout_success": "Logout berhasil"
  },
  "general": {
    "created": "Berhasil dibuat",
    "updated": "Berhasil diperbarui",
    "deleted": "Berhasil dihapus", 
    "retrieved": "Berhasil diambil",
    "not_found": "Resource tidak ditemukan",
    "internal_error": "Terjadi kesalahan server",
    "bad_request": "Permintaan tidak valid",
    "unauthorized": "Tidak terotorisasi",
    "forbidden": "Dilarang"
  }
}
`
