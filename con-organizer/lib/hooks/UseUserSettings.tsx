import { useEffect, useState } from 'react';
import { Subscription } from 'rxjs';
import { UserSettings } from '../../models/types';
import { userSettings$ } from '../observable';

export const useUserSettings = (id?: string) => {
    const [userSetting, setUserSetting] = useState<UserSettings>();
    const [loading, setLoading] = useState<boolean>(true);
    useEffect(() => {
        let userSettingsObservable: Subscription;
        if (id) {
            userSettingsObservable = userSettings$(id).subscribe((userSetting) => {
                setUserSetting(userSetting as UserSettings);
                setLoading(false);
            });
        }
        return () => {
            if (userSettingsObservable?.unsubscribe) {
                userSettingsObservable.unsubscribe();
            }
        };
    }, [id]);

    return { userSettings: userSetting, loading };
};
