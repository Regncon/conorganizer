import { useEffect, useState } from 'react';
import { doc } from 'firebase/firestore';
import { docData } from 'rxfire/firestore';
import { ConEvent } from '../../models/types';
import { db } from '../firebase';
export const eventRef = (id: string) => doc(db, `events/${id}`);
export function singleEvent$(id: string) {
    return docData(eventRef(id), { idField: 'id' });
}
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
