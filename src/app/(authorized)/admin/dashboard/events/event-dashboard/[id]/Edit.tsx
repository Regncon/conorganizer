'use client';
import Typography from '@mui/material/Typography';
import Paper from '@mui/material/Paper';
import TextField from '@mui/material/TextField';
import FormControl from '@mui/material/FormControl';
import FormControlLabel from '@mui/material/FormControlLabel';
import FormLabel from '@mui/material/FormLabel';
import Radio from '@mui/material/Radio';
import RadioGroup from '@mui/material/RadioGroup';
import TextareaAutosize from '@mui/material/TextareaAutosize';
import FormGroup from '@mui/material/FormGroup';
import Checkbox from '@mui/material/Checkbox';
import Button from '@mui/material/Button';
import Confetti from 'react-confetti';
import { useCallback, useEffect, useState, type ComponentProps, type FormEvent } from 'react';
import { doc, onSnapshot, updateDoc } from 'firebase/firestore';
import { db, firebaseAuth } from '$lib/firebase/firebase';

import Grid2 from '@mui/material/Unstable_Grid2/Grid2';
import { onAuthStateChanged, type Unsubscribe, type User } from 'firebase/auth';
import Slide from '@mui/material/Slide';
import Skeleton from '@mui/material/Skeleton';
import Snackbar from '@mui/material/Snackbar';
import type { MyNewEvent } from '$lib/types';
import Chip from '@mui/material/Chip';
import NavigateBeforeIcon from '@mui/icons-material/NavigateBefore';
import EventFromSkeleton from '$app/(authorized)/event/create/[id]/EventFormSkeleton';
import MainEvent from '$app/(public)/event/[id]/event';
import Box from '@mui/material/Box';
type Props = {
    id: string;
};
const Edit = ({ id }: Props) => {
    const [isExploding, setIsExploding] = useState(false);

    const [newEvent, setNewEvent] = useState<MyNewEvent>({});
    const [user, setUser] = useState<User | null>();

    const snackBarMessageInitial = 'Din endring er lagra!';
    const [snackBarMessage, setSnackBarMessage] = useState<string>(snackBarMessageInitial);
    const [isSnackBarOpen, setIsSnackBarOpen] = useState<boolean>(false);
    const [tags, setTags] = useState<{ name: keyof MyNewEvent; label: string; selected: boolean }[]>([
        { name: 'childFriendly', label: 'Arrangementet passer for barn', selected: newEvent?.childFriendly ?? false },
        {
            name: 'adultsOnly',
            label: 'Arrangementet passer berre for vaksne (18+)',
            selected: newEvent?.adultsOnly ?? false,
        },
        {
            name: 'beginnerFriendly',
            label: 'Arrangementet er nybyrjarvenleg',
            selected: newEvent?.beginnerFriendly ?? false,
        },
        {
            name: 'possiblyEnglish',
            label: 'Arrangementet kan haldast på engelsk',
            selected: newEvent?.possiblyEnglish ?? false,
        },
        {
            name: 'volunteersPossible',
            label: 'Andre kan halda arrangementet',
            selected: newEvent?.volunteersPossible ?? false,
        },
        {
            name: 'lessThanThreeHours',
            label: 'Eg trur arrangementet vil vare kortare enn 3 timer',
            selected: newEvent?.lessThanThreeHours ?? false,
        },
        {
            name: 'moreThanSixHours',
            label: 'Eg trur arrangementet vil vare lenger enn 6 timer',
            selected: newEvent?.moreThanSixHours ?? false,
        },
    ]);
    //const newEventDocRef = doc(db, 'users', user?.uid ?? '_', 'my-events', id);

    const updateDatabase = async (newEvent: Partial<MyNewEvent>) => {
        // await updateDoc(newEventDocRef, newEvent);
    };

    // useEffect(() => {
    //     let unsubscribeSnapshot: Unsubscribe | undefined;
    //     if (user) {
    //         unsubscribeSnapshot = onSnapshot(newEventDocRef, (snapshot) => {
    //             const newEventData = snapshot.data() as MyNewEvent;
    //             setNewEvent(newEventData);
    //             const newTags = [...tags].map((tag) => ({
    //                 ...tag,
    //                 selected: (newEventData[tag.name] as boolean) ?? false,
    //             }));
    //             setTags(newTags);
    //         });
    //     }

    //     const unsubscribeUser = onAuthStateChanged(firebaseAuth, (user) => {
    //         setUser(user);
    //     });

    //     return () => {
    //         unsubscribeSnapshot?.();
    //         unsubscribeUser();
    //     };
    // }, [user]);

    const event: Event = {
        id: '1',
        title: 'Dette er tittelen',
        system: 'DnD 5e',
        gameMaster: 'Ola Nordmann',
        shortDescription: 'Dette er en kort beskrivelse',
        icons: ['katt', 'hund', 'fugl', 'rollespill', 'nisse', 'visse', 'nisse2', 'nisse3', 'nisse4'],
    };

    const handleSnackBar = (event: React.SyntheticEvent | Event, reason?: string) => {
        if (reason === 'clickaway') {
            return;
        }

        setIsSnackBarOpen(false);
    };

    const handleBlur = useCallback(() => {
        setSnackBarMessage(snackBarMessageInitial);
        setIsSnackBarOpen(true);
    }, [isSnackBarOpen]);

    const handleSubmission = async () => {
        if (newEvent) {
            setIsSnackBarOpen(false);
            setIsExploding(!isExploding);
            await updateDatabase({ isSubmitted: !newEvent.isSubmitted });
            setSnackBarMessage(`Du har  ${newEvent.isSubmitted ? 'meldt av' : 'sendt inn'}  arrangementet`);
            setIsSnackBarOpen(true);
        }
    };
    const handleSubmit: ComponentProps<'form'>['onSubmit'] = async (e) => {
        e.preventDefault();
        if (!newEvent?.isSubmitted) {
            handleSubmission();
        }
    };

    const handleCancelSubmission: ComponentProps<'button'>['onClick'] = async (e) => {
        if (newEvent?.isSubmitted) {
            e.preventDefault();
            handleSubmission();
        }
    };

    const handleOnChange = (e: FormEvent<HTMLFormElement>) => {
        const { value: inputValue, name: inputName, checked, type } = e.target as HTMLInputElement;

        let value: string | boolean = inputValue;
        let name: string = inputName;

        if (type === 'checkbox') {
            value = checked;
        }

        if (type === 'radio') {
            name = 'gameType';
            value = inputName;
        }
        if (user?.email !== null) {
            let payload: Partial<MyNewEvent> = {
                [name]: value,
                updateAt: new Date(Date.now()).toString(),
                updatedBy: user?.email,
            };
            if (newEvent?.isSubmitted) {
                setIsSnackBarOpen(false);
                setSnackBarMessage('du må nå sende inn igjen skjemaet');
                payload = { ...payload, isSubmitted: false };
                setIsSnackBarOpen(true);
            }

            updateDatabase(payload);
        }
    };

    return (
        <>
            <Grid2
                sx={{ marginBlock: '1rem' }}
                noValidate={newEvent.isSubmitted}
                container
                component="form"
                spacing="2rem"
                onBlur={handleBlur}
                onChange={handleOnChange}
                onSubmit={handleSubmit}
            >
                <Grid2 xs={5}>
                    <MainEvent eventData={event} />
                </Grid2>
                <Grid2>
                    <Grid2 xs={5}>
                        <Paper>
                            <FormLabel sx={{ padding: '1rem' }}>Du kan bruke opptil 33 teikn</FormLabel>
                            <TextField
                                name="title"
                                label="Tittel på spelmodul / arrangement"
                                value={newEvent.title}
                                variant="outlined"
                                required
                                inputProps={{
                                    title: 'Minst 1 teikn og maks teikn er 33',
                                }}
                                margin="dense"
                            />
                        </Paper>
                    </Grid2>
                    <Grid2 xs={6}>
                        <Paper>
                            <FormLabel sx={{ padding: '1rem' }}>Du kan bruke opptil 50 teikn</FormLabel>
                            <TextField
                                name="subTitle"
                                value={newEvent.subTitle}
                                label="Kort oppsummering"
                                variant="outlined"
                                required
                                margin="dense"
                                inputProps={{
                                    title: 'Minst 1 teikn og maks 50 teikn',
                                }}
                            />
                        </Paper>
                    </Grid2>
                    <Grid2 xs={6} sm={6} md={3}>
                        <Paper>
                            <TextField
                                name="name"
                                value={newEvent.name}
                                label="Arrangørens namn (Ditt namn)"
                                variant="outlined"
                                required
                            />
                        </Paper>
                    </Grid2>
                    <Grid2 xs={6} sm={6} md={3}>
                        <Paper>
                            <TextField name="system" label="Spillsystem" value={newEvent.system} variant="outlined" />
                        </Paper>
                    </Grid2>
                    <Grid2 xs={6}>
                        <Paper sx={{ padding: '1rem' }}>
                            <FormControl>
                                <FormLabel>Skildring av modulen (tekst til programmet):</FormLabel>
                                <TextareaAutosize minRows={5} name="description" value={newEvent.description} />
                            </FormControl>
                        </Paper>
                    </Grid2>
                    <Grid2 xs={6} sm={4}>
                        <Paper sx={{ padding: '1rem' }}>
                            <FormControl>
                                <FormLabel>Kva type spel er det?</FormLabel>
                                <RadioGroup
                                    value={newEvent.gameType}
                                    aria-labelledby="demo-controlled-radio-buttons-group"
                                    name="controlled-radio-buttons-group"
                                >
                                    <FormControlLabel
                                        value="rolePlaying"
                                        control={<Radio name="rolePlaying" />}
                                        label="rollespel"
                                    />
                                    <FormControlLabel
                                        value="boardGame"
                                        control={<Radio name="boardGame" />}
                                        label="Brettspel"
                                    />
                                    <FormControlLabel
                                        value="cardGame"
                                        control={<Radio name="cardGame" />}
                                        label="Kortspel"
                                    />
                                    <FormControlLabel value="other" control={<Radio name="other" />} label="Annet" />
                                </RadioGroup>
                            </FormControl>
                        </Paper>
                    </Grid2>
                    <Grid2 xs={6}>
                        <Paper sx={{ padding: '1rem' }}>
                            <Typography sx={{ marginBlockStart: '1rem' }}>
                                Trykk på brikka som passar til spelet ditt:
                            </Typography>
                            <Box sx={{ maxWidth: '320px' }}>
                                {tags.map((tag) => (
                                    <Chip
                                        sx={{ marginBlock: '0.4rem', marginInlineEnd: '0.4rem' }}
                                        label={tag.label}
                                        key={tag.name}
                                        color={tag.selected ? 'primary' : 'secondary'}
                                        icon={<NavigateBeforeIcon />}
                                        variant={tag.selected ? 'filled' : 'outlined'}
                                        onClick={async () => {
                                            setTags((prev) =>
                                                prev.map((t) =>
                                                    t.name === tag.name ? { ...t, selected: !t.selected } : t
                                                )
                                            );
                                            await updateDatabase({ [tag.name]: !tag.selected });
                                            setIsSnackBarOpen(true);
                                        }}
                                    />
                                ))}
                            </Box>
                        </Paper>
                    </Grid2>
                    <Grid2 xs={6}>
                        <Paper sx={{ padding: '1rem' }}>
                            <FormControl>
                                <FormLabel>Merknader: Er det noko anna du vil vi skal vite?</FormLabel>
                                <TextareaAutosize
                                    minRows={3}
                                    name="additionalComments"
                                    value={newEvent.additionalComments}
                                />
                            </FormControl>
                        </Paper>
                    </Grid2>

                    <Grid2 xs={6}>
                        <Paper sx={{ padding: '1rem' }}>
                            <Typography>
                                Kladden vert lagra automatisk, men du må trykkje på knappen for å sende inn.
                            </Typography>
                            <Button type="submit" variant="contained" onClick={handleCancelSubmission}>
                                {newEvent.isSubmitted ? 'Meld av' : 'Send inn'}
                            </Button>
                        </Paper>
                    </Grid2>
                </Grid2>
            </Grid2>
            <Snackbar
                anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
                open={isSnackBarOpen}
                onClose={handleSnackBar}
                TransitionComponent={Slide}
                message={snackBarMessage}
                autoHideDuration={1200}
            />
        </>
    );
};
export default Edit;
