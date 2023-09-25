import { useEffect, useState } from 'react';
import {
    Alert,
    Box,
    Button,
    Checkbox,
    DialogActions,
    DialogContent,
    DialogTitle,
    Divider,
    FormControlLabel,
    FormGroup,
    InputLabel,
    MenuItem,
    Select,
    Switch,
    TextField,
} from '@mui/material';
import { doc, serverTimestamp, setDoc, updateDoc } from 'firebase/firestore';
import { eventsRef } from '@/lib/firebase';
import { GameType, Pool } from '@/models/enums';
import { ConEvent } from '@/models/types';

type Props = {
    conEvent: ConEvent;
};

const EditUi = ({ conEvent }: Props) => {
    const [title, setTitle] = useState(conEvent?.title || '');
    const [subtitle, setSubtitle] = useState(conEvent?.subtitle || '');
    const [imageUrl, setImageUrl] = useState(conEvent?.imageUrl || '');
    const [description, setDescription] = useState(conEvent?.description || '');
    const [errorMessage, setErrorMessage] = useState<string>();
    const [published, setPublished] = useState(conEvent?.published || false);
    const [eventPool, setEventPool] = useState(conEvent?.pool || Pool.none);
    const [eventType, setEventType] = useState(conEvent?.gameType || GameType.none);
    const [gameSystem, setGameSystem] = useState<string>(conEvent?.gameSystem || '');
    const [room, setRoom] = useState<string>(conEvent?.room || '');
    const [host, setHost] = useState<string>(conEvent?.host || '');
    const [hideEnrollment, setHideEnrollment] = useState<boolean>(conEvent?.hideEnrollment || false);
    const [sortingIndex, setSortingIndex] = useState<number>(conEvent?.sortingIndex || 0);
    const [childFriendly, setChildFriendly] = useState<boolean>(conEvent?.childFriendly || false);
    const [possiblyEnglish, setPossiblyEnglish] = useState<boolean>(conEvent?.possiblyEnglish || false);
    const [adultsOnly, setAdultsOnly] = useState<boolean>(conEvent?.adultsOnly || false);
    const [volunteersPossible, setVolunteersPossible] = useState<boolean>(conEvent?.volunteersPossible || false);
    const [lessThanThreeHours, setLessThanThreeHours] = useState<boolean>(conEvent?.lessThanThreeHours || false);
    const [moreThanSixHours, setMoreThanSixHours] = useState<boolean>(conEvent?.moreThanSixHours || false);
    const [beginnerFriendly, setBeginnerFriendly] = useState<boolean>(conEvent?.beginnerFriendly || false);

    useEffect(() => {
        setTitle(conEvent?.title || '');
        setDescription(conEvent?.description || '');
        setSubtitle(conEvent?.subtitle || '');
        setImageUrl(conEvent?.imageUrl || '');
        setPublished(conEvent?.published || false);
        setEventPool(conEvent?.pool || Pool.none);
        setEventType(conEvent?.gameType || GameType.none);
        setGameSystem(conEvent?.gameSystem || '');
        setChildFriendly(conEvent?.childFriendly || false);
        setPossiblyEnglish(conEvent?.possiblyEnglish || false);
        setAdultsOnly(conEvent?.adultsOnly || false);
        setVolunteersPossible(conEvent?.volunteersPossible || false);
        setLessThanThreeHours(conEvent?.lessThanThreeHours || false);
        setMoreThanSixHours(conEvent?.moreThanSixHours || false);
        setBeginnerFriendly(conEvent?.beginnerFriendly || false);
        setRoom(conEvent?.room || '');
        setHost(conEvent?.host || '');
        setHideEnrollment(conEvent?.hideEnrollment || false);
        setSortingIndex(conEvent?.sortingIndex || 0);
    }, [conEvent]);

    const addEvent = async () => {
        const newEvent = {
            title,
            description,
            subtitle,
            imageUrl,
            published: published,
            pool: eventPool,
            gameType: eventType,
            room: room,
            host: host,
            hideEnrollment: hideEnrollment,
            sortingIndex: sortingIndex,
            gameSystem: gameSystem,
            childFriendly,
            possiblyEnglish,
            adultsOnly,
            volunteersPossible,
            lessThanThreeHours,
            moreThanSixHours,
            beginnerFriendly,
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
            imageUrl,
            published: published,
            pool: eventPool,
            gameType: eventType,
            room: room,
            host: host,
            sortingIndex: sortingIndex,
            hideEnrollment: hideEnrollment,
            gameSystem: gameSystem,
            childFriendly: childFriendly,
            possiblyEnglish: possiblyEnglish,
            adultsOnly: adultsOnly,
            volunteersPossible: volunteersPossible,
            lessThanThreeHours: lessThanThreeHours,
            moreThanSixHours: moreThanSixHours,
            beginnerFriendly: beginnerFriendly,
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
    // throw new Error('test');
    return (
        <>
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
                    <Box sx={{ display: 'flex', flexDirection: 'row', justifyContent: 'space-between' }}>
                        <div>
                            <Switch checked={published} onChange={() => setPublished(!published)} />
                            <span>Publisert</span>
                        </div>
                        <div>
                            <Switch checked={hideEnrollment} onChange={() => setHideEnrollment(!hideEnrollment)} />
                            <span>Skjul påmelding</span>
                        </div>
                    </Box>
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
                            <TextField
                                autoFocus
                                margin="dense"
                                id="sortingIndex"
                                label="Sortering"
                                type="number"
                                fullWidth
                                variant="standard"
                                value={sortingIndex}
                                onChange={(e) => setSortingIndex(Number(e.target.value))}
                            />
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
                <FormGroup>
                    <FormControlLabel
                        control={
                            <Checkbox
                                onChange={(e) => setChildFriendly(Boolean(e.target.checked))}
                                checked={childFriendly}
                            />
                        }
                        label="Barnevennlig"
                    />
                    <FormControlLabel
                        control={
                            <Checkbox
                                onChange={(e) => setPossiblyEnglish(Boolean(e.target.checked))}
                                checked={possiblyEnglish}
                            />
                        }
                        label="Kan kjøres på engelsk"
                    />
                    <FormControlLabel
                        control={
                            <Checkbox onChange={(e) => setAdultsOnly(Boolean(e.target.checked))} checked={adultsOnly} />
                        }
                        label="Kun for voksne (18+)"
                    />
                    <FormControlLabel
                        control={
                            <Checkbox
                                onChange={(e) => setVolunteersPossible(Boolean(e.target.checked))}
                                checked={volunteersPossible}
                            />
                        }
                        label="Kan arrangeres av andre frivillige"
                    />
                    <FormControlLabel
                        control={
                            <Checkbox
                                onChange={(e) => setLessThanThreeHours(Boolean(e.target.checked))}
                                checked={lessThanThreeHours}
                            />
                        }
                        label="Mindre enn tre timer"
                    />
                    <FormControlLabel
                        control={
                            <Checkbox
                                onChange={(e) => setMoreThanSixHours(Boolean(e.target.checked))}
                                checked={moreThanSixHours}
                            />
                        }
                        label="Lengre enn seks timer"
                    />
                    <FormControlLabel
                        control={
                            <Checkbox
                                onChange={(e) => setBeginnerFriendly(Boolean(e.target.checked))}
                                checked={beginnerFriendly}
                            />
                        }
                        label="Nybegynnervennlig"
                    />
                </FormGroup>
                <DialogContent sx={{ width: 'auto', paddingTop: '0' }}>
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
                        id="imageUrl"
                        label="Bilde URL"
                        type="text"
                        fullWidth
                        variant="standard"
                        value={imageUrl}
                        onChange={(e) => setImageUrl(e.target.value)}
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
        </>
    );
};

export default EditUi;
