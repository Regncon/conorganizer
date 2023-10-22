import { useEffect, useState } from 'react';
import { collection } from 'firebase/firestore';
import { collectionData } from 'rxfire/firestore';
import { UserSettings } from '../../models/types';
import { db } from '../firebase';
export const allUserSettingsRef = collection(db, 'usersettings');
export const allUserSettings$ = collectionData(allUserSettingsRef, { idField: 'id' });
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
