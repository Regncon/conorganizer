'use client';

import { useEffect, useState } from 'react';
import CloseIcon from '@mui/icons-material/Close';
import {
    Alert,
    Box,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    Divider,
    IconButton,
    InputLabel,
    MenuItem,
    Select,
    Switch,
    TextField,
} from '@mui/material';
import { CollectionReference, doc, DocumentData, serverTimestamp, setDoc, updateDoc } from 'firebase/firestore';
import { gameType, pool } from '@/lib/enums';
import { ConEvent } from '@/lib/types';
import { Button } from '../lib/mui';
import EventUi from './eventUi';
import { type } from 'os';

type Props = {
    open: boolean;
    conEvent?: ConEvent;
    collectionRef: CollectionReference<DocumentData, DocumentData>;
    handleClose: () => void;
};

const EditDialog = ({ open, conEvent, collectionRef: collectionRef, handleClose }: Props) => {
    const [title, setTitle] = useState(conEvent?.title || '');
    const [subtitle, setSubtitle] = useState(conEvent?.subtitle || '');
    const [description, setDescription] = useState(conEvent?.description || '');
    const [showSelect, setShowSelect] = useState(true);
    const [errorMessage, setErrorMessage] = useState<string>();
    const [published, setPublished] = useState(false);
    const [eventPool, setEventPool] = useState(pool.none);
    const [eventType, setEventType] = useState(gameType.none);

    useEffect(() => {
        setTitle(conEvent?.title || '');
        setDescription(conEvent?.description || '');
        setSubtitle(conEvent?.subtitle || '');
    }, [conEvent]);

    const addevent = async () => {
        const newevent = {
            title,
            description,
            subtitle,
            published: published,
            createdAt: serverTimestamp(),
            lastUpdated: serverTimestamp(),
        };

        try {
            const eventRef = doc(collectionRef);
            await setDoc(eventRef, newevent);
        } catch (e) {
            const error = e as Error;
            setErrorMessage(error.message);
        }
    };

    async function editEvent(conEvent: ConEvent) {
        const updatedevent = {
            title: title,
            subtitle: subtitle,
            published: published,
            description: description,
            lastUpdated: serverTimestamp(),
        };

        try {
            const eventRef = doc(collectionRef, conEvent.id);
            updateDoc(eventRef, updatedevent);
        } catch (e) {
            console.error(e);
            const error = e as Error;
            setErrorMessage(error.message);
        }
    }
 
    console.log(conEvent, "conEvent") 

    return (
        <Dialog open={open} fullWidth={true} maxWidth="md">
            <Box sx={{ height: '900px' }} display="flex" flexDirection="row">
                <Box className="p-4" sx={{ width: '375px', height: '667px' }}>
                    <EventUi conEvent={conEvent || ({} as ConEvent)} showSelect={showSelect} />
                </Box>

                <Divider orientation="vertical" variant="middle" flexItem />

                <Box className="p-4">
                    <DialogTitle>{conEvent?.id ? 'Endre' : 'Legg til'}</DialogTitle>
                    <Divider />
                    <DialogContent
                        sx={{ display: 'flex', flexDirection: 'column', justifyContent: 'space-between', gap: '1rem' }}
                        >
                        <span>Opprettet: {conEvent?.createdAt ? conEvent.createdAt.toDate().toString() : ""} </span>
                        <span>Sist endret: {conEvent?.lastUpdated ? conEvent.lastUpdated.toDate().toString() : ""} </span>
                        <div>
                            <Switch checked={published} onChange={() => setPublished(!published)} />
                            <span>Publisert</span>
                        </div>

                        <div>
                            <InputLabel id="pool-select-label">Pulje</InputLabel>
                            <Select
                                labelId="pool-select-label"
                                id="pool-select"
                                value={eventPool}
                                label="Pulje"
                                onChange={(e) => setEventPool(e.target.value as pool)}
                            >
                                <MenuItem value={pool.none}>{pool.none}</MenuItem>
                                <MenuItem value={pool.FirdayEvening}>{pool.FirdayEvening}</MenuItem>
                                <MenuItem value={pool.SaturdayMorning}>{pool.SaturdayMorning}</MenuItem>
                                <MenuItem value={pool.SaturdayAfternoon}>{pool.SaturdayAfternoon}</MenuItem>
                                <MenuItem value={pool.SundayMorning}>{pool.SundayMorning}</MenuItem>
                            </Select>
                        </div>
                    </DialogContent>

                    <DialogContent sx={{ width: '375px' }}>
                        <TextField
                            autoFocus
                            margin="dense"
                            id="title"
                            label="Tittel"
                            type="text"
                            fullWidth
                            variant="standard"
                            value={title}
                            onChange={(e) => setTitle(e.target.value)}
                        />
                        <TextField
                            autoFocus
                            margin="dense"
                            id="subtitle"
                            label="Undertittel"
                            type="text"
                            fullWidth
                            variant="standard"
                            value={subtitle}
                            onChange={(e) => setSubtitle(e.target.value)}
                        />
                        <TextField
                            margin="dense"
                            id="description"
                            label="Beskrivelse"
                            type="text"
                            fullWidth
                            multiline
                            minRows={10}
                            variant="standard"
                            value={description}
                            onChange={(e) => setDescription(e.target.value)}
                        />
                    </DialogContent>
                    <DialogActions>
                        {conEvent?.id ? (
                            <Button onClick={() => editEvent(conEvent)}>Save</Button>
                        ) : (
                            <Button onClick={() => addevent()}>Add</Button>
                        )}
                    </DialogActions>
                    {!!errorMessage && <Alert severity="error">{errorMessage}</Alert>}
                </Box>

                <Box sx={{ position: 'absolute', top: 0, right: 0 }}>
                    <IconButton onClick={handleClose} aria-label="close">
                        <CloseIcon />
                    </IconButton>
                </Box>
            </Box>
        </Dialog>
    );
};

export default EditDialog;
