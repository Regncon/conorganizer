'use client';
import { ConEvent } from '@/lib/types';
import { faDiceD20 } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import {
    Card,
    Box,
    CardHeader,
    CardMedia,
    Chip,
    Collapse,
    IconButton,
    Typography,
    CardContent,
    FormControl,
    FormLabel,
    RadioGroup,
    Radio,
    CardActions,
    Button,
    useRadioGroup,
    styled,
    FormControlLabelProps,
    FormControlLabel,
    Divider,
} from '@mui/material';
import { collection, onSnapshot } from 'firebase/firestore';
import { useState, useEffect } from 'react';
import db from '../../../lib/firebase';
import parse from 'html-react-parser';
import CloseIcon from '@mui/icons-material/Close';
import { Edit } from '@mui/icons-material';
import EditDialog from '@/components/editDialog';
import EventUi from '@/components/eventUi';

type Props = { id: string };

const Event = ({ id }: Props) => {
    const colletionRef = collection(db, 'schools');
    const [conEvents, setconEvents] = useState([] as ConEvent[]);
    const [loading, setLoading] = useState(false);
    const [openEdit, setOpenEdit] = useState(false);

    useEffect(() => {
        setLoading(true);
        const unsub = onSnapshot(colletionRef, (querySnapshot) => {
            const items = [] as ConEvent[];
            querySnapshot.forEach((doc) => {
                items.push(doc.data() as ConEvent);
                items[items.length - 1].id = doc.id;
            });
            setconEvents(items);
            setLoading(false);
        });
        return () => {
            unsub();
        };
    }, []);

    const handleCloseEdit = () => {
        setOpenEdit(false);
    }

    const handleOpenEdit = () => { 
        setOpenEdit(true);
    }

    const conEvent = conEvents.find((conEvent) => conEvent.id === id);
    return (
        <>
        <EditDialog open={openEdit} handleClose={handleCloseEdit} colletionRef={colletionRef} conEvent={conEvent} /> 
        
        <EventUi conEvent={conEvent}/>

        <Card>
        <CardActions>
            <Button onClick={handleOpenEdit}
            >Endre</Button>
        </CardActions>
        </Card>
    </>
    );
};

export default Event;
