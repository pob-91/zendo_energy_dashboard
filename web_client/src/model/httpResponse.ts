export type HttpResponse<T> = {
  code: number;
  success: boolean;
  data?: T;
};
