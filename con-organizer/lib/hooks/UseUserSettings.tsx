import { useEffect, useState } from 'react';
import { doc } from 'firebase/firestore';
import { docData } from 'rxfire/firestore';
import { Subscription } from 'rxjs';
import { UserSettings } from '../../models/types';
import { db } from '../firebase';
export const userSettingsRef = (userId: string) => doc(db, `usersettings/${userId}`);
export function userSettings$(userId: string) {
    return docData(userSettingsRef(userId), { idField: 'id' });
}
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
