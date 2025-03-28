// https://gist.github.com/t3dotgg/a486c4ae66d32bf17c09c73609dacc5b?permalink_comment_id=5503910#gistcomment-5503910
export type PromiseResult<TData, TError = Error> =
  | [TError, null]
  | [null, TData];

export async function promiseResult<TData, TError = Error>(
  promise: Promise<TData>
): Promise<PromiseResult<TData, TError>> {
  try {
    const result = await promise;
    return [null, result];
  } catch (error) {
    return [
      error instanceof Error
        ? (error as TError)
        : (new Error(String(error)) as TError),
      null,
    ];
  }
}
