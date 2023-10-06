import { useEffect, useState } from "react";
import { Subscription } from "rxjs";
import { Enrollment } from "@/models/types";
import { participantEnrollments$ } from "../observable";


export const useSingleEnrollment = (eventId: string, userId?: string, participantId?: string ) => {
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