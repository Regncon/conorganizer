import { credential, firestore } from 'firebase-admin';
import { getApps, initializeApp, ServiceAccount } from 'firebase-admin/app';
import serviceAccount from './serviceAccountKey.json';

if (!getApps().length) {
    initializeApp({
        credential: credential.cert(serviceAccount as ServiceAccount),
    });
}

export const adminDb = firestore();
export const adminDoc = adminDb.collection('usersettings').get();
