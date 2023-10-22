import { useEffect, useState } from 'react';
import { collection } from 'firebase/firestore';
import { collectionData } from 'rxfire/firestore';
import { Subscription } from 'rxjs';
import { EnrollmentChoice } from '@/models/types';
import { db } from '../firebase';
export const enrollmentChoicesRef = (eventId: string) => collection(db, `events/${eventId}/enrollmentChoices/`);
export function allEnrollmentChoices$(eventId: string) {
    return collectionData(enrollmentChoicesRef(eventId), { idField: 'id' });
}
export const useAllEnrollmentChoices = (eventId?: string) => {
    const [enrollmentChoices, setEnrollmentChoices] = useState<EnrollmentChoice[]>();
    const [loadingEnrollmentChoices, setLoadingEnrollmentChoices] = useState<boolean>(true);

    useEffect(() => {
        let allEnrollmentChoicesObservable: Subscription;
        if (eventId)
            allEnrollmentChoicesObservable = allEnrollmentChoices$(eventId).subscribe((enrollmentChoices) => {
                setEnrollmentChoices(enrollmentChoices as EnrollmentChoice[] | undefined);
                setLoadingEnrollmentChoices(false);

                return () => {
                    allEnrollmentChoicesObservable.unsubscribe();
                };
            });
        return () => {
            null;
        };
    }, []);
    return { enrollmentChoices, loadingEnrollmentChoices };
};
