import { useEffect, useState } from 'react';
import { UserSettings } from '../../models/types';
import { allUserSettings$ } from '../observable';

export const useAllUserSettings = () => {
    const [allUserSettings, setUserSettings] = useState<UserSettings[]>();
    const [loadingUserSettings, setLoadingUserSettings] = useState<boolean>(true);

    useEffect(() => {
        const userSettingsObservable = allUserSettings$.subscribe((userSettings) => {
            setUserSettings(userSettings as UserSettings[] | undefined);
            setLoadingUserSettings(false);
        });

        return () => {
            userSettingsObservable.unsubscribe();
        };
    }, []);
    return { allUserSettings, loadingUserSettings };
};
