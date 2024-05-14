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
import DeleteNewEventButton from './DeleteNewEventButton';

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
    const docId = await createId(app, db);

    revalidatePath('/my-events');
    return (
        <>
            <Typography variant="h1">
                Ta en titt nedenfor for å se en liste over arrangementene du har satt sammen.
            </Typography>
            <Box sx={{ position: 'relative', marginTop: '2rem' }}>
                <NewEventButton docId={docId} />
                <Grid2 container spacing="2rem">
                    {docs
                        .sort((a, b) => {
                            return a.createdAt > b.createdAt ? 1 : -1;
                        })
                        .map((doc) => (
                            <Grid2 sx={{ textDecoration: 'none', position: 'relative' }} xs={5.7} md={3} key={doc.id}>
                                <Paper
                                    sx={{
                                        position: 'absolute',
                                        minWidth: '90px',
                                        textAlign: 'center',
                                        rotate: '-25deg',
                                        left: '0',
                                        zIndex: '2',
                                        backgroundColor: doc.isSubmitted ? 'success.main' : 'warning.main',
                                    }}
                                >
                                    <Typography>{doc.isSubmitted ? 'sendt inn' : 'kladd'}</Typography>
                                </Paper>
                                <Box component={Link} href={`/event/create/${doc.id}` as Route}>
                                    <EventCardBig
                                        title={doc.title}
                                        gameMaster={doc.name}
                                        shortDescription={doc.subTitle}
                                        system={doc.system}
                                        backgroundImage="my-events.jpg"
                                    />
                                </Box>
                                <DeleteNewEventButton docId={doc.id} />
                            </Grid2>
                        ))}
                </Grid2>
                <RealtimeMyEvents userId={user.uid} />
            </Box>
        </>
    );
};

export default MyEvents;
