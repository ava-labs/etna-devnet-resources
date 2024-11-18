import { useState, useCallback } from 'react';

// Custom hook to handle async operations
export function useAsync<T>(asyncFn: () => Promise<T>) {
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [data, setData] = useState<T | null>(null);

    const execute = useCallback(async () => {
        try {
            setLoading(true);
            setError(null);
            const result = await asyncFn();
            setData(result);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'An error occurred');
        } finally {
            setLoading(false);
        }
    }, [asyncFn]);

    return { loading, error, data, execute };
}
