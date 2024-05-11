import type { FirebaseOptions } from 'firebase/app';

export const firebaseConfig: FirebaseOptions = {
	apiKey: process.env.NEXT_PUBLIC_FIREBASE_DB_API_KEY,
	authDomain: 'regncon2024.firebaseapp.com',
	projectId: 'regncon2024',
	storageBucket: 'regncon2024.appspot.com',
	messagingSenderId: process.env.NEXT_PUBLIC_FIREBASE_DB_MESSAGING_SENDER_ID,
	appId: process.env.NEXT_PUBLIC_FIREBASE_DB_APP_ID,
};

export const firebaseAdminConfig = {
	type: 'service_account',
	project_id: 'regncon2024',
	private_key_id: process.env.FIREBASE_ADMIN_CLIENT_ID,
	private_key: process.env.FIREBASE_ADMIN_PRIVATE_KEY?.replace(/\\n/gm, '\n'),
	client_email: 'firebase-adminsdk-owlqn@regncon2024.iam.gserviceaccount.com',
	client_id: process.env.FIREBASE_ADMIN_CLIENT_ID,
	auth_uri: 'https://accounts.google.com/o/oauth2/auth',
	token_uri: 'https://oauth2.googleapis.com/token',
	auth_provider_x509_cert_url: 'https://www.googleapis.com/oauth2/v1/certs',
	client_x509_cert_url:
		'https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-owlqn%40regncon2024.iam.gserviceaccount.com',
	universe_domain: 'googleapis.com',
};
