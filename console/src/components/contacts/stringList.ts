// A JSON field holding a flat list of strings (the tags field, for one).
export function isStringList(json: unknown): json is string[] {
  return Array.isArray(json) && json.every((item) => typeof item === "string");
}
