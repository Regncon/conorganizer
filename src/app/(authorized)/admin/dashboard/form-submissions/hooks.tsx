import { useEffect, useState } from 'react';
import type { FormSubmission } from './types';
import { collectionGroup, onSnapshot, query } from 'firebase/firestore';
import { db } from '$lib/firebase/firebase';
import type { MyNewEvent } from '$lib/types';

export const useRealtimeTableData = () => {
    const [allSubmissions, setAllSubmissions] = useState<FormSubmission[]>([]);
    useEffect(() => {
        const myEventsQuery = query(collectionGroup(db, 'my-events'));

        const unsubscribe = onSnapshot(myEventsQuery, (snapshot) => {
            const getPriority = (row: FormSubmission) => {
                const unreadAndNotAccepted = row.isRead === false && row.isAccepted === false;
                const readAndNotAccepted = row.isRead === true && row.isAccepted === false;
                const readAndAccepted = row.isRead === true && row.isAccepted === true;
                const defaultAndUnexpectedValue = 4;
                if (unreadAndNotAccepted) return 1;
                if (readAndNotAccepted) return 2;
                if (readAndAccepted) return 3;
                return defaultAndUnexpectedValue;
            };
            const AllSubmissions = snapshot.docs
                .map((doc) => {
                    const data = doc.data() as MyNewEvent;
                    const submissions: FormSubmission = {
                        id: doc.id,
                        title: data.title,
                        subTitle: data.subTitle,
                        isRead: data.isRead ?? false,
                        isAccepted: data.isAccepted ?? false,
                    };

                    return submissions;
                })
                .sort((a, b) => {
                    const priorityA = getPriority(a);
                    const priorityB = getPriority(b);
                    return priorityA - priorityB;
                });
            setAllSubmissions(AllSubmissions);
        });

        return () => {
            unsubscribe();
        };
    }, []);
    return allSubmissions;
};
