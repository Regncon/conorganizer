import { useCallback, useEffect, useState } from 'react';
import useEventListener from './useEventListner';

declare global {
    interface WindowEventMap {
        'local-storage': CustomEvent;
    }
}

function useReadLocalStorage<T>(key: string): T | null {
    const readValue = useCallback((): T | null => {
        if (typeof window === 'undefined') {
            return null;
        }

        try {
            const item = window.localStorage.getItem(key);
            return item ? (JSON.parse(item) as T) : null;
        } catch (error) {
            console.warn(`Error reading localStorage key “${key}”:`, error);
            return null;
        }
    }, [key]);

    const [storedValue, setStoredValue] = useState<T | null>(readValue);

    useEffect(() => {
        setStoredValue(readValue());
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    const handleStorageChange = useCallback(
        (event: StorageEvent | CustomEvent) => {
            if ((event as StorageEvent)?.key && (event as StorageEvent).key !== key) {
                return;
            }
            setStoredValue(readValue());
        },
        [key, readValue]
    );

    useEventListener('storage', handleStorageChange);

    useEventListener('local-storage', handleStorageChange);

    return storedValue;
}

export default useReadLocalStorage;
