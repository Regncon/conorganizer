import { auth,credential, firestore } from 'firebase-admin';
import { getApps, initializeApp, ServiceAccount } from 'firebase-admin/app';

console.log(process.env.FIREBASE_CLIENT_ID);

const firebaseAdminConfig = {
    type: 'service_account',
    project_id: 'regncon2023',
    private_key_id: process.env.FIREBASE_ADMIN_CLIENT_ID,
    private_key: process.env.FIREBASE_ADMIN_PRIVATE_KEY?.replace(/\\n/gm, '\n'),
    client_email: 'firebase-adminsdk-7gxfb@regncon2023.iam.gserviceaccount.com',
    client_id: process.env.FIREBASE_ADMIN_CLIENT_ID,
    auth_uri: 'https://accounts.google.com/o/oauth2/auth',
    token_uri: 'https://oauth2.googleapis.com/token',
    auth_provider_x509_cert_url: 'https://www.googleapis.com/oauth2/v1/certs',
    client_x509_cert_url:
        'https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-7gxfb%40regncon2023.iam.gserviceaccount.com',
    universe_domain: 'googleapis.com',
};

if (!getApps().length) {
    initializeApp({
        credential: credential.cert(firebaseAdminConfig as ServiceAccount),
    });
}

export const adminDb = firestore();
export const adminUser = auth();
