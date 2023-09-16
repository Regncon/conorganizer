'use client';
import { useEffect, useState } from 'react';
import { Button, Card, CardActions } from '@mui/material';
import { collection, onSnapshot } from 'firebase/firestore';
import EditDialog from '@/components/editDialog';
import EventUi from '@/components/eventUi';
import { ConEvent } from '@/lib/types';
import db from '../../../lib/firebase';

type Props = { id: string };

const Event = ({ id }: Props) => {
    const collectionRef = collection(db, 'schools');
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
    }, [collectionRef]);

    const handleCloseEdit = () => {
        setOpenEdit(false);
    };

    const handleOpenEdit = () => {
        setOpenEdit(true);
    };

    const conEvent = conEvents.find((conEvent) => conEvent.id === id) || ({} as ConEvent);
    return (
        <>
            {loading && <h1>Loading...</h1> }
            <EditDialog
                open={openEdit}
                handleClose={handleCloseEdit}
                collectionRef={collectionRef}
                conEvent={conEvent}
            />

            <EventUi conEvent={conEvent} showSelect={true} />

            <Card sx={{ maxWidth: '440px' }}>
                <CardActions>
                    <Button onClick={handleOpenEdit}>Endre</Button>
                </CardActions>
            </Card>
        </>
    );
};

export default Event;
