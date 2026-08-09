import { useMutation } from "@tanstack/react-query";
import { register } from "../api/auth-api";
import type { RegisterInput } from "../types";

export function useRegister() {
  return useMutation({ mutationFn: (input: RegisterInput) => register(input) });
}
