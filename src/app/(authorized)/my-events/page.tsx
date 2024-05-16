import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import Grid2 from '@mui/material/Unstable_Grid2/Grid2';
import { collection, doc, type Firestore } from 'firebase/firestore';
import NewEventButton from './NewEventButton';
import type { FirebaseApp } from 'firebase/app';
import Box from '@mui/material/Box';
import { getAllMyEvents } from './actions';
import RealtimeMyEvents from './RealtimeMyEvents';
import type { Route } from 'next';
import { Link, Paper, Typography } from '@mui/material';
import { revalidatePath } from 'next/cache';
import EventCardBig from '$app/(public)/EventCardBig';
import DynamicLink from './DynamicLink';

const createId = async (app: FirebaseApp, db: Firestore) => {
    const collectionRef = collection(db, '_');
    const docRef = doc(collectionRef);
    return docRef.id;
};

const MyEvents = async () => {
    const { app, user, auth, db } = await getAuthorizedAuth();

    if (app === null || user === null || auth === null || db === null) {
        throw new Error('not Logged inn');
    }

    const docs = await getAllMyEvents(db, user);
    const newDocumentId = await createId(app, db);

    revalidatePath('/my-events');
    return (
        <>
            <Typography variant="h1">Sjå under for ein oversikt over arrangementa du har registrert.</Typography>
            <Box sx={{ position: 'relative', marginTop: '2rem' }}>
                <NewEventButton newDocumentId={newDocumentId} />
                <Grid2 container rowGap="0.35rem">
                    {docs
                        .sort((a, b) => {
                            return a.createdAt > b.createdAt ? 1 : -1;
                        })
                        .map((doc) => (
                            <Grid2 sx={{ textDecoration: 'none', position: 'relative' }} xs={12} md={3} key={doc.id}>
                                <DynamicLink docId={doc.id}>
                                    <EventCardBig
                                        title={doc.title}
                                        gameMaster={doc.name}
                                        shortDescription={doc.subTitle}
                                        system={doc.system}
                                        backgroundImage="my-events.jpg"
                                        myEventBar
                                        myEventBarSubmitted={doc.isSubmitted}
                                        myEventDocId={doc.id}
                                    />
                                </DynamicLink>
                            </Grid2>
                        ))}
                </Grid2>
                <RealtimeMyEvents userId={user.uid} />
            </Box>
        </>
    );
};

export default MyEvents;
