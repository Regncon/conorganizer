import { useEffect, useState } from "react";
import { Subscription } from "rxjs";
import { Enrollment } from "@/models/types";
import { userEnrollments$ } from "../observable";


export const useSingleEnrollment = (eventId: string, userId?: string) => {
    const [enrollments, setEnrollments] = useState<Enrollment>();
    const [loading, setLoading] = useState<boolean>(true);

    useEffect(() => {
        let enrollmentsObservable: Subscription;
        if (eventId && userId) {
            enrollmentsObservable = userEnrollments$(eventId, userId).subscribe((enrollments) => {
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