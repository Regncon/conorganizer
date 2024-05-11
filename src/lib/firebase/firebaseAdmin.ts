import { credential, firestore } from 'firebase-admin';
import {
	getApps,
	initializeApp as initializeAdminApp,
	refreshToken,
	ServiceAccount,
	type App,
} from 'firebase-admin/app';
import { getAuth as getAdminAuth } from 'firebase-admin/auth';

import { cookies } from 'next/headers';

import { getAuth, signInWithCustomToken } from 'firebase/auth';
import { initializeApp, type FirebaseOptions } from 'firebase/app';
import { firebaseAdminConfig, firebaseConfig } from './config';
import { getFirestore } from 'firebase-admin/firestore';

//console.log(process.env.FIREBASE_CLIENT_ID);

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
	const noSessionReturn = { app: null, currentUser: null };

	if (!session) {
		const idToken = await getAppRouterSession();

		if (idToken) {
			session = idToken;
		}

		if (!idToken || !session) {
			return noSessionReturn;
		}
	}

	const decodedIdToken = await adminAuth.verifyIdToken(session);
	const app = initializeAuthenticatedApp(decodedIdToken.uid);
	const auth = getAuth(app);

	const isRevoked = !(await adminAuth.verifyIdToken(session, true).catch((e) => console.error(e.message)));
	if (isRevoked) return noSessionReturn;

	if (auth.currentUser?.uid !== decodedIdToken.uid) {
		const customToken = await adminAuth
			.createCustomToken(decodedIdToken.uid)
			.catch((e) => console.error(e.message));

		if (!customToken) return noSessionReturn;

		await signInWithCustomToken(auth, customToken);
	}

	return { app, currentUser: auth.currentUser };
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

export const adminDb = getFirestore(app);

// const adminApp =
// 	getApps().find((app) => app.name === 'admin') ||
// 	initializeApp(
// 		{
// 			credential: credential.cert(firebaseAdminConfig as ServiceAccount),
// 		},
// 		'admin'
// 	);
// export const adminDb = firestore(adminApp);
