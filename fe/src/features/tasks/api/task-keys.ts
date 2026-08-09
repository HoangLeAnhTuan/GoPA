export const taskKeys = { all: ["tasks"] as const, list: () => [...taskKeys.all, "list"] as const };
