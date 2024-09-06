import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';

import { collection, doc, type Firestore } from 'firebase/firestore';
import type { FirebaseApp } from 'firebase/app';
import Box from '@mui/material/Box';
import { getAllMyEvents } from './actions';
import RealtimeMyEvents from './RealtimeMyEvents';
import { Grid2, Typography } from '@mui/material';
import { revalidatePath } from 'next/cache';
import DynamicLink from './DynamicLink';
import AddEventCard from './AddEventCard';
import EventCardBig from '$app/(public)/components/EventCardBig';
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
            <Typography sx={{ textAlign: 'center', fontSize: '1.8rem' }}>
                Sjå under for ein oversikt over arrangementa du har registrert.
            </Typography>
            <Box sx={{ position: 'relative', marginTop: '2rem' }}>
                <Grid2 container spacing="2rem">
                    {docs
                        .sort((a, b) => {
                            return new Date(a.createdAt) > new Date(b.createdAt) ? 1 : -1;
                        })
                        .map((doc) => (
                            <Grid2
                                sx={{
                                    textDecoration: 'none',
                                    position: 'relative',
                                    display: 'flex',
                                    placeContent: 'center',
                                    placeItems: 'center',
                                }}
                                key={doc.id}
                                size={{
                                    xl: 2.7,
                                }}
                            >
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
                    <Grid2
                        sx={{
                            display: 'flex',
                            placeContent: 'center',
                            placeItems: 'center',
                        }}
                        size={{
                            xl: 2.7,
                        }}
                    >
                        <AddEventCard newDocumentId={newDocumentId} />
                    </Grid2>
                </Grid2>
                <RealtimeMyEvents userId={user.uid} />
            </Box>
        </>
    );
};

export default MyEvents;
