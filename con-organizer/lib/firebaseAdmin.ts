import { credential, firestore, initializeApp } from 'firebase-admin';
import { getApps, ServiceAccount } from 'firebase-admin/app';
import serviceAccount from './serviceAccountKey.json';

if (!getApps().length) {
    initializeApp({
        credential: credential.cert(serviceAccount as ServiceAccount),
    });
}

export const adminDb = firestore();
export const adminDoc = adminDb.collection('usersettings').get();
