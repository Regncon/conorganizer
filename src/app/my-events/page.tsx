import CardBase from '$app/dashboard/CardBase';
import type { NewEvent } from '$app/types';
import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import Paper from '@mui/material/Paper';
import Grid2 from '@mui/material/Unstable_Grid2/Grid2';
import { collection, getFirestore, getDocs, doc, type Firestore } from 'firebase/firestore';
import NewEventButton from './NewEventButton';
import type { FirebaseApp } from 'firebase/app';
import Box from '@mui/material/Box';
import { createUserWithEmailAndPassword, onAuthStateChanged, signInWithEmailAndPassword } from 'firebase/auth';
import { getAuth } from 'firebase-admin/auth';

const createId = async (app: FirebaseApp, db: Firestore) => {
	const collectionRef = collection(db, '_');
	const docRef = doc(collectionRef);
	return docRef.id;
};

const MyEvents = async () => {
	// const newEventDocRef = adminDb.collection('/users').get();
	const { app, currentUser, auth } = await getAuthorizedAuth();
	const user = currentUser;

	if (app === null || currentUser === null || auth === null) {
		throw new Error('asdfg');
	}

	const db = getFirestore(app);
	const ref = collection(db, '/users', currentUser.uid, 'my-events');
	const docs = await getDocs(ref);
	const docId = await createId(app, db);
	return (
		<Box sx={{ position: 'relative', marginTop: '0.5rem' }}>
			<NewEventButton docId={docId} />
			<Grid2 container spacing="2rem">
				{docs.docs.map((doc) => {
					const document = { ...doc.data(), id: doc.id } as NewEvent;
					return (
						<Grid2 xs={5.7} md={3} key={document.id}>
							<CardBase
								href={`/event/create/${document.id}`}
								description={document.description}
								img="/my-events.jpg"
								imgAlt="Mine arrangementer"
								title={document.title}
							/>
						</Grid2>
					);
				})}
			</Grid2>
		</Box>
	);
};

export default MyEvents;
