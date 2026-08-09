export interface UserDto {
  id: string;
  email: string;
  role: "USER" | "ADMIN";
  created_at: string;
}

export interface User {
  id: string;
  email: string;
  role: UserDto["role"];
  createdAt: Date;
}

export interface AuthDto {
  user: UserDto;
  access_token: string;
}

export interface AuthResult {
  user: User;
  accessToken: string;
}

export interface LoginInput {
  email: string;
  password: string;
}

export interface RegisterInput extends LoginInput {
  confirmPassword: string;
}
