import { useEffect, useState } from 'react';
import { collection } from 'firebase/firestore';
import { collectionData } from 'rxfire/firestore';
import { ConEvent } from '../../models/types';
import { db } from '../firebase';
export const eventsRef = collection(db, 'events');
export const allEvents$ = collectionData(eventsRef, { idField: 'id' });
export const useAllEvents = () => {
    const [events, setEvents] = useState<ConEvent[]>();
    const [loading, setLoading] = useState<boolean>(true);

    useEffect(() => {
        const eventsObservable = allEvents$.subscribe((events) => {
            setEvents(events as ConEvent[] | undefined);
            setLoading(false);
        });

        return () => {
            eventsObservable.unsubscribe();
        };
    }, []);

    return { events, loading };
};
