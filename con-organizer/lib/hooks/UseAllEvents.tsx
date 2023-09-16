import { useEffect, useState } from 'react';
import { collection, query } from 'firebase/firestore';
import { collectionData } from 'rxfire/firestore';
import { tap } from 'rxjs';
import db from '@/lib/firebase';
import { ConEvent } from '../types';

export const useAllEvents = () => {
    const [events, setEvents] = useState<ConEvent[]>();
    const [loading, setLoading] = useState<boolean>(true);
    const eventRef = collection(db, 'events');

    useEffect(() => {
        console.log('in useeffect');

        const conEvents = collectionData<ConEvent[]>(eventRef)
            .pipe(tap((events) => console.log('This is just an observable!')))
            .subscribe((events) => {
                console.log(events);

                setEvents(events);
                setLoading(false);
            });

        console.log(conEvents);

        return () => {
            conEvents.unsubscribe();
        };
    }, []);
    return { events, loading };
};
