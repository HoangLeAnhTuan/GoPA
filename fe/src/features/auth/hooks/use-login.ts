import { useMutation } from "@tanstack/react-query";
import { login } from "../api/auth-api";
import type { LoginInput } from "../types";

export function useLogin() {
  return useMutation({ mutationFn: (input: LoginInput) => login(input) });
}
