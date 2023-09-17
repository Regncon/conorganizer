import { useEffect, useState } from 'react';
import { ConEvent } from '../../models/types';
import { allEvents$ } from '../observables/AllEvents';

export const useAllEvents = () => {
    const [events, setEvents] = useState<ConEvent[]>();
    const [loading, setLoading] = useState<boolean>(true);

    useEffect(() => {
        const eventsObservable = allEvents$.subscribe((events) => {
            if (events) {
                setEvents(events as ConEvent[]);
                setLoading(false);
            }
        });

        return () => {
            eventsObservable.unsubscribe();
        };
    }, []);
    return { events, loading };
};
