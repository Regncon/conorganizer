import { useEffect, useState } from 'react';
import { Subscription } from 'rxjs';
import { ConAuthorization } from '../../models/types';
import { userSettings$ } from '../observable';

export const useUserSettings = (id?: string) => {
    const [conAuthorization, setConAuthorization] = useState<ConAuthorization>();
    const [loading, setLoading] = useState<boolean>(true);
    useEffect(() => {
        let conAuthorizationObservable: Subscription;
        if (id) {
            conAuthorizationObservable = userSettings$(id).subscribe((userSetting) => {
                setConAuthorization(userSetting as ConAuthorization);
                setLoading(false);
            });
        }
        return () => {
            if (conAuthorizationObservable?.unsubscribe) {
                conAuthorizationObservable.unsubscribe();
            }
        };
    }, [id]);

    return { conAuthorization, loading };
};
