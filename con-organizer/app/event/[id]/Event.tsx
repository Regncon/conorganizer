'use client';
import { useEffect, useState } from 'react';
import { Box, Button, Card, CardActions } from '@mui/material';
import { collection, onSnapshot } from 'firebase/firestore';
import EditDialog from '@/components/editDialog';
import EventUi from '@/components/eventUi';
import { ConEvent } from '@/lib/types';
import db from '../../../lib/firebase';
import MainNavigator from '@/components/mainNavigator';

type Props = { id: string };

const Event = ({ id }: Props) => {
    const collectionRef = collection(db, 'events');
    const [conEvents, setconEvents] = useState([] as ConEvent[]);
    const [openEdit, setOpenEdit] = useState(false);
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        setLoading(true);
        const unSub = onSnapshot(collectionRef, (querySnapshot) => {
            const items = [] as ConEvent[];
            querySnapshot.forEach((doc) => {
                items.push(doc.data() as ConEvent);
                items[items.length - 1].id = doc.id;
            });
            setconEvents(items);
            setLoading(false);
        });
        return () => {
            unSub();
        };
    }, []);

    const handleCloseEdit = () => {
        setOpenEdit(false);
    };

    const handleOpenEdit = () => {
        setOpenEdit(true);
    };

    const conEvent = conEvents.find((conEvent) => conEvent.id === id) || ({} as ConEvent);
    return (
        <Box sx={{ maxWidth: '1080px', margin: '0 auto' }}>
            {loading && <h1>Loading...</h1>}
            <EditDialog
                open={openEdit}
                handleClose={handleCloseEdit}
                collectionRef={collectionRef}
                conEvent={conEvent}
            />

            <Card>
                <Button onClick={() => window.history.go(-1)}>Tilbake</Button>
            </Card>

            <EventUi conEvent={conEvent} showSelect={true} />

            <Card>
                <CardActions>
                    <Button onClick={handleOpenEdit}>Endre</Button>
                </CardActions>
            </Card>
            <MainNavigator />
        </Box>
    );
};

export default Event;
