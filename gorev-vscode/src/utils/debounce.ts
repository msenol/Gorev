/**
 * Debouncing utility for controlling function execution frequency
 * Prevents rapid successive calls to expensive operations
 */

export interface DebounceOptions {
    delay?: number;
    immediate?: boolean;
    maxWait?: number;
}

export interface DebouncedFunction<T extends (...args: any[]) => any> {
    (...args: Parameters<T>): Promise<ReturnType<T>>;
    cancel(): void;
    flush(): Promise<ReturnType<T> | undefined>;
    pending(): boolean;
}

/**
 * Creates a debounced version of the provided function
 * @param func Function to debounce
 * @param options Debounce configuration
 */
export function debounce<T extends (...args: any[]) => any>(
    func: T,
    options: DebounceOptions = {}
): DebouncedFunction<T> {
    const { delay = 500, immediate = false, maxWait } = options;

    let timeoutId: NodeJS.Timeout | undefined;
    let maxTimeoutId: NodeJS.Timeout | undefined;
    let lastCallTime: number | undefined;
    let lastArgs: Parameters<T> | undefined;
    let result: ReturnType<T> | undefined;
    let isInvoking = false;

    // Promise tracking for async support
    let resolvePromise: ((value: ReturnType<T>) => void) | undefined;
    let rejectPromise: ((reason: any) => void) | undefined;
    let activePromise: Promise<ReturnType<T>> | undefined;

    function shouldInvoke(time: number): boolean {
        if (lastCallTime === undefined) return true;
        if (time - lastCallTime >= delay) return true;
        if (maxWait !== undefined && time - lastCallTime >= maxWait) return true;
        return false;
    }

    function invokeFunc(): ReturnType<T> {
        const args = lastArgs!;
        lastArgs = undefined;
        isInvoking = true;

        try {
            result = func.apply(null, args) as ReturnType<T>;
            isInvoking = false;
            return result!;
        } catch (error) {
            isInvoking = false;
            throw error;
        }
    }

    function leadingEdge(): ReturnType<T> | undefined {
        lastCallTime = Date.now();

        if (maxWait !== undefined) {
            maxTimeoutId = setTimeout(timerExpired, maxWait);
        }

        return immediate ? invokeFunc() : result;
    }

    function timerExpired(): void {
        const time = Date.now();
        if (shouldInvoke(time)) {
            trailingEdge();
        } else {
            timeoutId = setTimeout(timerExpired, delay - (time - lastCallTime!));
        }
    }

    function trailingEdge(): void {
        timeoutId = undefined;

        if (lastArgs) {
            try {
                const invokeResult = invokeFunc();
                if (resolvePromise) {
                    resolvePromise(invokeResult);
                    resolvePromise = undefined;
                    rejectPromise = undefined;
                    activePromise = undefined;
                }
            } catch (error) {
                if (rejectPromise) {
                    rejectPromise(error);
                    resolvePromise = undefined;
                    rejectPromise = undefined;
                    activePromise = undefined;
                }
            }
        }

        lastCallTime = undefined;
        if (maxTimeoutId) {
            clearTimeout(maxTimeoutId);
            maxTimeoutId = undefined;
        }
    }

    function cancel(): void {
        if (timeoutId) {
            clearTimeout(timeoutId);
            timeoutId = undefined;
        }
        if (maxTimeoutId) {
            clearTimeout(maxTimeoutId);
            maxTimeoutId = undefined;
        }

        if (rejectPromise) {
            rejectPromise(new Error('Debounced function cancelled'));
            resolvePromise = undefined;
            rejectPromise = undefined;
            activePromise = undefined;
        }

        lastCallTime = undefined;
        lastArgs = undefined;
        isInvoking = false;
    }

    async function flush(): Promise<ReturnType<T> | undefined> {
        if (timeoutId) {
            clearTimeout(timeoutId);
            timeoutId = undefined;

            if (lastArgs) {
                try {
                    return invokeFunc();
                } catch (error) {
                    throw error;
                }
            }
        }
        return result;
    }

    function pending(): boolean {
        return timeoutId !== undefined || isInvoking;
    }

    function debounced(...args: Parameters<T>): Promise<ReturnType<T>> {
        const time = Date.now();
        const isInvokeNow = shouldInvoke(time);

        lastArgs = args;
        lastCallTime = time;

        // Return existing promise if one is active
        if (activePromise && !isInvokeNow) {
            return activePromise;
        }

        // Create new promise for this invocation
        activePromise = new Promise<ReturnType<T>>((resolve, reject) => {
            resolvePromise = resolve;
            rejectPromise = reject;
        });

        if (isInvokeNow) {
            if (timeoutId === undefined) {
                try {
                    const invokeResult = leadingEdge();
                    if (!immediate && resolvePromise) {
                        // For non-immediate mode, we need to wait for trailing edge
                        timeoutId = setTimeout(timerExpired, delay);
                    } else if (immediate) {
                        // Immediate mode - resolve right away
                        if (resolvePromise && invokeResult !== undefined) {
                            resolvePromise(invokeResult);
                            resolvePromise = undefined;
                            rejectPromise = undefined;
                            activePromise = undefined;
                        }
                    }
                } catch (error) {
                    if (rejectPromise) {
                        rejectPromise(error);
                        resolvePromise = undefined;
                        rejectPromise = undefined;
                        activePromise = undefined;
                    }
                }
            }
        } else {
            timeoutId = setTimeout(timerExpired, delay);
        }

        return activePromise!;
    }

    debounced.cancel = cancel;
    debounced.flush = flush;
    debounced.pending = pending;

    return debounced;
}

/**
 * Specialized debounce for refresh operations
 * Optimized for tree view refreshes with sensible defaults
 */
export function debounceRefresh<T extends (...args: any[]) => any>(
    func: T,
    delay = 500
): DebouncedFunction<T> {
    return debounce(func, {
        delay,
        immediate: false,
        maxWait: delay * 3 // Ensure execution within 3x delay at most
    });
}

/**
 * Specialized debounce for configuration changes
 * Uses immediate execution to provide responsive UI feedback
 */
export function debounceConfig<T extends (...args: any[]) => any>(
    func: T,
    delay = 100
): DebouncedFunction<T> {
    return debounce(func, {
        delay,
        immediate: true,
        maxWait: delay * 2
    });
}