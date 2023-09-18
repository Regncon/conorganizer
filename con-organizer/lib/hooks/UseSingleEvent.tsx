import { useEffect, useState } from 'react';
import { ConEvent } from '../../models/types';
import { singleEvent$ } from '../observable';

export const useSingleEvents = (id: string) => {
    const [event, setEvent] = useState<ConEvent>();
    const [loading, setLoading] = useState<boolean>(true);

    useEffect(() => {
        const eventObservable = singleEvent$(id).subscribe((event) => {
            setEvent(event as ConEvent);
            setLoading(false);
        });

        return () => {
            eventObservable.unsubscribe();
        };
    }, []);

    return { event, loading };
};
