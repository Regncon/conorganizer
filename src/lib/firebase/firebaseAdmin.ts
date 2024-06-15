import { credential } from 'firebase-admin';
import { getApps, initializeApp as initializeAdminApp, ServiceAccount, type App } from 'firebase-admin/app';
import { getAuth as getAdminAuth } from 'firebase-admin/auth';
import { cookies } from 'next/headers';
import { getAuth, signInWithCustomToken } from 'firebase/auth';
import { initializeApp } from 'firebase/app';
import { firebaseAdminConfig, firebaseConfig } from './config';
import { getFirestore as getAdminFirestore } from 'firebase-admin/firestore';
import { getFirestore } from 'firebase/firestore';

export const SESSION_COOKIE_NAME = '__session';
export const getAuthorizedAuth = async () => {
    let session: string | undefined;
    const adminApp =
        getApps().find((app) => app.name === 'admin') ||
        initializeAdminApp(
            {
                credential: credential.cert(firebaseAdminConfig as ServiceAccount),
            },
            'admin'
        );
    const adminAuth = getAdminAuth(adminApp);
    const noSessionReturn = { app: null, user: null, auth: null, db: null };

    if (!session) {
        session = await getAppRouterSession();

        if (!session) {
            return noSessionReturn;
        }
    }

    try {
        //Litt usikker på navn, men dette er noe annet en adminApp dette er en "app, auth og db for en logget inn  bruker som kan brukes server components"
        const decodedIdToken = await adminAuth.verifySessionCookie(session);
        const app = initializeAuthenticatedApp(decodedIdToken.uid);
        const auth = getAuth(app);
        const db = getFirestore(app);

        const isRevoked = !(await adminAuth.verifySessionCookie(session, true).catch((e) => console.error(e.message)));
        if (isRevoked) return noSessionReturn;
        // To signIn user so we get access to auth.currentUser
        if (auth.currentUser?.uid !== decodedIdToken.uid) {
            const customToken = await adminAuth
                .createCustomToken(decodedIdToken.uid)
                .catch((e) => console.error(e.message));

            if (!customToken) return noSessionReturn;

            await signInWithCustomToken(auth, customToken);
        }

        return { app, user: auth.currentUser, auth, db };
    } catch (error) {
        return noSessionReturn;
    }
};

async function getAppRouterSession() {
    const cookieStore = cookies();

    try {
        return cookieStore.get(SESSION_COOKIE_NAME)?.value;
    } catch (error) {
        return undefined;
    }
}

function initializeAuthenticatedApp(uid: string) {
    const random = Math.random().toString(36).split('.')[1];
    const appName = `authenticated-context:${uid}:${random}`;

    const app = initializeApp(firebaseConfig, appName);

    return app;
}

let app: App =
    getApps().find((app) => app.name === 'admin') ||
    initializeAdminApp(
        {
            credential: credential.cert(firebaseAdminConfig as ServiceAccount),
        },
        'admin'
    );

// if (!getApps().length && !getApps().some((app) => app.name.includes('[DEFAULT]'))) {
// 	app = initializeAdminApp({
// 		credential: credential.cert(firebaseAdminConfig as ServiceAccount),
// 	});
// }

export const adminDb = getAdminFirestore(app);
export const adminAuth = getAdminAuth(app);

// const adminApp =
// 	getApps().find((app) => app.name === 'admin') ||
// 	initializeApp(
// 		{
// 			credential: credential.cert(firebaseAdminConfig as ServiceAccount),
// 		},
// 		'admin'
// 	);
// export const adminDb = firestore(adminApp);
