import CardBase from '$app/dashboard/CardBase';
import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import Grid2 from '@mui/material/Unstable_Grid2/Grid2';
import { collection, doc, type Firestore } from 'firebase/firestore';
import NewEventButton from './NewEventButton';
import type { FirebaseApp } from 'firebase/app';
import Box from '@mui/material/Box';
import { getAllMyEvents } from './actions';
import RealtimeMyEvents from './RealtimeMyEvents';

const createId = async (app: FirebaseApp, db: Firestore) => {
	const collectionRef = collection(db, '_');
	const docRef = doc(collectionRef);
	return docRef.id;
};

const MyEvents = async () => {
	// const newEventDocRef = adminDb.collection('/users').get();
	const { app, user, auth, db } = await getAuthorizedAuth();

	if (app === null || user === null || auth === null || db === null) {
		throw new Error('not Logged inn');
	}

	const docs = await getAllMyEvents(db, user);
	const docId = await createId(app, db);
	return (
		<Box sx={{ position: 'relative', marginTop: '2rem' }}>
			<NewEventButton docId={docId} />
			<Grid2 container spacing="2rem">
				{docs.map((doc) => (
					<Grid2 xs={5.7} md={3} key={doc.id}>
						<CardBase
							href={`/event/create/${doc.id}`}
							description={doc.description}
							img="/my-events.jpg"
							imgAlt="Mine arrangementer"
							title={doc.title}
						/>
					</Grid2>
				))}
			</Grid2>
			<RealtimeMyEvents userId={user.uid} />
		</Box>
	);
};

export default MyEvents;
