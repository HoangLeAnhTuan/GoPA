export interface UserDto {
  id: string;
  email: string;
  role: "USER" | "ADMIN";
  display_name?: string;
  avatar_url?: string;
  created_at: string;
}

export interface User {
  id: string;
  email: string;
  role: UserDto["role"];
  displayName?: string;
  avatarUrl?: string;
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
  displayName?: string;
}

export interface UpdateProfileInput {
  display_name?: string;
  avatar_url?: string;
}
