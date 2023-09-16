'use client';

import { useEffect, useState } from 'react';
import CloseIcon from '@mui/icons-material/Close';
import { Box, Dialog, DialogActions, DialogContent, DialogTitle, Divider, IconButton, TextField } from '@mui/material';
import { CollectionReference, doc, DocumentData, serverTimestamp, setDoc, updateDoc } from 'firebase/firestore';
import { ConEvent } from '@/lib/types';
import { Button } from '../lib/mui';
import EventUi from './eventUi';

interface Props {
    open: boolean;
    conEvent?: ConEvent;
    collectionRef: CollectionReference<DocumentData, DocumentData>;
    handleClose: () => void;
}

const EditDialog = ({ open, conEvent, collectionRef: collectionRef, handleClose }: Props) => {
    const [title, setTitle] = useState(conEvent?.title || '');
    const [subtitle, setSubtitle] = useState(conEvent?.subtitle || '');
    const [description, setDescription] = useState(conEvent?.description || '');
    const [showSelect, setShowSelect] = useState(true);

    useEffect(() => {
        setTitle(conEvent?.title || '');
        setDescription(conEvent?.description || '');
        setSubtitle(conEvent?.subtitle || '');
    }, [conEvent]);

    const addSchool = async () => {
        const newSchool = {
            title,
            description,
            createdAt: serverTimestamp(),
            lastUpdate: serverTimestamp(),
        };

        try {
            const schoolRef = doc(collectionRef);
            await setDoc(schoolRef, newSchool);
        } catch (error) {
            console.error(error);
        }
    };

    async function editEvent(conEvent: ConEvent) {
        const updatedSchool = {
            title: title,
            subtitle: subtitle,
            description: description,
            lastUpdate: serverTimestamp(),
        };

        try {
            const schoolRef = doc(collectionRef, conEvent.id);
            updateDoc(schoolRef, updatedSchool);
        } catch (error) {
            console.error(error);
        }
    }

    return (
        <Dialog open={open} fullWidth={true} maxWidth="md">
            <Box sx={{ height: '900px' }} display="flex" flexDirection="row">
                <Box className="p-4" sx={{ width: '375px', height: '667px' }}>
                    <EventUi conEvent={conEvent || ({} as ConEvent) } showSelect={showSelect} />
                </Box>

                <Divider orientation="vertical" variant="middle" flexItem />

                <Box className="p-4">
                    <DialogTitle>{conEvent?.id ? 'Endre' : 'Legg til'}</DialogTitle>
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
                            <Button onClick={() => addSchool()}>Add</Button>
                        )}
                    </DialogActions>
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
