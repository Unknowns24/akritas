export interface ServiceData<T> {
  data: T;
}

export function requireApiData<T>(data: T | undefined, error: unknown): T {
  if (error) {
    throw error;
  }

  if (data === undefined || data === null) {
    throw new Error("No data returned");
  }

  return data;
}
