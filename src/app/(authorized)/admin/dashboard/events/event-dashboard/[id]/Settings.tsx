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
type Props = {
    id?: string;
};
const Settings = ({ id }: Props) => {
    const [isExploding, setIsExploding] = useState(false);

    const [newEvent, setNewEvent] = useState<MyNewEvent>({} as MyNewEvent);
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
                <Grid2 xs={12} sm={6} md={3}>
                    <Paper>
                        <TextField
                            type="email"
                            name="email"
                            value={newEvent.email}
                            label="E-postadresse"
                            variant="outlined"
                            fullWidth
                            disabled
                        />
                    </Paper>
                </Grid2>
                <Grid2 xs={12} sm={6} md={3}>
                    <Paper>
                        <TextField
                            type="phone"
                            name="phone"
                            value={newEvent.phone}
                            label="Kva telefonnummer kan vi nå deg på?"
                            variant="outlined"
                            required
                            fullWidth
                        />
                    </Paper>
                </Grid2>
                <Grid2 xs={12} sm={4}>
                    <Paper sx={{ padding: '1rem', height: '100%' }}>
                        <TextField
                            type="number"
                            name="participants"
                            value={newEvent.participants}
                            label="Maks antall deltakere"
                            variant="outlined"
                            required
                            fullWidth
                        />
                    </Paper>
                </Grid2>
                <Grid2 xs={12} sm={4}>
                    <Paper sx={{ padding: '1rem' }}>
                        <FormGroup>
                            <FormLabel>Kva for pulje kan du arrangere i?</FormLabel>
                            <FormControlLabel
                                control={<Checkbox checked={newEvent.fridayEvening} />}
                                name="fridayEvening"
                                label="Fredag Kveld"
                            />
                            <FormControlLabel
                                control={<Checkbox checked={newEvent.saturdayMorning} />}
                                name="saturdayMorning"
                                label="Lørdag Morgen"
                            />
                            <FormControlLabel
                                control={<Checkbox checked={newEvent.saturdayEvening} />}
                                name="saturdayEvening"
                                label="Lørdag Kveld"
                            />
                            <FormControlLabel
                                control={<Checkbox checked={newEvent.sundayMorning} />}
                                name="sundayMorning"
                                label="Søndag Morgen"
                            />
                        </FormGroup>
                    </Paper>
                </Grid2>
                <Grid2 xs={12}>
                    <Paper sx={{ padding: '1rem' }}>
                        <FormGroup>
                            <FormLabel>Kryss av for det som gjeld</FormLabel>
                            <FormControlLabel
                                control={<Checkbox name="moduleCompetition" checked={newEvent.moduleCompetition} />}
                                label="Eg vil vere med på modulkonkurransen"
                            />
                            <Typography>
                                husk å sende modulen til <a href="mailto:moduler@regncon.no">moduler@regncon.no </a>{' '}
                                innen første september!
                            </Typography>
                        </FormGroup>
                    </Paper>
                </Grid2>
                <Grid2 xs={12}>
                    <Paper sx={{ padding: '1rem' }}>
                        <FormControl fullWidth>
                            <FormLabel>Merknader: Er det noko anna du vil vi skal vite?</FormLabel>
                            <TextareaAutosize
                                minRows={3}
                                name="additionalComments"
                                value={newEvent.additionalComments}
                            />
                        </FormControl>
                    </Paper>
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
export default Settings;
