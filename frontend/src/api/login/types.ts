export interface UserLoginType {
  username: string
  password: string
}

export interface UserType {
  username: string
  password: string
  role: string
  roleId: string
}

export interface Token {
  access_token: string
}

export interface changePassword {
  newPassword: string
}

export interface changePasswordResponse {
  message: string
}
