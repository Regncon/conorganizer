import { useEffect, useState } from 'react';
import { ConAuthorization } from '../../models/types';
import { conAuthorizationRef$ } from '../observables/authorization';

export const useAuthorizationHook = (id: string) => {
    const [conAuthorization, setConAuthorization] = useState<ConAuthorization>();
    const [loading, setLoading] = useState<boolean>(true);

    useEffect(() => {
        const conAuthorizationObservable = conAuthorizationRef$(id).subscribe((snapshot) => {
            if (snapshot.data()) {
                setConAuthorization({ ...(snapshot.data() as ConAuthorization), id: snapshot.id });
                setLoading(false);
            }
            setLoading(false);
        });

        return () => {
            conAuthorizationObservable.unsubscribe();
        };
    }, []);

    return { conAuthorization, loading };
};
