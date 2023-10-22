import { useEffect, useState } from 'react';
import { doc } from 'firebase/firestore';
import { docData } from 'rxfire/firestore';
import { Subscription } from 'rxjs';
import { Enrollment } from '@/models/types';
import { db } from '../firebase';
export const participantEnrollmentsRef = (eventId: string, userId: string, participantId: string) =>
    doc(db, `events/${eventId}`, `/enrollments/${userId}`, `/eventParticipants/${participantId}`);

export function participantEnrollments$(eventId: string, userId: string, participantId: string) {
    return docData(participantEnrollmentsRef(eventId, userId, participantId), { idField: 'id' });
}
export const useSingleEnrollment = (eventId: string, userId?: string, participantId?: string) => {
    const [enrollments, setEnrollments] = useState<Enrollment>();
    const [loading, setLoading] = useState<boolean>(true);

    useEffect(() => {
        let enrollmentsObservable: Subscription;
        if (eventId && userId && participantId) {
            enrollmentsObservable = participantEnrollments$(eventId, userId, participantId).subscribe((enrollments) => {
                setEnrollments(enrollments as Enrollment);
                setLoading(false);
            });
        }
        return () => {
            if (enrollmentsObservable?.unsubscribe) {
                enrollmentsObservable.unsubscribe();
            }
        };
    }, [eventId, userId]);

    return { enrollments, loading };
};
