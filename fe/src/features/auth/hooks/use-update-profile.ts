import { useMutation, useQueryClient } from "@tanstack/react-query";
import { updateProfile } from "../api/auth-api";
import type { UpdateProfileInput } from "../types";

export function useUpdateProfile() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (input: UpdateProfileInput) => updateProfile(input),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["auth", "me"] });
    },
  });
}
