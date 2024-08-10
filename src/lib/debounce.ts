/**
 * Debounces a function, creating a new function that does the same as the original, but will not actually run before
 * a specified amount of time has passed since it was last called.
 *
 * @param fn The function to debounce
 * @param delay Number of milliseconds to wait since the last call to the function to actually run it
 *
 * @returns A function that does the same as `fn`, but won't actually run before `delay` milliseconds has passed since
 * its last invocation. Its return value will be wrapped in a promise
 */
const debounce = <P extends unknown[], R>(
    fn: (...args: P) => R | Promise<R>,
    delay: Parameters<typeof setTimeout>[1]
): ((...args: P) => Promise<R>) => {
    let timer: ReturnType<typeof setTimeout> | null = null;

    type Reject = Parameters<ConstructorParameters<typeof Promise<R>>[0]>[1];
    let prevReject: Reject = () => { };

    return (...args: P): Promise<R> =>
        new Promise((resolve, reject) => {
            if (timer !== null) {
                clearTimeout(timer);
                prevReject('Aborted by debounce');
            }

            prevReject = reject;

            timer = setTimeout(async () => {
                timer = null;

                try {
                    const result = await fn(...args);
                    resolve(result);
                } catch (err) {
                    reject(err);
                }
            }, delay);
        });
};
export default debounce;
