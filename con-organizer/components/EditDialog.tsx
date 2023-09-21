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
import { doc, serverTimestamp, setDoc, updateDoc } from 'firebase/firestore';
import { GameType, Pool } from '@/lib/enums';
import { eventsRef } from '@/lib/firebase';
import { ConEvent } from '@/models/types';
import { Button } from '../lib/mui';
import EventUi from './EventUi';

type Props = {
    open: boolean;
    conEvent?: ConEvent;
    handleClose: () => void;
};

const EditDialog = ({ open, conEvent, handleClose }: Props) => {
    const [title, setTitle] = useState(conEvent?.title || '');
    const [subtitle, setSubtitle] = useState(conEvent?.subtitle || '');
    const [description, setDescription] = useState(conEvent?.description || '');
    const [errorMessage, setErrorMessage] = useState<string>();
    const [published, setPublished] = useState(conEvent?.published || false);
    const [eventPool, setEventPool] = useState(conEvent?.pool || Pool.none);
    const [eventType, setEventType] = useState(conEvent?.gameType || GameType.none);
    const [gameSystem, setGameSystem] = useState<string>(conEvent?.gameSystem || '');
    const [room, setRoom] = useState<string>(conEvent?.room || '');
    const [host, setHost] = useState<string>(conEvent?.host || '');

    useEffect(() => {
        setTitle(conEvent?.title || '');
        setDescription(conEvent?.description || '');
        setSubtitle(conEvent?.subtitle || '');
        setPublished(conEvent?.published || false);
        setEventPool(conEvent?.pool || Pool.none);
        setEventType(conEvent?.gameType || GameType.none);
        setGameSystem(conEvent?.gameSystem || '');
        setRoom(conEvent?.room || '');
        setHost(conEvent?.host || '');
    }, [conEvent]);

    const addEvent = async () => {
        const newEvent = {
            title,
            description,
            subtitle,
            published: published,
            pool: eventPool,
            gameType: eventType,
            room: room,
            host: host,
            gameSystem: gameSystem,
            createdAt: serverTimestamp(),
            lastUpdated: serverTimestamp(),
        };

        try {
            const eventRef = doc(eventsRef);
            await setDoc(eventRef, newEvent);
        } catch (e) {
            const error = e as Error;
            setErrorMessage(error.message);
        }
    };

    async function editEvent(conEvent: ConEvent) {
        const updatedEvent = {
            title,
            description,
            subtitle,
            published: published,
            pool: eventPool,
            gameType: eventType,
            room: room,
            host: host,
            gameSystem: gameSystem,
            createdAt: serverTimestamp(),
            lastUpdated: serverTimestamp(),
        };

        try {
            const eventRef = doc(eventsRef, conEvent.id);
            updateDoc(eventRef, updatedEvent);
        } catch (e) {
            console.error(e);
            const error = e as Error;
            setErrorMessage(error.message);
        }
    }

    return (
        <Dialog open={open} fullWidth={true} maxWidth="md">
            <Box sx={{ minHeight: '900px' }} display="flex" flexDirection="row">
                <Box className="p-4" sx={{ width: '375px', height: '667px' }}>
                    <EventUi conEvent={conEvent || ({} as ConEvent)} />
                </Box>

                <Divider orientation="vertical" variant="middle" flexItem />

                <Box className="p-4">
                    <DialogTitle sx={{ paddingBottom: '0px', paddingLeft: '0px', paddingRight: '0px' }}>
                        {conEvent?.id ? 'Endre arangement' : 'Legg til nytt arangement'}
                    </DialogTitle>
                    <Box sx={{ fontSize: '0.8rem' }}>
                        <div>Opprettet: {conEvent?.createdAt ? conEvent.createdAt.toDate().toLocaleString() : ''} </div>
                        <div>
                            Sist endret: {conEvent?.lastUpdated ? conEvent.lastUpdated.toDate().toLocaleString() : ''}{' '}
                        </div>
                    </Box>

                    <Divider />

                    <DialogContent
                        sx={{
                            display: 'flex',
                            flexDirection: 'column',
                            justifyContent: 'space-between',
                            gap: '1rem',
                        }}
                    >
                        <div>
                            <Switch checked={published} onChange={() => setPublished(!published)} />
                            <span>Publisert</span>
                        </div>
                        <Box sx={{ display: 'flex', flexDirection: 'row', justifyContent: 'space-between' }}>
                            <div>
                                <InputLabel id="pool-select-label">Pulje</InputLabel>
                                <Select
                                    labelId="pool-select-label"
                                    id="pool-select"
                                    value={eventPool}
                                    label="Pulje"
                                    onChange={(e) => setEventPool(e.target.value as Pool)}
                                >
                                    <MenuItem value={Pool.none}>{Pool.none}</MenuItem>
                                    <MenuItem value={Pool.FridayEvening}>{Pool.FridayEvening}</MenuItem>
                                    <MenuItem value={Pool.SaturdayMorning}>{Pool.SaturdayMorning}</MenuItem>
                                    <MenuItem value={Pool.SaturdayEvening}>{Pool.SaturdayEvening}</MenuItem>
                                    <MenuItem value={Pool.SundayMorning}>{Pool.SundayMorning}</MenuItem>
                                </Select>
                            </div>

                            <div>
                                <InputLabel id="type-select-label">Type</InputLabel>
                                <Select
                                    labelId="type-select-label"
                                    id="type-select"
                                    value={eventType}
                                    label="Type"
                                    onChange={(e) => setEventType(e.target.value as GameType)}
                                >
                                    <MenuItem value={GameType.none}>{GameType.none}</MenuItem>
                                    <MenuItem value={GameType.roleplaying}>{GameType.roleplaying}</MenuItem>
                                    <MenuItem value={GameType.boardgame}>{GameType.boardgame}</MenuItem>
                                    <MenuItem value={GameType.other}>{GameType.other}</MenuItem>
                                </Select>
                            </div>
                        </Box>
                    </DialogContent>

                    <DialogContent sx={{ width: '375px', paddingTop: '0' }}>
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
                            autoFocus
                            margin="dense"
                            id="host"
                            label="Arrangør"
                            type="text"
                            fullWidth
                            variant="standard"
                            value={host}
                            onChange={(e) => setHost(e.target.value)}
                        />

                        <Box sx={{ display: 'flex', flexDirection: 'row', justifyContent: 'space-between' }}>
                            <TextField
                                autoFocus
                                margin="dense"
                                id="gameSystem"
                                label="Spillsystem"
                                type="text"
                                fullWidth
                                variant="standard"
                                value={gameSystem}
                                onChange={(e) => setGameSystem(e.target.value)}
                            />
                            <Divider orientation="vertical" variant="middle" flexItem />
                            <TextField
                                autoFocus
                                margin="dense"
                                id="room"
                                label="Rom"
                                type="text"
                                fullWidth
                                variant="standard"
                                value={room}
                                onChange={(e) => setRoom(e.target.value)}
                            />
                        </Box>

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
                            <Button onClick={() => editEvent(conEvent)} type="submit">
                                Save
                            </Button>
                        ) : (
                            <Button onClick={() => addEvent()}>Add</Button>
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
